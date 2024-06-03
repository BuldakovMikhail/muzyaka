package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/musician/usecase"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"strconv"
)

// @Summary CreateMusician
// @Security ApiKeyAuth
// @Tags musician
// @Description create musician
// @ID create-musician
// @Accept  json
// @Produce  json
// @Param input body dto.MusicianWithoutId true "musician info"
// @Success 200 {object} dto.CreateMusicianResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician [post]
func CreateMusician(musicianUseCase usecase.MusicianUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.MusicianWithoutId
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := musicianUseCase.AddMusician(dto.ToModelMusicianWithoutId(&req, 0))
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, dto.CreateMusicianResponse{
			Id: id,
		})
	}
}

// @Summary MusicianUpdate
// @Security ApiKeyAuth
// @Tags musician
// @Description update musician
// @ID update-musician
// @Accept  json
// @Produce  json
// @Param input body dto.MusicianWithoutId true "musician info"
// @Param musician_id   path      int  true  "Musician ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician/{musician_id} [put]
func UpdateMusician(musicianUseCase usecase.MusicianUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.MusicianWithoutId
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := chi.URLParam(r, "musician_id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = musicianUseCase.UpdatedMusician(dto.ToModelMusicianWithoutId(&req, aid))
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary MusicianDelete
// @Security ApiKeyAuth
// @Tags musician
// @Description delete musician
// @ID delete-musician
// @Accept  json
// @Produce  json
// @Param musician_id   path      int  true  "Musician ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician/{musician_id} [delete]
func DeleteMusician(musicianUseCase usecase.MusicianUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "musician_id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = musicianUseCase.DeleteMusician(aid)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary MusicianGet
// @Security ApiKeyAuth
// @Tags musician
// @Description get musician
// @ID get-musician
// @Accept  json
// @Produce  json
// @Param musician_id   path      int  true  "Musician ID"
// @Success 200 {object} dto.Musician
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician/{musician_id} [get]
func GetMusician(musicianUseCase usecase.MusicianUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "musician_id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mus, err := musicianUseCase.GetMusician(aid)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, dto.ToDtoMusician(mus))
	}
}
