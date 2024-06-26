package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"strconv"
)

var (
	musicianPath = "http://localhost:8080/api/musician/"
	merchPath    = "http://localhost:8080/api/merch/"
)

func CreateMerch(client *http.Client, query dto.MerchWithoutId, musicianId uint64, jwt string) error {
	url := musicianPath + strconv.FormatUint(musicianId, 10) + "/merch"

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

func GetAllMerch(client *http.Client, musicianId uint64, jwt string) ([]*dto.Merch, error) {
	url := musicianPath + strconv.FormatUint(musicianId, 10) + "/merch"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+jwt)

	respGot, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer respGot.Body.Close()

	var resp dto.MerchCollection
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return nil, err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return nil, errors.New(resp.Error)
	}

	return resp.Items, nil
}

func UpdateMerch(client *http.Client, query dto.MerchWithoutId, merchId uint64, jwt string) error {
	url := merchPath + strconv.FormatUint(merchId, 10)

	bodyAsByteArr, err := json.Marshal(query)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyAsByteArr))
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
	if respGot.StatusCode != http.StatusOK {
		return errors.New(resp.Error)
	}

	return nil
}

func DeleteMerch(client *http.Client, merchId uint64, jwt string) error {
	url := merchPath + strconv.FormatUint(merchId, 10)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
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

func GetMerch(client *http.Client, merchId uint64, jwt string) (*dto.MerchWithMusician, error) {
	url := musicianPath + "merch/" + strconv.FormatUint(merchId, 10)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+jwt)

	respGot, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer respGot.Body.Close()

	var resp dto.MerchWithMusician
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return nil, err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return nil, errors.New(resp.Error)
	}

	return &resp, nil
}
