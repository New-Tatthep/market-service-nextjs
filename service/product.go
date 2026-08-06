package service

import (
	"market-service/custom_error"
	"market-service/response"

	"github.com/New-Tatthep/microservice"
)

type ProductServiceAction interface {
	FilterProduct(request microservice.FilterRequest) ([]microservice.Field, error)
}

func (sv *service) ProductAction() ProductServiceAction {
	return sv
}

func (sv *service) FilterProduct(request microservice.FilterRequest) ([]microservice.Field, error) {
	resp := make([]response.FilterProductResponse, 0)
	// var total int64 = 0
	// if err := sv.store.Do(func(action datastore.IAction) error {
	datalist, totals, err := sv.store.ProductAction().FilterProduct(request.GetFilters(), request.GetOption())
	if err != nil {
		return nil, custom_error.Wrap(err)
	}

	prepareImageInfo := response.UploadFile{}

	for _, data := range datalist {
		prepareImageInfo.PublicURL = data.Image.PublicURL
		resp = append(resp, response.FilterProductResponse{
			Code:        data.Code,
			Name:        data.Name,
			Description: data.Description,
			Price:       data.Price,
			Status:      data.Status,
			Quantity:    data.Quantity,
			Image:       prepareImageInfo,
		})
	}

	return []microservice.Field{
		{
			Key:   "datas",
			Value: resp,
		},
		{
			Key:   "total",
			Value: totals,
		},
	}, nil
}
