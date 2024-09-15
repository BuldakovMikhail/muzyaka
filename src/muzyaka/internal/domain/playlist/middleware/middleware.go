package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/auth/middleware"
	usecase2 "src/internal/domain/auth/usecase"
	"src/internal/domain/playlist/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"strconv"
)

func CheckPlaylistOwnership(next http.Handler, useCase usecase.PlaylistUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playlistID := chi.URLParam(r, "id")
		playlistIDUint, err := strconv.ParseUint(playlistID, 10, 64)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		userInfo, isOk := r.Context().Value(middleware.ValuesFromContext).(middleware.ContextValues)
		if !isOk {
			render.JSON(w, r, response.Error(models.ErrInvalidContext.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}
		if userInfo.Role == usecase2.AdminRole {
			next.ServeHTTP(w, r)
			return
		}

		isAllowed, err := useCase.IsPlaylistOwned(playlistIDUint, userInfo.Id)
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
