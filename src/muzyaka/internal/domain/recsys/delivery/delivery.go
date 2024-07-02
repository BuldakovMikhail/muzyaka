package delivery

import (
	"github.com/go-chi/render"
	"net/http"
	usecase2 "src/internal/domain/recsys/usecase"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
	"strconv"
)

// @Summary GetRecommendedTracks
// @Security ApiKeyAuth
// @Tags track
// @Description get recommended tracks
// @ID get-recommended-tracks
// @Accept  json
// @Produce  json
// @Param        id    query     uint64  true  "id of song"
// @Param        page    query     int  true  "number of page from 1"
// @Param        page_size    query     int  true  "size of page"
// @Success 200 {object} dto.TracksMetaCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/track/recs [get]
func GetRecommendedTracks(useCase usecase2.RecSysUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pageSizeStr := r.URL.Query().Get("page_size")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tracks, err := useCase.GetSameTracks(id, page, pageSize)

		if err != nil {
			render.JSON(w, r, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var res []*dto.TrackMeta
		for _, v := range tracks {
			res = append(res, dto.ToDtoTrackMeta(v))
		}

		render.JSON(w, r, dto.TracksMetaCollection{Tracks: res})
	}
}
