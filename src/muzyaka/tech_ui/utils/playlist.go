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
	playlistPath = "http://localhost:8080/api/playlist/"
)

func CreatePlaylist(client *http.Client,
	query dto.PlaylistWithoutId,
	userId uint64,
	jwt string) error {

	url := userPath + strconv.FormatUint(userId, 10) + "/playlist"

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

	var resp dto.CreatePlaylistResponse
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

//func GetAllAlbums(client *http.Client,
//	musicianId uint64,
//	jwt string) ([]*dto.Album, error) {
//
//	url := musicianPath + strconv.FormatUint(musicianId, 10) + "/album"
//
//	request, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return nil, err
//	}
//	request.Header.Set("Authorization", "Bearer "+jwt)
//
//	respGot, err := client.Do(request)
//	if err != nil {
//		return nil, err
//	}
//	defer respGot.Body.Close()
//
//	var resp dto.AlbumsCollection
//	err = render.DecodeJSON(respGot.Body, &resp)
//
//	if err != nil {
//		return nil, err
//	}
//	if respGot.StatusCode != http.StatusOK {
//		var resp response.Response
//		err = render.DecodeJSON(respGot.Body, &resp)
//		return nil, errors.New(resp.Error)
//	}
//
//	return resp.Albums, nil
//}

func GetAllTracksFromPlaylist(client *http.Client,
	playlistId uint64,
	jwt string) ([]*dto.TrackMeta, error) {

	url := playlistPath + strconv.FormatUint(playlistId, 10) + "/track"

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

	var resp dto.TracksMetaCollection
	err = render.DecodeJSON(respGot.Body, &resp)

	if err != nil {
		return nil, err
	}
	if respGot.StatusCode != http.StatusOK {
		var resp response.Response
		err = render.DecodeJSON(respGot.Body, &resp)
		return nil, errors.New(resp.Error)
	}

	return resp.Tracks, nil
}

func UpdatePlaylist(client *http.Client, query dto.PlaylistWithoutId, playlistId uint64, jwt string) error {
	url := playlistPath + strconv.FormatUint(playlistId, 10)

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

func DeletePlaylist(client *http.Client, playlistId uint64, jwt string) error {
	url := playlistPath + strconv.FormatUint(playlistId, 10)

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
