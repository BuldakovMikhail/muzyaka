package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/track/usecase"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"strconv"
)

// @Summary UpdateTrack
// @Security ApiKeyAuth
// @Tags track
// @Description update track
// @ID update-track
// @Accept  json
// @Produce  json
// @Param id path int true "track ID"
// @Param input body dto.TrackObjectWithoutId true "track info"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/track/{id} [put]
func UpdateTrack(useCase usecase.TrackUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackID := chi.URLParam(r, "id")
		trackIDUint, err := strconv.ParseUint(trackID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req dto.TrackObjectWithoutId
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = useCase.UpdateTrack(dto.ToModelTrackObjectWithoutId(&req, trackIDUint))
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary GetTrack
// @Security ApiKeyAuth
// @Tags track
// @Description get track by ID
// @ID get-track
// @Accept  json
// @Produce  json
// @Param id path int true "track ID"
// @Success 200 {object} dto.TrackObject
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/track/{id} [get]
func GetTrack(useCase usecase.TrackUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackID := chi.URLParam(r, "id")
		trackIDUint, err := strconv.ParseUint(trackID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		track, err := useCase.GetTrack(trackIDUint)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, dto.ToDtoTrackObject(track))
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
func DeleteTrack(useCase usecase.TrackUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackID := chi.URLParam(r, "id")
		trackIDUint, err := strconv.ParseUint(trackID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = useCase.DeleteTrack(trackIDUint)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}
