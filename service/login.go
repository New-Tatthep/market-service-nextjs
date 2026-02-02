package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"market-service/api"
	"market-service/custom_error"
	"market-service/request"
	"market-service/response"
	"strings"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/log"
	"github.com/New-Tatthep/microservice/util/dateutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/New-Tatthep/microservice/util/uuid"
	"github.com/jinzhu/copier"
)

type CloudSessionAction interface {
	CloudLogin(req request.LoginRequest) (*microservice.Field, error)
}

func (sv *service) CloudSession() CloudSessionAction {
	return sv
}

func (sv *service) CloudLogin(req request.LoginRequest) (*microservice.Field, error) {
	var logErr error
	var loginType string
	defer func() {
		err := sv.InsertSignonLog(req.Username, loginType, ActionLogin, req.IpAddress, req.UserAgent, logErr)
		if err != nil {
			sv.ctx.Logger().Debugf("InsertSignOnLog err : %s", err.Error())
		}
	}()

	req.Username = strings.ToLower(req.Username)
	// get user data from db
	userData, err := sv.store.Session().GetUserProfile(req.Username)
	if err != nil {
		logErr = err
		return nil, custom_error.Wrap(err)
	}

	// check user exist
	if stringutil.IsEmptyString(userData.UserCode) {
		err := custom_error.Wrap(errors.New(custom_error.InvalidUsernamePassword))
		logErr = err
		return nil, custom_error.Wrap(err)
	}

	if userData.UserType != userTypeAdmin {
		// validate user status
		if userData.Status == UserStatusInactive {
			err := custom_error.Wrap(errors.New(custom_error.UserInactive))
			logErr = err
			return nil, custom_error.Wrap(err)
		}

	}

	// check exist
	loginExist, err := sv.CheckSessionExist(userData.UserCode)
	if err != nil {
		logErr = err
		return nil, custom_error.Wrap(err)
	}

	countSession, err := sv.CountUserSession(userData.UserCode)
	if err != nil {
		logErr = err
		return nil, custom_error.Wrap(err)
	}

	if loginExist && !req.ForceLogin && countSession >= sv.authConfig.Concurrent {
		err := custom_error.Wrap(errors.New(custom_error.AlreadyLogin))
		logErr = err
		return nil, custom_error.Wrap(err)
	}

	// reset count
	err = sv.store.Session().UpdateLoginFail(nil, userData.UserCode, 0, UserStatusActive)
	if err != nil {
		logErr = err
		return nil, custom_error.Wrap(err)
	}

	var userCache UserCache

	// prepare session data
	sessionId := uuid.NewUUID()
	token := sv.CreateToken(userData.UserCode, sessionId)
	subscribeChannel := fmt.Sprintf("%s:%s#%s", sv.customConfig.SocketSessionNamespace, userData.UserCode, sessionId)

	userCache.Sessions = Sessions{
		Session{
			SessionId:        sessionId,
			LoginType:        LoginLocal,
			IpAddress:        sv.ctx.WebContext().RealIP(),
			MachineId:        "",
			UserAgent:        sv.ctx.WebContext().Request().UserAgent(),
			LoginTime:        dateutil.GetCurrentEpochTime(),
			LastActivityTime: dateutil.GetCurrentEpochTime(),
			LoginFromLocal:   false,
			SubscribeChannel: subscribeChannel,
			Token:            token,
		},
	}

	forceChangePassword := false
	if userData.UserType == userTypeUser || userData.UserType == userTypeAdmin {
		forceChangePassword = userData.ForceChangePassword
	}

	// set user profile
	var userProfile UserProfile
	if err := copier.Copy(&userProfile, userData); err != nil {
		logErr = err
		return nil, custom_error.Wrap(err)
	}
	userCache.UserProfile = userProfile

	// create user cache
	if err := sv.CloudCreateUserCache(userCache); err != nil {
		logErr = err
		return nil, custom_error.Wrap(err)
	}

	// return authen data
	resp := response.LoginResponse{
		Token:               token,
		SubscribeChannel:    subscribeChannel,
		ForceChangePassword: forceChangePassword,
	}

	return &microservice.Field{
		Key:   "data",
		Value: resp,
	}, nil
}

