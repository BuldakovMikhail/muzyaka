package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/playlist/usecase"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"strconv"
)

// @Summary PlaylistCreate
// @Security ApiKeyAuth
// @Tags playlist
// @Description create playlist
// @ID create-playlist
// @Accept  json
// @Produce  json
// @Param user_id   path      int  true  "user ID"
// @Param input body dto.PlaylistWithoutId true "playlist info"
// @Success 200 {object} dto.CreatePlaylistResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id}/playlist [post]
func PlaylistCreate(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "user_id")
		aid, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.PlaylistWithoutId
		err = render.DecodeJSON(r.Body, &req)

		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		id, err := useCase.AddPlaylist(dto.ToModelPlaylistWithoutId(&req, 0), aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.CreatePlaylistResponse{
			Id: id,
		})
	}
}

// @Summary PlaylistUpdate
// @Security ApiKeyAuth
// @Tags playlist
// @Description update playlist
// @ID update-playlist
// @Accept  json
// @Produce  json
// @Param id   path      int  true  "playlist ID"
// @Param input body dto.PlaylistWithoutId true "playlist info"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/playlist/{id} [put]
func UpdatePlaylist(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.PlaylistWithoutId
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.UpdatedPlaylist(dto.ToModelPlaylistWithoutId(&req, aid))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary PlaylistDelete
// @Security ApiKeyAuth
// @Tags playlist
// @Description delete playlist
// @ID delete-playlist
// @Accept  json
// @Produce  json
// @Param id   path      int  true  "playlist ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/playlist/{id} [delete]
func DeletePlaylist(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.DeletePlaylist(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary PlaylistGet
// @Security ApiKeyAuth
// @Tags playlist
// @Description get playlist
// @ID get-playlist
// @Accept  json
// @Produce  json
// @Param id   path      int  true  "playlist ID"
// @Success 200 {object} dto.PlaylistWithUser
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/playlist/{id} [get]
func GetPlaylist(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		playlist, err := useCase.GetPlaylist(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		userId, err := useCase.GetUserForPlaylist(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.ToDtoPlaylistWithUser(playlist, userId))
	}
}

// @Summary AddTrackPlaylist
// @Security ApiKeyAuth
// @Tags playlist
// @Description add track to playlist
// @ID add-track-playlist
// @Accept  json
// @Produce  json
// @Param id path int true "playlist ID"
// @Param input body dto.AddTrackPlaylistRequest true "track info"
// @Success 200 {object} dto.AddTrackPlaylistResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/playlist/{id}/track [post]
func AddTrack(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playlistID := chi.URLParam(r, "id")

		playlistIDUint, err := strconv.ParseUint(playlistID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.AddTrackPlaylistRequest
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.AddTrack(playlistIDUint, req.TrackId)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.AddTrackPlaylistResponse{
			Status: response.StatusOK,
		})
	}
}

// @Summary DeleteTrackPlaylist
// @Security ApiKeyAuth
// @Tags playlist
// @Description delete track from playlist
// @ID delete-track-playlist
// @Accept  json
// @Produce  json
// @Param id path int true "playlist ID"
// @Param track_id path int true "track ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/playlist/{id}/track/{track_id} [delete]
func DeleteTrack(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playlistID := chi.URLParam(r, "id")
		trackID := chi.URLParam(r, "track_id")

		playlistIDUint, err := strconv.ParseUint(playlistID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		trackIDUint, err := strconv.ParseUint(trackID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.DeleteTrack(playlistIDUint, trackIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary GetAllTracksPlaylist
// @Security ApiKeyAuth
// @Tags playlist
// @Description get all tracks from playlist
// @ID get-all-tracks-playlist
// @Accept  json
// @Produce  json
// @Param playlist_id path int true "playlist ID"
// @Success 200 {array} dto.TracksMetaCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/playlist/{playlist_id}/track [get]
func GetAllTracksForPlaylist(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playlistID := chi.URLParam(r, "playlist_id")

		playlistIDUint, err := strconv.ParseUint(playlistID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		tracks, err := useCase.GetAllTracks(playlistIDUint)
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

// @Summary GetAllPlaylistForUser
// @Security ApiKeyAuth
// @Tags playlist
// @Description get all playlists for user
// @ID get-all-playlists-for-user
// @Accept  json
// @Produce  json
// @Param user_id path int true "user ID"
// @Success 200 {array} dto.TracksMetaCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id}/playlist [get]
func GetAllPlaylists(useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")

		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		playlists, err := useCase.GetAllPlaylistsForUser(userIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var res []*dto.Playlist

		for _, v := range playlists {
			res = append(res, dto.ToDtoPlaylist(v))
		}

		render.JSON(w, r, dto.PlaylistsCollection{Playlists: res})
	}
}
