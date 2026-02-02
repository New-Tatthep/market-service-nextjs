package service

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"market-service/custom_error"
	"market-service/datastore"
	"market-service/response"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/dateutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/New-Tatthep/microservice/util/uuid"
)

const (
	// userType
	userTypeAd    = "ad"
	userTypeAdmin = "admin"
	userTypeUser  = "user"

	forgotpasswordTypeCheck  = "check"
	forgotpasswordTypeReset  = "reset"
	forgotpasswordTypeVerify = "verify"

	// redis key name
	RedisSessionField        = "session"
	RedisUserProfileField    = "user_profile"
	RedisPermissionField     = "permission"
	RedisForgotPasswordField = "forgot_password"

	// login type
	LoginLocal = "local"
	LoginCloud = "cloud"

	// user status
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusLock     = "lock"

	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusDelete   = "deleted"
	StatusLock     = "lock"

	// kafka topic
	ForceLogoutTopic                    = "ForceLogout"
	KafkaSyncCounterDataName            = "CounterData"
	KafkaSyncLocalUserPasswordDataTopic = "SyncLocalUserPassword"
	KafkaSyncCloudUserPasswordDataTopic = "SyncCloudUserPassword"

	// action
	ActionLogin  = "login"
	ActionLogout = "logout"

	CompanyMainTypeCodeMain = "main"
	CompanyMainTypeCodeSub  = "sub"
)

func (sv *service) CheckInternet() bool {
	// timeout := time.Duration(time.Duration(sv.customConfig.DialTimeout) * time.Second)
	// conn, err := net.DialTimeout("udp", sv.customConfig.DialToIp, timeout)
	// if err != nil {
	// 	return false
	// }
	// defer conn.Close()

	// fmt.Println()
	// TODO remove mock no internet
	// return false

	timeout := time.Duration(time.Duration(sv.customConfig.DialTimeout) * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get(sv.customConfig.DialToIp)

	return err == nil
}

func countDigit(number int64) int {
	count := 0
	for number != 0 {
		number /= 10
		count++
	}
	return count
}

func DefaultEpoch2DateString(epoch int64) (string, error) {
	epochTime := time.Unix(epoch/int64(time.Microsecond), 0)
	return fmt.Sprintf("%d-%02d-%02d",
		epochTime.Year(),
		epochTime.Month(),
		epochTime.Day(),
	), nil
}

func validateEmployeeExpireTime(endDate int64) bool {

	currentTime := dateutil.GetCurrentEpochTime()
	count := countDigit(endDate)

	if count != 13 && count != 10 {
		return false
	}

	endDateStr, _ := DefaultEpoch2DateString(endDate)
	endDateStr = fmt.Sprintf("%s 00:00:00", endDateStr)

	dateEpouch, err := dateutil.DefaultDateTimeString2Epoch(endDateStr)
	if err != nil {
		return false
	}

	if dateEpouch > currentTime {
		return true
	}
	return false
}

func (sv *service) CreateToken(userCode string, sessionId string) string {
	plainToken := fmt.Sprintf("%s_%s", userCode, sessionId)
	token := base64.URLEncoding.EncodeToString([]byte(plainToken))

	return token
}

func (sv *service) ExtractToken(token string) (string, string, error) {
	plainTokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", "", err
	}

	splitedToken := strings.Split(string(plainTokenBytes), "_")
	if len(splitedToken) != 2 {
		return "", "", errors.New("invalid token")
	}

	return splitedToken[0], splitedToken[1], nil
}

func (sv *service) UpdateUserProfileCache(userCode string, sessionDuration time.Duration, userProfileData UserProfile) error {
	if err := sv.cache.HSetS(RedisUserKey(userCode), RedisUserProfileField, userProfileData.String(), sessionDuration); err != nil {
		return custom_error.Wrap(err)
	}

	return nil
}

func (sv *service) UpdateUserPermissionCache(userCode string, sessionDuration time.Duration, permission map[string]map[string]map[string]map[string]bool) error {
	permissionDataBytes, err := json.Marshal(permission)
	if err != nil {
		return custom_error.Wrap(err)
	}
	permissionData := string(permissionDataBytes)

	if err := sv.cache.HSetS(RedisUserKey(userCode), RedisPermissionField, permissionData, sessionDuration); err != nil {
		return custom_error.Wrap(err)
	}

	return nil
}

