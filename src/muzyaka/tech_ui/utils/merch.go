package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"strconv"
)

var (
	merchPath = "http://localhost:8080/api/musician/"
)

func CreateMerch(client *http.Client, query dto.MerchWithoutId, musicianId uint64, jwt string) error {
	url := merchPath + strconv.FormatUint(musicianId, 10) + "/merch"

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
	defer respGot.Body.Close()
	if err != nil {
		return err
	}

	var resp dto.CreateMerchResponse
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return errors.New(resp.Error)
	}

	return nil
}
