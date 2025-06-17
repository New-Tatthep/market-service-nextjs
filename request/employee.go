package request

import (
	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type EmployeeRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
