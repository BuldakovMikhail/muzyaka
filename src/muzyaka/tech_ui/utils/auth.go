package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
)

var (
	authPath = "http://localhost:8080/api/auth/"
)

func SignIn(client *http.Client, query dto.SignIn) (string, error) {
	url := authPath + "sign-in"

	reqBody := query

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
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return "", errors.New(resp.Error)
	}

	return resp.Token, nil
}

func SignUpAsUser(client *http.Client, query dto.SignUp) (string, error) {
	url := authPath + "sign-up/user"

	reqBody := query

	reqBodyJson, _ := json.Marshal(reqBody)
	respGot, err := client.Post(url, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", err
	}
	defer respGot.Body.Close()

	var resp dto.SignUpResponse
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return "", err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return "", errors.New(resp.Error)
	}

	return resp.Token, nil
}

func SignUpAsMusician(client *http.Client, query dto.SignUpMusician) (string, error) {
	url := authPath + "sign-up/musician"

	reqBody := query

	reqBodyJson, _ := json.Marshal(reqBody)
	respGot, err := client.Post(url, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", err
	}
	defer respGot.Body.Close()

	var resp dto.SignUpResponse
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return "", err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return "", errors.New(resp.Error)
	}

	return resp.Token, nil
}
