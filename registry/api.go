package registry

import (
	"new-service/handler"

	"github.com/New-Tatthep/microservice"
)

const (
	ServiceV1Endpoint = "/new-service/v1"
)

const (
	cacheContextName = "rediscache"
	kafkaContextName = "mq"
)

func APIRegister(ms microservice.IMicroservice) {
	ms.POST(ServiceV1Endpoint+"/employee/filter", handler.FilterEmployeeHandler)
	ms.POST(ServiceV1Endpoint+"/employee/create", handler.CreateEmployeeHandler)
	ms.POST(ServiceV1Endpoint+"/employee/update", handler.UpdateEmployeeHandler)
	ms.POST(ServiceV1Endpoint+"/employee/delete", handler.DeleteEmployeeHandler)

}
