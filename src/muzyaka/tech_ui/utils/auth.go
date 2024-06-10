package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	"src/internal/models/dto"
)

var (
	authPath = "http://localhost:8080/api/auth/"
)

func SignIn(client *http.Client, login string, password string) (string, error) {
	url := authPath + "sign-in"

	reqBody := dto.SignIn{
		Email:    login,
		Password: password,
	}

	reqBodyJson, _ := json.Marshal(reqBody)
	respGot, err := client.Post(url, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", err
	}
	defer respGot.Body.Close()

	var resp dto.SignInResponse
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return "", err
	}
	if respGot.StatusCode != http.StatusOK {
		return "", errors.New(respGot.Status)
	}

	return resp.Token, nil
}
