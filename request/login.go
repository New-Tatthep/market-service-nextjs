package request

import (
	"errors"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	MachineId string `json:"machine_id"`

	ForceLogin bool `json:"force_login"`
}

func (req *LoginRequest) String() string {
	return stringutil.Json(*req)
}

func (req *LoginRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *LoginRequest) Validate() error {
	var errs microservice.Errors

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

func (req *LoginRequest) ValidateLocal() error {
	var errs microservice.Errors

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

type LoginWithLocalRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	BranchNo   string `json:"branch_no"`
	IpAddress  string `json:"ip_address"`
	MachineId  string `json:"machine_id"`
	UserAgent  string `json:"user_agent"`
	ForceLogin bool   `json:"force_login"`

	// FORLOGOUT
	IsLogout  bool   `json:"is_logout"`
	UserCode  string `json:"user_code"`
	SessionId string `json:"session_id"`

	// IsCheckMachineId
	CompanyCode      string `json:"company_code"`
	IsCheckMachineId bool   `json:"is_check_machine_id"`
}

func (req *LoginWithLocalRequest) String() string {
	return stringutil.Json(*req)
}

func (req *LoginWithLocalRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *LoginWithLocalRequest) Validate() error {
	var errs microservice.Errors

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

type UpdateUserCacheRequest struct {
	Username string `json:"username"`
}

func (req *UpdateUserCacheRequest) String() string {
	return stringutil.Json(*req)
}

func (req *UpdateUserCacheRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *UpdateUserCacheRequest) Validate() error {
	var errs microservice.Errors

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

type LogoutRequest struct {
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

func (req *LogoutRequest) String() string {
	return stringutil.Json(*req)
}

func (req *LogoutRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *LogoutRequest) Validate() error {
	var errs microservice.Errors

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

type CheckUserPasswordRequest struct {
	Password string `json:"password"`
	UserCode string `json:"user_code"`
}

func (req *CheckUserPasswordRequest) String() string {
	return stringutil.Json(*req)
}

func (req *CheckUserPasswordRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *CheckUserPasswordRequest) Validate() error {
	var errs microservice.Errors

	if stringutil.IsEmptyString(req.Password) {
		return errors.New("password is required")
	}

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

type GetUserInfoForSocketEventRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *GetUserInfoForSocketEventRequest) String() string {
	return stringutil.Json(*req)
}

func (req *GetUserInfoForSocketEventRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *GetUserInfoForSocketEventRequest) Validate() error {
	var errs microservice.Errors

	if stringutil.IsEmptyString(req.Username) {
		return errors.New("username is required")
	}

	if stringutil.IsEmptyString(req.Password) {
		return errors.New("password is required")
	}

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

type ForceLogoutRequest struct {
	MessageId string `json:"message_id"`
	UserCode  string `json:"user_code"`
	SessionId string `json:"session_id"`
	BranchNo  string `json:"branch_no"`
	Sender    string `json:"sender"`
}

func (req *ForceLogoutRequest) String() string {
	return stringutil.Json(*req)
}

func (msg ForceLogoutRequest) ToJson() string {
	return stringutil.Json(msg)
}

func (req *ForceLogoutRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *ForceLogoutRequest) Validate() error {
	var errs microservice.Errors

	if stringutil.IsEmptyString(req.UserCode) {
		errs = append(errs, microservice.NewError("REQUIRED_USER_CODE", "user_code is required"))
	}

	if stringutil.IsEmptyString(req.SessionId) {
		errs = append(errs, microservice.NewError("REQUIRED_SESSION_ID", "session_id is required"))
	}

	if len(errs) > 0 {
		return &errs
	}
	return nil
}
