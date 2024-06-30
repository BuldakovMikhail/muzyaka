package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"strconv"
)

const (
	userPath  = "http://localhost:8080/api/user/"
	basicPath = "http://localhost:8080/api/"
)

func GetMe(client *http.Client, jwt string) (*dto.GetMeResponse, error) {
	url := basicPath + "get-me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	respGot, err := client.Do(req)
	defer respGot.Body.Close()
	if err != nil {
		return nil, err
	}

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

func LikeTrack(client *http.Client,
	query dto.Like,
	userId uint64,
	jwt string) error {

	url := userPath + strconv.FormatUint(userId, 10) + "/favorite"

	bodyAsByteArr, err := json.Marshal(query)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyAsByteArr))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+jwt)

	respGot, err := client.Do(request)
	if err != nil {
		return err
	}
	defer respGot.Body.Close()

	var resp response.Response
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return err
	}
	if respGot.StatusCode == http.StatusNotFound {
		return errors.New(models.ErrNotFound.Error())
	}
	if respGot.StatusCode != http.StatusOK {
		return errors.New(resp.Error)
	}

	return nil
}
