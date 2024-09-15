package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"strconv"
)

func GetMusician(client *http.Client, musicianId uint64, jwt string) (*dto.Musician, error) {
	url := musicianPath + strconv.FormatUint(musicianId, 10)

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

	data, err := io.ReadAll(respGot.Body)
	if err != nil {
		return nil, err
	}
	respFlow := bytes.NewReader(data)

	var resp dto.Musician
	err = render.DecodeJSON(respFlow, &resp)

	if err != nil {
		return nil, err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		respFlow := bytes.NewReader(data)
		err = render.DecodeJSON(respFlow, &resp)
		return nil, errors.New(resp.Error)
	}

	return &resp, nil
}

func UpdateMusician(client *http.Client, query dto.MusicianWithoutId, musicianId uint64, jwt string) error {
	url := musicianPath + strconv.FormatUint(musicianId, 10)

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