func (sv *service) ExtendUserSessionDuration(userCode string, sessionId string) error {
	userSessions, err := sv.GetUserSession(userCode)
	if err != nil {
		return custom_error.Wrap(err)
	}

	// find session
	found, _, idx := userSessions.FindSession(sessionId)
	if !found {
		return custom_error.Wrap(errors.New("session not found"))
	}

	userSessions[idx].LastActivityTime = dateutil.GetCurrentEpochTime()

	// save to redis
	sessionDuration := time.Duration(sv.authConfig.SessionDuration) * time.Second
	if err := sv.cache.HSetS(RedisUserKey(userCode), RedisSessionField, userSessions.String(), sessionDuration); err != nil {
		return custom_error.Wrap(err)
	}

	return nil
}

// check expired & find session index & extend user session (stamp new last activity)
func (sv *service) CheckExpiredSession(userCode string, sessionId string) (*UserCache, int, error) {
	userCache, err := sv.GetUserCache(userCode)
	if err != nil {
		return nil, -1, err
	}

	found, userSession, idx := userCache.Sessions.FindSession(sessionId)
	if !found {
		return nil, -1, custom_error.Wrap(errors.New(custom_error.SessionNotFound))
	}

	nextExpiredTime := time.Unix(userSession.LastActivityTime/int64(time.Microsecond), 0).Add(time.Duration(sv.authConfig.SessionDuration) * time.Second)
	result := nextExpiredTime.Compare(time.Now())
	if result == -1 {
		return nil, -1, custom_error.Wrap(errors.New(custom_error.SessionExpired))
	}

	// extend cache dutaion
	if err := sv.ExtendUserSessionDuration(userCode, sessionId); err != nil {
		return nil, -1, custom_error.Wrap(err)
	}

	return userCache, idx, nil
}

func (sv *service) GetUserCache(userCode string) (*UserCache, error) {
	sessionsJson, err := sv.cache.HGet(RedisUserKey(userCode), RedisSessionField)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	var sessions Sessions
	if stringutil.IsNotEmptyString(sessionsJson) {
		if err := json.Unmarshal([]byte(sessionsJson), &sessions); err != nil {
			return nil, custom_error.Wrap(err)
		}
	}

	userProfileJson, err := sv.cache.HGet(RedisUserKey(userCode), RedisUserProfileField)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	var userProfile UserProfile
	if stringutil.IsNotEmptyString(userProfileJson) {
		if err := json.Unmarshal([]byte(userProfileJson), &userProfile); err != nil {
			return nil, custom_error.Wrap(err)
		}
	}

	permissionJson, err := sv.cache.HGet(RedisUserKey(userCode), RedisPermissionField)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	var permission map[string]map[string]map[string]map[string]bool
	if stringutil.IsNotEmptyString(permissionJson) && permissionJson != "null" {
		if err := json.Unmarshal([]byte(permissionJson), &permission); err != nil {
			return nil, custom_error.Wrap(err)
		}
	}

	userCache := &UserCache{
		Sessions:    sessions,
		UserProfile: userProfile,
		Permission:  permission,
	}

	return userCache, nil
}

func (sv *service) GetUserSession(userCode string) (Sessions, error) {
	sessionsJson, err := sv.cache.HGet(RedisUserKey(userCode), RedisSessionField)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	// unmarshal
	var sessions Sessions
	if err := json.Unmarshal([]byte(sessionsJson), &sessions); err != nil {
		return nil, custom_error.Wrap(err)
	}

	return sessions, nil
}

func (sv *service) GetUserProfile(userCode string) (*UserProfile, error) {
	userProfileJson, err := sv.cache.HGet(RedisUserKey(userCode), RedisUserProfileField)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	var userProfile UserProfile
	if stringutil.IsNotEmptyString(userProfileJson) {
		if err := json.Unmarshal([]byte(userProfileJson), &userProfile); err != nil {
			return nil, custom_error.Wrap(err)
		}
	}

	return &userProfile, nil
}

func (sv *service) RemoveSession(userCode string, sessionId string) (*Session, error) {
	// get user session with token from local
	oldSessionsJson, err := sv.cache.HGet(RedisUserKey(userCode), RedisSessionField)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	if stringutil.IsEmptyString(oldSessionsJson) {
		return nil, nil
	}

	// unmarshal
	var oldSessions Sessions
	if err := json.Unmarshal([]byte(oldSessionsJson), &oldSessions); err != nil {
		return nil, custom_error.Wrap(err)
	}

	// remove session
	newSessions := oldSessions
	var removeSession Session
	var found bool
	for idx, oldSession := range oldSessions {
		if oldSession.SessionId == sessionId {
			newSessions = append(newSessions[:idx], newSessions[idx+1:]...)
			removeSession = oldSession
			found = true
			break
		}
	}

	// case no session in local
	if !found {
		return nil, nil
	}

	// update session
	if len(newSessions) == 0 {
		// remove all
		fields, _, err := sv.cache.HScan(RedisUserKey(userCode), 0, "", 100)
		if err != nil {
			return nil, custom_error.Wrap(err)
		}

		if err := sv.cache.HDel(RedisUserKey(userCode), fields...); err != nil {
			return nil, custom_error.Wrap(err)
		}
	} else {
		// update
		sessionDuration := time.Duration(sv.authConfig.SessionDuration) * time.Second

		if err := sv.cache.HSetS(RedisUserKey(userCode), RedisSessionField, newSessions.String(), sessionDuration); err != nil {
			return nil, custom_error.Wrap(err)
		}
	}

	return &removeSession, nil
}