func (sv *service) CheckSessionExist(userCode string) (bool, error) {
	jsonData, err := sv.cache.HGet(RedisUserKey(userCode), RedisSessionField)
	if err != nil {
		return false, custom_error.Wrap(err)
	}

	if stringutil.IsEmptyString(jsonData) {
		// not found
		return false, nil
	} else {
		// found
		return true, nil
	}
}

func (sv *service) CountUserSession(userCode string) (int64, error) {
	jsonData, err := sv.cache.HGet(RedisUserKey(userCode), RedisSessionField)
	if err != nil {
		return 0, custom_error.Wrap(err)
	}

	var sessions Sessions
	if stringutil.IsNotEmptyString(jsonData) {
		if err := json.Unmarshal([]byte(jsonData), &sessions); err != nil {
			return 0, custom_error.Wrap(err)
		}

		lenSession := len(sessions)
		return int64(lenSession), nil
	}
	return 0, nil
}

func (sv *service) CloudCreateUserCache(data UserCache) error {
	if len(data.Sessions) != 1 {
		return custom_error.Wrap(errors.New("session length must be 1"))
	}

	// sessionDuration := time.Duration(sv.authConfig.SessionDuration) * time.Second

	// get old session
	oldSessionJson, err := sv.cache.HGet(RedisUserKey(data.UserProfile.UserCode), RedisSessionField)
	if err != nil {
		return custom_error.Wrap(err)
	}

	var newSessions Sessions
	if stringutil.IsEmptyString(oldSessionJson) {
		sv.ctx.Logger().Debugf("user %s login : no session found, create session id %s", data.UserProfile.UserCode, data.Sessions[0].SessionId)
		newSessions = data.Sessions
	} else {
		// unmarshal
		var oldSession Sessions
		if err := json.Unmarshal([]byte(oldSessionJson), &oldSession); err != nil {
			return custom_error.Wrap(err)
		}

		var removeSessions Sessions
		if sv.authConfig.Concurrent == 1 {
			newSessions = data.Sessions
			removeSessions = oldSession

		} else {
			// check length
			if len(oldSession) >= int(sv.authConfig.Concurrent) {
				// find oldest session
				minIdx := 0
				for idx, session := range oldSession {
					if session.LoginTime < oldSession[minIdx].LoginTime {
						minIdx = idx
					}
				}

				// append to remove list
				removeSessions = append(removeSessions, oldSession[minIdx])

				// remove min index (oldest session)
				newSessions = append(oldSession[:minIdx], oldSession[minIdx+1:]...)

				// append new session
				newSessions = append(newSessions, data.Sessions...)

			} else {
				// new = old + new input
				newSessions = append(oldSession, data.Sessions...)
			}
		}

		sv.ctx.Logger().Debugf("user %s login : found %d/%d session, create session id %s",
			data.UserProfile.UserCode,
			len(oldSession),
			sv.authConfig.Concurrent,
			data.Sessions[0].SessionId,
		)

		// force logout
		for _, removeSession := range removeSessions {
			_, err := sv.RemoveSession(data.UserProfile.UserCode, removeSession.SessionId)
			if err != nil {
				return custom_error.Wrap(err)
			}

			sv.ctx.Logger().Debugf("user %s login : remove session id %s",
				data.UserProfile.UserCode,
				removeSession.SessionId,
			)

			if data.Sessions[0].BranchNo != removeSession.BranchNo {
				if removeSession.LoginFromLocal {
					// case login cloud -> kick local
					// input.BranchNo != removeSession.BranchNo
					// removeSession.LoginFromLocal == true
					// send kafka

					// case login local -> kick other local
					// input.BranchNo != removeSession.BranchNo
					// removeSession.LoginFromLocal == true
					// send kafka

					err := sv.SendForceLogoutMessage(data.UserProfile.UserCode, data.UserProfile.CompanyCode, removeSession)
					if err != nil {
						return custom_error.Wrap(err)
					}

				} else {
					// case login local -> kick cloud
					// input.BranchNo != removeSession.BranchNo
					// removeSession.LoginFromLocal == false
					// send socket

					err = sv.SendLogoutMessage(removeSession.SubscribeChannel)
					if err != nil {
						return custom_error.Wrap(err)
					}

					sv.ctx.Logger().Debugf("user %s login : send logout message to subscribe channel %s",
						data.UserProfile.UserCode,
						removeSession.SessionId,
						removeSession.SubscribeChannel,
					)
				}
			} else {
				if removeSession.LoginFromLocal {
					// case login local -> kick local
					// input.BranchNo == removeSession.BranchNo
					// removeSession.LoginFromLocal == true

					// manage by itself
					sv.ctx.Logger().Debugf("user %s login : concurrent from same channel",
						data.UserProfile.UserCode,
					)
				} else {
					// case login cloud -> kick cloud
					// input.BranchNo == removeSession.BranchNo
					// removeSession.LoginFromLocal == false
					// send socket

					err = sv.SendLogoutMessage(removeSession.SubscribeChannel)
					if err != nil {
						return custom_error.Wrap(err)
					}

					sv.ctx.Logger().Debugf("user %s login : send logout message to subscribe channel %s",
						data.UserProfile.UserCode,
						removeSession.SessionId,
						removeSession.SubscribeChannel,
					)
				}
			}
		}
	}

	// set new session
	// if err := sv.cache.HSetS(RedisUserKey(data.UserProfile.UserCode), RedisSessionField, newSessions.String(), sessionDuration); err != nil {
	// 	return custom_error.Wrap(err)
	// }

	// // update user cache
	// if err := sv.UpdateUserProfileCache(data.UserProfile.UserCode, sessionDuration, data.UserProfile); err != nil {
	// 	return custom_error.Wrap(err)
	// }

	// if err := sv.UpdateUserPermissionCache(data.UserProfile.UserCode, sessionDuration, data.Permission); err != nil {
	// 	return custom_error.Wrap(err)
	// }

	return nil
}

