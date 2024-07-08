package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/album/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"strconv"
)

// @Summary GetAlbum
// @Security ApiKeyAuth
// @Tags album
// @Description get album by ID
// @ID get-album
// @Accept  json
// @Produce  json
// @Param id path int true "album ID"
// @Success 200 {object} dto.Album
// @Failure 400,404,405 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/album/{id} [get]
func GetAlbum(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		albumID := chi.URLParam(r, "id")
		albumIDUint, err := strconv.ParseUint(albumID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		album, err := useCase.GetAlbum(albumIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.ToDtoAlbum(album))
	}
}

// @Summary UpdateAlbum
// @Security ApiKeyAuth
// @Tags album
// @Description update album
// @ID update-album
// @Accept  json
// @Produce  json
// @Param id path int true "album ID"
// @Param input body dto.AlbumWithoutId true "album info"
// @Success 200 {object} response.Response
// @Failure 400,404,405 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/album/{id} [put]
func UpdateAlbum(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		albumID := chi.URLParam(r, "id")
		albumIDUint, err := strconv.ParseUint(albumID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.AlbumWithoutId
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.UpdateAlbum(dto.ToModelAlbumWithId(albumIDUint, &req))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary DeleteAlbum
// @Security ApiKeyAuth
// @Tags album
// @Description delete album
// @ID delete-album
// @Accept  json
// @Produce  json
// @Param id path int true "album ID"
// @Success 200 {object} response.Response
// @Failure 400,404,405 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/album/{id} [delete]
func DeleteAlbum(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		albumID := chi.URLParam(r, "id")
		albumIDUint, err := strconv.ParseUint(albumID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.DeleteAlbum(albumIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary AddAlbumWithTracks
// @Security ApiKeyAuth
// @Tags musician
// @Description add album with tracks
// @ID add-album-with-tracks
// @Accept  json
// @Produce  json
// @Param musician_id path int true "musician ID"
// @Param input body dto.AlbumWithTracks true "album info"
// @Success 200 {object} dto.CreateAlbumResponse
// @Failure 400,404,405 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician/{musician_id}/album [post]
func AddAlbumWithTracks(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		musicianID := chi.URLParam(r, "musician_id")
		musicianIDUint, err := strconv.ParseUint(musicianID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.AlbumWithTracks
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var modelTracks []*models.TrackObject
		for _, v := range req.Tracks {
			modelTracks = append(modelTracks, dto.ToModelTrackObjectWithoutId(v, 0, ""))
		}

		albumID, err := useCase.AddAlbumWithTracks(dto.ToModelAlbumWithId(0, &req.AlbumWithoutId), modelTracks, musicianIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.CreateAlbumResponse{Id: albumID})
	}
}

// @Summary CreateTrack
// @Security ApiKeyAuth
// @Tags album
// @Description add track to album
// @ID add-track-to-album
// @Accept  json
// @Produce  json
// @Param id path int true "album ID"
// @Param input body dto.TrackObjectWithoutId true "track info"
// @Success 200 {object} dto.CreateTrackResponse
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/album/{id}/tracks [post]
func CreateTrack(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		albumID := chi.URLParam(r, "id")
		albumIDUint, err := strconv.ParseUint(albumID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.TrackObjectWithoutId
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		trackID, err := useCase.AddTrack(albumIDUint, dto.ToModelTrackObjectWithoutId(&req, 0, ""))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.CreateTrackResponse{Id: trackID})
	}
}

// @Summary GetAllTracks
// @Security ApiKeyAuth
// @Tags album
// @Description get all tracks from album
// @ID get-all-tracks-from-album
// @Accept  json
// @Produce  json
// @Param id path int true "album ID"
// @Success 200 {object} dto.TracksMetaCollection
// @Failure 400,404,405 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/album/{id}/tracks [get]
func GetAllTracks(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		albumID := chi.URLParam(r, "id")
		albumIDUint, err := strconv.ParseUint(albumID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		tracks, err := useCase.GetAllTracks(albumIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var res []*dto.TrackMeta
		for _, v := range tracks {
			res = append(res, dto.ToDtoTrackMeta(v))
		}

		render.JSON(w, r, dto.TracksMetaCollection{Tracks: res})
	}
}

// @Summary GetAllAlbums
// @Security ApiKeyAuth
// @Tags musician
// @Description get all albums
// @ID get-albums-all
// @Accept  json
// @Produce  json
// @Param musician_id   path      int  true  "Musician ID"
// @Success 200 {object} dto.AlbumsCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician/{musician_id}/album [get]
func GetAllAlbumForMusician(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "musician_id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		merch, err := useCase.GetAllAlbumsForMusician(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var res []*dto.Album

		for _, v := range merch {
			res = append(res, dto.ToDtoAlbum(v))
		}

		render.JSON(w, r, dto.AlbumsCollection{Albums: res})
	}
}

// @Summary DeleteTrack
// @Security ApiKeyAuth
// @Tags track
// @Description delete track
// @ID delete-track
// @Accept  json
// @Produce  json
// @Param id path int true "track ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/track/{id} [delete]
func DeleteTrack(useCase usecase.AlbumUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackID := chi.URLParam(r, "id")
		trackIDUint, err := strconv.ParseUint(trackID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.DeleteTrack(trackIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
