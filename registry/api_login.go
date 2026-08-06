package registry

import (
	"market-service/handler"

	"github.com/New-Tatthep/microservice"
)

const (
	LoginEndpoint  = ServiceV1Endpoint + "/login"
	LogoutEndpoint = ServiceV1Endpoint + "/logout"
)

func ApiLoginRegister(ms microservice.IMicroservice) {
	// both cloud and local
	ms.POST(LoginEndpoint, handler.LoginHandler)
	ms.POST(LogoutEndpoint, handler.LogoutHandler)
}
