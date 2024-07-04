package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/track/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
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

		err = useCase.UpdateTrack(dto.ToModelTrackObjectWithoutId(&req, trackIDUint, ""))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
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
// @Success 200 {object} dto.TrackObjectWithSource
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/track/{id} [get]
func GetTrack(useCase usecase.TrackUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackID := chi.URLParam(r, "id")
		trackIDUint, err := strconv.ParseUint(trackID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		track, err := useCase.GetTrack(trackIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.ToDtoTrackObjectWithSource(track))
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

// @Summary FindTracks
// @Security ApiKeyAuth
// @Tags track
// @Description find tracks
// @ID find-tracks
// @Accept  json
// @Produce  json
// @Param        q    query     string  true  "name search by q"
// @Param        page    query     int  true  "number of page from 1"
// @Param        page_size    query     int  true  "size of page"
// @Success 200 {object} dto.TracksMetaCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/track [get]
func FindTracks(useCase usecase.TrackUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("q")
		if name == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(models.ErrInvalidParameter.Error()))
			return
		}

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		pageSizeStr := r.URL.Query().Get("page_size")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		tracks, err := useCase.GetTracksByPartName(name, page, pageSize)
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
