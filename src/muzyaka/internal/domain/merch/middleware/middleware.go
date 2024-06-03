package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/auth/middleware"
	usecase3 "src/internal/domain/auth/usecase"
	"src/internal/domain/merch/usecase"
	usecase2 "src/internal/domain/musician/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"strconv"
)

func CheckMerchOwnership(next http.Handler,
	useCase usecase.MerchUseCase,
	musicianUseCase usecase2.MusicianUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		merchID := chi.URLParam(r, "id")
		merchIDUint, err := strconv.ParseUint(merchID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userInfo, isOk := r.Context().Value(middleware.ValuesFromContext).(middleware.ContextValues)
		if !isOk {
			render.JSON(w, r, response.Error(models.ErrInvalidContext.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}
		if userInfo.Role == usecase3.AdminRole {
			next.ServeHTTP(w, r)
			return
		}

		musicianId, err := musicianUseCase.GetMusicianIdForUser(userInfo.Id)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			render.Status(r, http.StatusInternalServerError)
			return
		}

		isAllowed, err := useCase.IsMerchOwned(merchIDUint, musicianId)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			render.Status(r, http.StatusInternalServerError)
			return
		}

		if !isAllowed {
			render.JSON(w, r, response.Error(models.ErrAccessDenied.Error()))
			render.Status(r, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	}
}
