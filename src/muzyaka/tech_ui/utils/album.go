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

func CreateAlbum(client *http.Client,
	query dto.AlbumWithTracks,
	musicianId uint64,
	jwt string) error {

	url := musicianPath + strconv.FormatUint(musicianId, 10) + "/album"

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

	var resp dto.CreateAlbumResponse
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
