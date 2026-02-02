package request

import (
	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type ForgotpasswordRequest struct {
	//use for first time
	Username string `json:"username"`
	//use for resend and verify otp
	Code  string `json:"code"`
	Token string `json:"token"`
	// verify otp
	OtpCode    string `json:"otp_code"`
	RefOtpCode string `json:"ref_otp_code"`
	// reset Password
	Password string `json:"password"`

	MachineId string `json:"machine_id"`
	IpAddress string `json:"ip_address"`
}

func (req *ForgotpasswordRequest) String() string {
	return stringutil.Json(*req)
}

func (req *ForgotpasswordRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *ForgotpasswordRequest) Validate() error {
	var errs microservice.Errors

	// TODO

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

type ResetPasswordWithLocalRequest struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Password string `json:"password"`

	// AGENT
	BranchNo  string `json:"branch_no"`
	IpAddress string `json:"ip_address"`
	MachineId string `json:"machine_id"`
	UserAgent string `json:"user_agent"`
}

func (req *ResetPasswordWithLocalRequest) String() string {
	return stringutil.Json(*req)
}

func (req *ResetPasswordWithLocalRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *ResetPasswordWithLocalRequest) Validate() error {
	var errs microservice.Errors

	// TODO

	if len(errs) > 0 {
		return &errs
	}
	return nil
}
