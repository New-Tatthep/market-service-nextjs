package service

import (
	"new-service/custom_error"
	"new-service/datastore"
	"new-service/request"

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

	if stringutil.IsEmptyString(req.Name) {
		return nil, custom_error.New("employee name is required")
	}

	err := sv.store.InsertEmployee(datastore.EmployeeModel{
		ID:   uuid.NewUUID(),
		Name: req.Name,
	})
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{}, nil
}

func (sv *service) UpdateEmployee(req *request.EmployeeRequest) ([]microservice.Field, error) {

	if stringutil.IsEmptyString(req.ID) {
		return nil, custom_error.New("employee ID is required")
	}

	if stringutil.IsEmptyString(req.Name) {
		return nil, custom_error.New("employee name is required")
	}

	err := sv.store.UpdateEmployee(datastore.EmployeeModel{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{}, nil
}

func (sv *service) DeleteEmployee(req *request.EmployeeRequest) ([]microservice.Field, error) {

	if stringutil.IsEmptyString(req.ID) {
		return nil, custom_error.New("employee ID is required")
	}

	err := sv.store.DeleteEmployee(req.ID)
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	return []microservice.Field{}, nil
}
