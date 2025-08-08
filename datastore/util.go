package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/New-Tatthep/microservice"
)

type FilterData struct {
	Filters map[string]interface{} `json:"filters"`
}

func (obj *FilterData) Validate() error {
	return nil
}

func queryOptionBuilder(searchQuery squirrel.SelectBuilder, option microservice.IQueryOption) squirrel.SelectBuilder {
	if option != nil {
		if option.GetLimit() > 0 {
			searchQuery = searchQuery.Limit(uint64(option.GetLimit()))
		}
		if option.GetOffset() > 0 {
			searchQuery = searchQuery.Offset(uint64(option.GetOffset()))
		}

		//sorting
		if len(option.GetSort()) > 0 {
			orders := make([]string, 0)
			for _, order := range option.GetSort() {
				orderStr := fmt.Sprintf("%s %s", order.GetField(), string(order.GetOrder()))
				orders = append(orders, orderStr)
			}
			searchQuery = searchQuery.OrderBy(orders...)
		}
	}

	return searchQuery
}

type UploadFile struct {
	RequestID        string `json:"request_id"`
	Context          string `json:"context"`
	Folder           string `json:"folder"`
	FileName         string `json:"file_name"`
	OriginalFileName string `json:"original_file_name"`
	Key              string `json:"key"`
	FileSize         int64  `json:"file_size"`
	PublicURL        string `json:"public_url"`
	CdnURL           string `json:"cdn_url"`
}

func (value UploadFile) Value() (driver.Value, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (value *UploadFile) Scan(src interface{}) error {
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
