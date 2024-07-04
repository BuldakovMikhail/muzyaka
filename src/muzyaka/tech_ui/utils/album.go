package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"strconv"
)

const (
	albumPath = "http://localhost:8080/api/album/"
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

	data, err := io.ReadAll(respGot.Body)
	if err != nil {
		return err
	}
	respFlow := bytes.NewReader(data)

	var resp dto.CreateAlbumResponse
	err = render.DecodeJSON(respFlow, &resp)

	if err != nil {
		return err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		respFlow := bytes.NewReader(data)
		err = render.DecodeJSON(respFlow, &resp)
		return errors.New(resp.Error)
	}

	return nil
}

func GetAllAlbums(client *http.Client,
	musicianId uint64,
	jwt string) ([]*dto.Album, error) {

	url := musicianPath + strconv.FormatUint(musicianId, 10) + "/album"

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

	var resp dto.AlbumsCollection
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

	return resp.Albums, nil
}

func GetAllTracks(client *http.Client,
	albumId uint64,
	jwt string) ([]*dto.TrackMeta, error) {

	url := albumPath + strconv.FormatUint(albumId, 10) + "/tracks"

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

	var resp dto.TracksMetaCollection
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

	return resp.Tracks, nil
}

func UpdateAlbum(client *http.Client, query dto.AlbumWithoutId, albumId uint64, jwt string) error {
	url := albumPath + strconv.FormatUint(albumId, 10)

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

func DeleteAlbum(client *http.Client, albumId uint64, jwt string) error {
	url := albumPath + strconv.FormatUint(albumId, 10)

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
