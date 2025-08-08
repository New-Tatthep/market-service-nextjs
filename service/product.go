package service

import (
	"market-service/custom_error"
	"market-service/datastore"
	"market-service/response"

	"github.com/New-Tatthep/microservice"
)

func (sv *service) FilterProduct(request microservice.FilterRequest) ([]*response.FilterProductResponse, int64, error) {
	resp := make([]*response.FilterProductResponse, 0)
	var total int64 = 0
	if err := sv.store.Do(func(action datastore.IAction) error {
		datalist, totals, err := sv.store.FilterProduct(request.GetFilters(), request.GetOption())
		if err != nil {
			return custom_error.Wrap(err)
		}

		prepareImageInfo := response.UploadFile{}

		for _, data := range datalist {
			prepareImageInfo.PublicURL = data.Image.PublicURL
			resp = append(resp, &response.FilterProductResponse{
				Code:        data.Code,
				Name:        data.Name,
				Description: data.Description,
				Price:       data.Price,
				Status:      data.Status,
				Quantity:    data.Quantity,
				Image:       prepareImageInfo,
			})
		}
		total = totals
		return nil

	}); err != nil {
		return resp, -1, err
	}

	return resp, total, nil
}
