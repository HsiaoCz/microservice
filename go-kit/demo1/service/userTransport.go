package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// 获取请求
func DecodeUserRequest(c context.Context, r *http.Request) (any, error) {
	if r.URL.Query().Get("uid") != "" {
		uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
		return UserRequest{
			Uid: uid,
		}, nil
	}
	return nil, errors.New("参数错误")
}

// 对响应进行编码
func EncodeUserResponse(c context.Context, w http.ResponseWriter, response any) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
