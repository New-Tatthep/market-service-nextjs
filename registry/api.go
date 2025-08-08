package registry

import (
	"market-service/handler"

	"github.com/New-Tatthep/microservice"
)

const (
	ServiceV1Endpoint = "/market-service/v1"
)

const (
	cacheContextName = "rediscache"
	kafkaContextName = "mq"
)

func APIRegister(ms microservice.IMicroservice) {
	ms.POST(ServiceV1Endpoint+"/employee/filter", handler.FilterEmployeeHandler)
	ms.POST(ServiceV1Endpoint+"/employee/get", handler.FilterEmployeeHandler)
	ms.POST(ServiceV1Endpoint+"/employee/create", handler.CreateEmployeeHandler)
	ms.POST(ServiceV1Endpoint+"/employee/update", handler.UpdateEmployeeHandler)
	ms.POST(ServiceV1Endpoint+"/employee/delete", handler.DeleteEmployeeHandler)
}
