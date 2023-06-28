package utils

import (
	"errors"

	"github.com/imroc/req/v3"
)

var client = req.C()

func getClient() *req.Client {
	client.BaseURL = API_URL_V1

	return client
}

func GetMe(token string) (*MeRequest, error) {
	client = getClient()

	var me MeRequest
	var errorResult ErrorAPIResponse
	resp, err := client.R().SetHeader("Authorization", "Bearer "+token).SetErrorResult(&errorResult).SetSuccessResult(&me).Get("/users/@me")

	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, errors.New(errorResult.Error.Message)
	}

	return &me, nil
}
