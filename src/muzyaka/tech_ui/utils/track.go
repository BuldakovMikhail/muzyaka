package utils

import (
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"strconv"
)

var (
	trackPath = "http://localhost:8080/api/track/"
)

func GetTrack(client *http.Client,
	trackId uint64,
	jwt string) (*dto.TrackObject, error) {

	url := trackPath + strconv.FormatUint(trackId, 10)

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

	var resp dto.TrackObject
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