func (sv *service) GetToken() (string, error) {
	token := sv.ctx.WebContext().Request().Header[http.CanonicalHeaderKey("x-auth-token")]

	if token == nil {
		return "", custom_error.Wrap(errors.New(custom_error.TokenNotFound))
	}

	if stringutil.IsEmptyString(token[0]) {
		return "", custom_error.Wrap(errors.New(custom_error.TokenNotFound))
	}

	return token[0], nil
}

func (sv *service) LoginFailCount(userData datastore.EmployeeModel) (*microservice.Field, error) {
	count := userData.CountLoginFail + 1
	lock := false

	err := sv.store.Session().UpdateLoginFail(nil, userData.UserCode, count, UserStatusActive)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	resp := &microservice.Field{
		Key: "data",
		Value: response.LoginFailResponse{
			CountLoginFail:  count,
			LoginChanceLeft: sv.authConfig.MaxLoginFail - count,
			Lock:            lock,
		},
	}

	if lock {
		return resp, custom_error.Wrap(errors.New(custom_error.UserLocked))
	} else {
		return resp, custom_error.Wrap(errors.New(custom_error.InvalidUsernamePassword))
	}
}

func ToNullString(str string, canNull bool) sql.NullString {
	if canNull {
		if stringutil.IsNotEmptyString(str) {
			return sql.NullString{
				String: str,
				Valid:  true,
			}
		} else {
			return sql.NullString{}
		}
	} else {
		return sql.NullString{
			String: str,
			Valid:  true,
		}
	}
}

func (sv *service) InsertSignonLog(username string, loginType string, action string, ipAddress string, userAgent string, err error) error {
	var status, errMsg string
	if err != nil {
		status = "fail"
		errMsg = err.Error()
	} else {
		status = "success"
	}

	input := datastore.InsertSignOnLogInput{
		LogCode:  uuid.NewUUID(),
		Username: username,
		IpAddress: sql.NullString{
			Valid:  true,
			String: ipAddress,
		},
		UserAgent: sql.NullString{
			Valid:  true,
			String: userAgent,
		},
		MachineId:  sql.NullString{},
		LoginType:  ToNullString(loginType, true),
		Action:     action,
		Status:     status,
		LastUpdate: dateutil.GetCurrentEpochTime(),
		Message:    ToNullString(errMsg, true),
	}

	dbErr := sv.store.Log().InsertSignOnLog(nil, input)
	if dbErr != nil {
		return custom_error.Wrap(dbErr)
	}

	return nil
}

func RedisUserKey(userCode string) string {
	return fmt.Sprintf("user_%s", userCode)
}

func RedisUserForgotPasswordKey(token string) string {
	return fmt.Sprintf("user_forgotpassword_%s", token)
}

func RedisCloudReqUserForgotPasswordKey(token string) string {
	return fmt.Sprintf("user_cloud_forgotpassword_%s", token)
}

func generateChar(length int) string {
	rand.Seed(time.Now().UnixNano())

	min := 97  // ASCII code for 'a'
	max := 122 // ASCII code for 'z'

	var randomString string

	for i := 0; i < length; i++ {
		randomInt := rand.Intn(max-min+1) + min
		randomChar := rune(randomInt)
		randomString += string(randomChar)
	}

	return randomString
}

func IsCurrentDay(epochTimestamp int64) bool {
	// Convert the epoch timestamp to time.Time
	timestampTime := time.UnixMilli(epochTimestamp)

	// Get the current time
	currentTime := time.Now()

	// Compare the year, month, and day
	return timestampTime.Year() == currentTime.Year() &&
		timestampTime.Month() == currentTime.Month() &&
		timestampTime.Day() == currentTime.Day()
}

func ConvertEpochToDateStringNoFormat(epochTimestamp int64) string {
	// Convert the epoch timestamp to time.Time
	timestampTime := time.UnixMilli(epochTimestamp)

	resultData := fmt.Sprintf("%d%d%d", timestampTime.Year(), int(timestampTime.Month()), timestampTime.Day())

	// Compare the year, month, and day
	return resultData
}
