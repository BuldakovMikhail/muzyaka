package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"io"
	"net/http"
	url2 "net/url"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"strconv"
)

const (
	trackPath       = "http://localhost:8080/api/track/"
	trackSearchPath = "http://localhost:8080/api/track"
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

	data, err := io.ReadAll(respGot.Body)
	if err != nil {
		return nil, err
	}
	respFlow := bytes.NewReader(data)

	var resp dto.TrackObject
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

func AddTrack(client *http.Client,
	query dto.TrackObjectWithoutId,
	albumId uint64,
	jwt string) error {

	url := albumPath + strconv.FormatUint(albumId, 10) + "/tracks"

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

	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return errors.New(resp.Error)
	}

	return nil
}

func DeleteTrack(client *http.Client, trackId uint64, jwt string) error {
	url := trackPath + strconv.FormatUint(trackId, 10)

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

func FindTracks(client *http.Client,
	query string,
	page int,
	pageSize int,
	jwt string) ([]*dto.TrackMeta, error) {

	url := trackSearchPath

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+jwt)
	request.URL.RawQuery = url2.Values{
		"q":         {query},
		"page":      {strconv.Itoa(page)},
		"page_size": {strconv.Itoa(pageSize)},
	}.Encode()

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

func GetSameTracks(client *http.Client,
	trackId uint64,
	page int,
	pageSize int,
	jwt string) ([]*dto.TrackMeta, error) {

	url := trackPath + "recs"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+jwt)
	request.URL.RawQuery = url2.Values{
		"id":        {strconv.FormatUint(trackId, 10)},
		"page":      {strconv.Itoa(page)},
		"page_size": {strconv.Itoa(pageSize)},
	}.Encode()

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

func GetGenres(client *http.Client, jwt string) ([]string, error) {
	url := trackPath + "genres"

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

	var resp dto.Genres

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

	return resp.Genres, nil
}
