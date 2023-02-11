package main

import (
	"mircoservice/go-kit/demo1/service"
	"net/http"

	httpTransport "github.com/go-kit/kit/transport/http"
)

func main() {
	user := service.UserService{}
	endp := service.GenUserEndpoint(user)

	serverHandler := httpTransport.NewServer(endp, service.DecodeUserRequest, service.EncodeUserResponse)
	http.ListenAndServe(":3202", serverHandler)
}
