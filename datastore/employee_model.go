package datastore

import (
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type EmployeeModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (req *EmployeeModel) String() string {
	return stringutil.Json(*req)
}

func (req *EmployeeModel) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}
