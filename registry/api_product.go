package registry

import (
	"market-service/handler"

	"github.com/New-Tatthep/microservice"
)

const (
	ProductFilterTransactionEndpoint = ServiceV1Endpoint + "/product/filter"
	ProductGetTransactionEndpoint    = ServiceV1Endpoint + "/product/get"
)

func APIProductRegister(ms microservice.IMicroservice) {
	ms.POST(ProductFilterTransactionEndpoint, handler.ProductFilterHandler)
	ms.POST(ProductGetTransactionEndpoint, handler.ProductGetHandler)
}
