package utils

import (
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	"src/internal/models/dto"
)

func GetMe(client *http.Client, jwt string) (*dto.GetMeResponse, error) {
	url := "http://localhost:8080/api/get-me"

	respGot, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer respGot.Body.Close()

	var resp dto.GetMeResponse
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return nil, err
	}
	if respGot.StatusCode != http.StatusOK {
		return nil, errors.New(respGot.Status)
	}

	return &resp, nil
}
