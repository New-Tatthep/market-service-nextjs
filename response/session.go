package response

import (
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type LoginResponse struct {
	Token               string `json:"token"`
	SubscribeChannel    string `json:"subscribe_channel"`
	ForceChangePassword bool   `json:"force_change_password"`
}

func (req *LoginResponse) String() string {
	return stringutil.Json(*req)
}

func (req *LoginResponse) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type LoginWithLocalResponse struct {
	Token               string  `json:"token"`
	Session             Session `json:"session"`
	ForceChangePassword bool    `json:"force_change_password"`
}

func (req *LoginWithLocalResponse) String() string {
	return stringutil.Json(*req)
}

func (req *LoginWithLocalResponse) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type LoginFailResponse struct {
	CountLoginFail  int64 `json:"count_login_fail"`
	LoginChanceLeft int64 `json:"login_chance_left"`
	Lock            bool  `json:"lock"`
}

func (req *LoginFailResponse) String() string {
	return stringutil.Json(*req)
}

func (req *LoginFailResponse) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type Session struct {
	SessionId        string `json:"session_id"`
	LoginType        string `json:"login_type"`
	IpAddress        string `json:"ip_address"`
	MachineId        string `json:"machine_id"`
	DeviceName       string `json:"device_name"`
	UserAgent        string `json:"user_agent"`
	LoginTime        int64  `json:"login_time"`
	LastActivityTime int64  `json:"last_activity_time"`
	LoginFromLocal   bool   `json:"login_from_local"`
	SubscribeChannel string `json:"subscribe_channel"`
	BranchNo         string `json:"branch"`
	Token            string `json:"token"`
}

func (req *Session) String() string {
	return stringutil.Json(*req)
}

func (req *Session) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type GetUserCacheResponse struct {
	UserProfile           UserProfile                                      `json:"user_profile"`
	Permission            map[string]map[string]map[string]map[string]bool `json:"permissions"`
	FlagIsCheckout        bool                                             `json:"flag_is_checkout"`
	FlagDisableCloseChip  bool                                             `json:"flag_disable_close_chip"`
	FlagDisableCheckStock bool                                             `json:"flag_disable_check_stock"`
	UrlRedirectPageConfig string                                           `json:"url_redirect_page_config"`
	MenuSequence          map[string]string                                `json:"menu_sequence"`
}

func (req *GetUserCacheResponse) String() string {
	return stringutil.Json(*req)
}

func (req *GetUserCacheResponse) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type UserProfile struct {
	UserCode            string            `json:"user_code"`
	UserType            string            `json:"user_type"`
	Username            string            `json:"username"`
	FirstName           string            `json:"first_name"`
	LastName            string            `json:"last_name"`
	Email               string            `json:"email"`
	ProfileImage        File              `json:"profile_image"`
	EndDate             int64             `json:"end_date"`
	EmployeeCode        string            `json:"employee_code"`
	CountLoginFail      int64             `json:"count_login_fail"`
	CompanyList         []Company         `json:"company_list"`
	CompanyCode         string            `json:"company_code"`
	ActiveCheckin       bool              `json:"active_checkin"`
	ForceChangePassword bool              `json:"force_change_password"`
	BranchGroupList     []BranchGroupData `json:"branch_group_list"`
	BranchCashierList   []BranchData      `json:"branch_cashier_list"`

	CloudCounterNames map[string]string `json:"cloud_counter_names"`
}

func (req *UserProfile) String() string {
	return stringutil.Json(*req)
}

func (req *UserProfile) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type File struct {
	RequestID    string `json:"request_id"`
	PublicURL    string `json:"public_url"`
	CdnURL       string `json:"cdn_url"`
	FileSize     int64  `json:"file_size"`
	Folder       string `json:"folder"`
	FileName     string `json:"file_name"`      // uuid
	RealFileName string `json:"real_file_name"` // original file name
	ContextName  string `json:"context_name"`
}

func (req *File) String() string {
	return stringutil.Json(*req)
}

func (req *File) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type Company struct {
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
	PublicURL   string `json:"public_url"`
}

func (req *Company) String() string {
	return stringutil.Json(*req)
}

func (req *Company) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type BranchGroupData struct {
	CompanyCode     string                     `json:"company_code"`
	BranchGroupNo   string                     `json:"branch_group_no"`
	BranchGroupName string                     `json:"branch_group_name"`
	Permission      map[string]map[string]bool `json:"permission"`
	BranchList      []BranchData               `json:"branch_list"`
}

type BranchData struct {
	CompanyCode string `json:"company_code"`
	BranchNo    string `json:"branch_no"`
	BranchName  string `json:"branch_name"`
}
