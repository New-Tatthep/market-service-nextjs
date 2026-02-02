package service

import (
	"market-service/custom_error"
	"market-service/datastore"
	"market-service/request"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/New-Tatthep/microservice/util/uuid"
)

type EmployeeServiceAction interface {
	FilterEmployee(req *datastore.FilterData) ([]microservice.Field, error)
	CreateEmployee(input *request.EmployeeRequest) ([]microservice.Field, error)
	UpdateEmployee(input *request.EmployeeRequest) ([]microservice.Field, error)
	DeleteEmployee(req *request.EmployeeRequest) ([]microservice.Field, error)
}

func (sv *service) EmployeeAction() EmployeeServiceAction {
	return sv
}

func (sv *service) FilterEmployee(req *datastore.FilterData) ([]microservice.Field, error) {
	datalist, totals, err := sv.store.EmployeeAction().FilterEmployee(req)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{
		{
			Key:   "datas",
			Value: datalist,
		},
		{
			Key:   "total",
			Value: totals,
		},
	}, nil
}

func (sv *service) CreateEmployee(req *request.EmployeeRequest) ([]microservice.Field, error) {

	if stringutil.IsEmptyString(req.UserName) {
		return nil, custom_error.New("user name is required")
	}

	prepareEmployeeDatas := datastore.EmployeeModel{
		UserCode:   uuid.NewUUID(),
		UserName:   req.UserName,
		UserType:   req.UserType,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Status:     "active",
		Password:   req.Password,
		CreateCode: "TT",
		// CreateTime: dateutil.GetCurrentEpochTime(),
		UpdateCode: "TT",
		// UpdateTime: dateutil.GetCurrentEpochTime(),
	}

	err := sv.store.EmployeeAction().InsertEmployee(prepareEmployeeDatas)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{}, nil
}

func (sv *service) UpdateEmployee(req *request.EmployeeRequest) ([]microservice.Field, error) {

	if stringutil.IsEmptyString(req.UserCode) {
		return nil, custom_error.New("employee UserCode is required")
	}

	if stringutil.IsEmptyString(req.FirstName) {
		return nil, custom_error.New("employee FirstName is required")
	}

	err := sv.store.EmployeeAction().UpdateEmployee(datastore.EmployeeModel{
		UserCode:  req.UserCode,
		FirstName: req.FirstName,
	})
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{}, nil
}

func (sv *service) DeleteEmployee(req *request.EmployeeRequest) ([]microservice.Field, error) {

	if stringutil.IsEmptyString(req.UserCode) {
		return nil, custom_error.New("employee ID is required")
	}

	err := sv.store.EmployeeAction().DeleteEmployee(req.UserCode)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{}, nil
}