func (sv *service) SendLogoutMessage(channel string) error {
	serverEndpoint := sv.customConfig.SocketServerEndpoint
	apiKey := sv.customConfig.SocketApiKey

	sks, err := microservice.NewSocketServer(
		microservice.WithSocketServerEndpoint(serverEndpoint),
		microservice.WithSocketServerApiKey(apiKey))
	if err != nil {
		return custom_error.Wrap(err)
	}

	if _, err := sks.Publish(channel, &SocketMessageRequest{
		Action:    "logout",
		Timestamp: dateutil.GetCurrentEpochTime(),
	}); err != nil {
		return custom_error.Wrap(err)
	}

	return nil
}

func (sv *service) SendForceLogoutMessage(userCode string, companyCode string, removeSession Session) error {
	rest, err := api.New()
	if err != nil {
		log.Errorf("%s", err.Error())
	}

	reqSyncData := api.ProduceAnyMessageSpecificBranchRequest{
		Topic:       ForceLogoutTopic,
		CompanyCode: companyCode,
		BranchList:  []string{removeSession.BranchNo},
		Messages: []microservice.KafkaMessage{
			{
				Value: &request.ForceLogoutRequest{
					MessageId: uuid.NewUUID(),
					UserCode:  userCode,
					SessionId: removeSession.SessionId,
				},
			},
		},
	}

	err = rest.SendToCloudSpecificBranch(reqSyncData)
	if err != nil {
		log.Errorf("failed to produce force logout message reason %s", err.Error())
	} else {
		sv.ctx.Logger().Debugf("user %s login : send force logout message to branch no %s",
			userCode,
			removeSession.BranchNo,
		)
	}

	return nil
}
