package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndpoint(UserService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := "nothing"
		if r.Method == "GET" {
			result = UserService.GetName(r.Uid)
		}
		if r.Method == "DELETE" {
			if err := UserService.DeleteUser(r.Uid); err != nil {
				result = err.Error()
				return result, nil
			}
			result = fmt.Sprintf("userid为%d的删除成功", r.Uid)
		}
		return UserResponse{Result: result}, nil
	}
}
