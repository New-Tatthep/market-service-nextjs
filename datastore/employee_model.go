package datastore

import (
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type EmployeeModel struct {
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

func (req *EmployeeModel) String() string {
	return stringutil.Json(*req)
}

func (req *EmployeeModel) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}
