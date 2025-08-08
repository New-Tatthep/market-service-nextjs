package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/shopspring/decimal"
)

type Product struct {
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Status      string          `json:"status"`
	Quantity    int64           `json:"quantity"`
	Image       UploadFile      `json:"image"`
}

func (md *Product) String() string {
	return stringutil.Json(*md)
}

func (md *Product) ToMap() map[string]any {
	return convutil.Obj2Map(*md)
}

func (value Product) Value() (driver.Value, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (value *Product) Scan(src interface{}) error {
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
