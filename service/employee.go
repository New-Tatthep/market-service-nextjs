package service

import (
	"market-service/custom_error"
	"market-service/datastore"
	"market-service/request"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/New-Tatthep/microservice/util/uuid"
)

type ServiceAction interface {
	FilterEmployee(req *datastore.FilterData) ([]microservice.Field, error)
	CreateEmployee(input *request.EmployeeRequest) ([]microservice.Field, error)
	UpdateEmployee(input *request.EmployeeRequest) ([]microservice.Field, error)
	DeleteEmployee(req *request.EmployeeRequest) ([]microservice.Field, error)
}

func (sv *service) Connector() ServiceAction {
	return sv
}

func (sv *service) FilterEmployee(req *datastore.FilterData) ([]microservice.Field, error) {
	datalist, totals, err := sv.store.FilterEmployee(req)
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

	if stringutil.IsEmptyString(req.FirstName) {
		return nil, custom_error.New("employee name is required")
	}

	req.UpdateCode = uuid.NewUUID()

	prepareEmployeeDatas := datastore.EmployeeModel{
		UserCode:  uuid.NewUUID(),
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Status:    "active",
	}

	err := sv.store.InsertEmployee(prepareEmployeeDatas)
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

	err := sv.store.UpdateEmployee(datastore.EmployeeModel{
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

	err := sv.store.DeleteEmployee(req.UserCode)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{}, nil
}
