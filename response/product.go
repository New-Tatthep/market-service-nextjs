package response

import (
	"github.com/New-Tatthep/microservice/util/convutil"
	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/shopspring/decimal"
)

type FilterProductResponse struct {
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Status      string          `json:"status"`
	Quantity    int64           `json:"quantity"`
	Image       UploadFile      `json:"image"`
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

func (md *FilterProductResponse) String() string {
	return stringutil.Json(*md)
}

func (md *FilterProductResponse) ToMap() map[string]any {
	return convutil.Obj2Map(*md)
}
