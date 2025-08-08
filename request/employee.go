package request

import (
	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type EmployeeRequest struct {
	UserCode     string `json:"user_code"`
	UserName     string `json:"user_name"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Status       string `json:"status"`
	ProfileImage string `json:"profile_image"`
	UpdateCode   string `json:"update_code"`
	UpdateTime   int8   `json:"update_time"`
	MobileNo     string `json:"mobile_no"`
	Password     string `json:"password"`
}

func (req *EmployeeRequest) String() string {
	return stringutil.Json(*req)
}

func (req *EmployeeRequest) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

func (req *EmployeeRequest) Validate() error {
	var errs microservice.Errors

	if len(errs) > 0 {
		return &errs
	}

	return nil
}
