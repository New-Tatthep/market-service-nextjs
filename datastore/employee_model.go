package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
)

type EmployeeModel struct {
	UserCode            string `json:"user_code"`
	UserName            string `json:"user_name"`
	UserType            string `json:"user_type"`
	Password            string `json:"password"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	ProfileImage        File
	MobileNo            string `json:"mobile_no"`
	BirthDate           int64  `json:"birth_date"`
	CreateCode          string `json:"create_code"`
	CreateTime          int64  `json:"create_time"`
	UpdateCode          string `json:"update_code"`
	UpdateTime          int64  `json:"update_time"`
	Status              string `json:"status"`
	CountLoginFail      int64
	ForceChangePassword bool
	// Additional fields can be added here
}

func (req *EmployeeModel) String() string {
	return stringutil.Json(*req)
}

func (req *EmployeeModel) ToMap() map[string]any {
	return convutil.Obj2Map(*req)
}

type File struct {
	RequestID    string `json:"request_id"`
	PublicURL    string `json:"public_url"`
	CdnURL       string `json:"cdn_url"`
	FileSize     int64  `json:"file_size"`
	Folder       string `json:"folder"`
	FileName     string `json:"file_name"`
	RealFileName string `json:"real_file_name"`
	ContextName  string `json:"context_name"`
}

func (value File) Value() (driver.Value, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (value *File) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	data, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte but got %T", src)
	}

	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	return nil
}
