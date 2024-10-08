package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/merch/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"strconv"
)

// @Summary MerchCreate
// @Security ApiKeyAuth
// @Tags musician
// @Description create merch
// @ID create-merch
// @Accept  json
// @Produce  json
// @Param input body dto.MerchWithoutId true "merch info"
// @Param musician_id   path      int  true  "Musician ID"
// @Success 200 {object} dto.CreateMerchResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician/{musician_id}/merch [post]
func MerchCreate(merchUseCase usecase.MerchUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		musicianID := chi.URLParam(r, "musician_id")
		musicianIDUint, err := strconv.ParseUint(musicianID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.MerchWithoutId
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		id, err := merchUseCase.AddMerch(dto.ToModelMerchWithoutId(&req, 0), musicianIDUint)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.CreateMerchResponse{
			Id: id,
		})
	}
}

// @Summary MerchUpdate
// @Security ApiKeyAuth
// @Tags merch
// @Description update merch
// @ID update-merch
// @Accept  json
// @Produce  json
// @Param input body dto.MerchWithoutId true "merch info"
// @Param id   path      int  true  "Merch ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/merch/{id} [put]
func UpdateMerch(useCase usecase.MerchUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchID := chi.URLParam(r, "id")
		merchIDUint, err := strconv.ParseUint(merchID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var req dto.MerchWithoutId
		err = render.DecodeJSON(r.Body, &req)

		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.UpdateMerch(dto.ToModelMerchWithoutId(&req, merchIDUint))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary MerchDelete
// @Security ApiKeyAuth
// @Tags merch
// @Description delete merch
// @ID delete-merch
// @Accept  json
// @Produce  json
// @Param id   path      int  true  "Merch ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/merch/{id} [delete]
func DeleteMerch(useCase usecase.MerchUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		err = useCase.DeleteMerch(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary MerchGet
// @Security ApiKeyAuth
// @Tags merch
// @Description get merch
// @ID get-merch
// @Accept  json
// @Produce  json
// @Param id   path      int  true  "Merch ID"
// @Success 200 {object} dto.MerchWithMusician
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/merch/{id} [get]
func GetMerch(useCase usecase.MerchUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		merch, err := useCase.GetMerch(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		musicianId, err := useCase.GetMusicianForMerch(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.ToDtoMerchWithMusician(merch, musicianId))
	}
}

// @Summary MerchGetAll
// @Security ApiKeyAuth
// @Tags musician
// @Description get all merch
// @ID get-merch-all
// @Accept  json
// @Produce  json
// @Param id   path      int  true  "Musician ID"
// @Success 200 {object} dto.MerchCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/musician/{musician_id}/merch [get]
func GetAllMerchForMusician(useCase usecase.MerchUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "musician_id")
		aid, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		merch, err := useCase.GetAllMerchForMusician(aid)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var res []*dto.Merch

		for _, v := range merch {
			res = append(res, dto.ToDtoMerch(v))
		}

		render.JSON(w, r, dto.MerchCollection{Items: res})
	}
}

// @Summary FindMerch
// @Security ApiKeyAuth
// @Tags merch
// @Description find merch
// @ID find-merch
// @Accept  json
// @Produce  json
// @Param        q    query     string  true  "name search by q"
// @Param        page    query     int  true  "number of page from 1"
// @Param        page_size    query     int  true  "size of page"
// @Success 200 {object} dto.MerchCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/merch [get]
func FindMerch(useCase usecase.MerchUseCase) http.HandlerFunc {
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

		merch, err := useCase.GetMerchByPartName(name, page, pageSize)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var res []*dto.Merch
		for _, v := range merch {
			res = append(res, dto.ToDtoMerch(v))
		}

		render.JSON(w, r, dto.MerchCollection{Items: res})
	}
}
