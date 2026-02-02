package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/New-Tatthep/microservice"
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/go-resty/resty/v2"
)

type ProduceAnyMessageAllBranchRequest struct {
	Topic       string                      `json:"topic"`
	CompanyCode string                      `json:"company_code"`
	Messages    []microservice.KafkaMessage `json:"messages"`
}

func (msg ProduceAnyMessageAllBranchRequest) String() string {
	return stringutil.Json(msg)
}

func (msg ProduceAnyMessageAllBranchRequest) ToJson() string {
	return stringutil.Json(msg)
}

func (msg ProduceAnyMessageAllBranchRequest) ToMap() map[string]any {
	return convutil.Obj2Map(msg)
}

type ProduceAnyMessageSpecificBranchRequest struct {
	Topic       string                      `json:"topic"`
	CompanyCode string                      `json:"company_code"`
	BranchList  []string                    `json:"branch_list"`
	Messages    []microservice.KafkaMessage `json:"messages"`
}

func (msg ProduceAnyMessageSpecificBranchRequest) String() string {
	return stringutil.Json(msg)
}

func (msg ProduceAnyMessageSpecificBranchRequest) ToJson() string {
	return stringutil.Json(msg)
}

func (msg ProduceAnyMessageSpecificBranchRequest) ToMap() map[string]any {
	return convutil.Obj2Map(msg)
}

func (api *API) SendToCloudAllBranch(body ProduceAnyMessageAllBranchRequest) error {
	endpoint := api.customConfig.ProduceCloudAllBranchUrl

	if stringutil.IsEmptyString(endpoint) {
		return errors.New("endpoint is required")
	}

	client := resty.New()

	resp, err := client.R().
		SetBody(body).
		Post(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%d : %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}

func (api *API) SendToCloudSpecificBranch(body ProduceAnyMessageSpecificBranchRequest) error {
	endpoint := api.customConfig.ProduceCloudSpecificBranchUrl

	if stringutil.IsEmptyString(endpoint) {
		return errors.New("endpoint is required")
	}

	client := resty.New()

	resp, err := client.R().
		SetBody(body).
		Post(endpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%d : %s", resp.StatusCode(), string(resp.Body()))
	}

	return nil
}
