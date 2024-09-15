package utils

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"io"
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

func DislikeTrack(client *http.Client,
	query dto.Dislike,
	userId uint64,
	jwt string) error {

	url := userPath + strconv.FormatUint(userId, 10) + "/favorite"

	bodyAsByteArr, err := json.Marshal(query)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("DELETE", url, bytes.NewBuffer(bodyAsByteArr))
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

func GetLikedTracks(client *http.Client,
	userId uint64,
	jwt string) ([]*dto.TrackMeta, error) {

	url := userPath + strconv.FormatUint(userId, 10) + "/favorite"

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
	if respGot.StatusCode == http.StatusNotFound {
		return nil, errors.New(models.ErrNotFound.Error())
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		respFlow := bytes.NewReader(data)
		err = render.DecodeJSON(respFlow, &resp)
		return nil, errors.New(resp.Error)
	}

	return resp.Tracks, nil
}

func IsTrackLiked(client *http.Client,
	userId uint64,
	trackId uint64,
	jwt string) (bool, error) {

	url := userPath + strconv.FormatUint(userId, 10) + "/favorite/" + strconv.FormatUint(trackId, 10)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	request.Header.Set("Authorization", "Bearer "+jwt)

	respGot, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer respGot.Body.Close()

	data, err := io.ReadAll(respGot.Body)
	if err != nil {
		return false, err
	}
	respFlow := bytes.NewReader(data)

	var resp dto.IsLikedResponse
	err = render.DecodeJSON(respFlow, &resp)

	if err != nil {
		return false, err
	}
	if respGot.StatusCode == http.StatusNotFound {
		return false, errors.New(models.ErrNotFound.Error())
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		respFlow := bytes.NewReader(data)
		err = render.DecodeJSON(respFlow, &resp)
		return false, errors.New(resp.Error)
	}

	return resp.IsLiked, nil
}
