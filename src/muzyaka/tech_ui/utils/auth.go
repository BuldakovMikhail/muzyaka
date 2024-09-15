package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
)

const (
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

	data, err := io.ReadAll(respGot.Body)
	if err != nil {
		return "", err
	}
	respFlow := bytes.NewReader(data)

	var resp dto.SignInResponse
	err = render.DecodeJSON(respFlow, &resp)

	if err != nil {
		return "", err
	}
	if respGot.StatusCode != http.StatusOK {
		respFlow := bytes.NewReader(data)
		var resp response.Response
		err = render.DecodeJSON(respFlow, &resp)
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

	data, err := io.ReadAll(respGot.Body)
	if err != nil {
		return "", err
	}
	respFlow := bytes.NewReader(data)

	var resp dto.SignUpResponse
	err = render.DecodeJSON(respFlow, &resp)

	if err != nil {
		return "", err
	}
	if respGot.StatusCode != http.StatusOK {
		respFlow := bytes.NewReader(data)
		var resp response.Response
		err = render.DecodeJSON(respFlow, &resp)
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

	data, err := io.ReadAll(respGot.Body)
	if err != nil {
		return "", err
	}
	respFlow := bytes.NewReader(data)

	var resp dto.SignUpResponse
	err = render.DecodeJSON(respFlow, &resp)

	if err != nil {
		return "", err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		respFlow := bytes.NewReader(data)
		err = render.DecodeJSON(respFlow, &resp)
		return "", errors.New(resp.Error)
	}

	return resp.Token, nil
}
