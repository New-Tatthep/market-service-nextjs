package registry

import (
	"market-service/handler"

	"github.com/New-Tatthep/microservice"
)

const (
	ProductFilterTransactionEndpoint = ServiceV1Endpoint + "/product/filter"
)

func APIProductRegister(ms microservice.IMicroservice) {
	ms.POST(ProductFilterTransactionEndpoint, handler.ProductFilterHandler)

}
