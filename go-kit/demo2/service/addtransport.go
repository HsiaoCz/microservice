package service

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SumRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
