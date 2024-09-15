package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/auth/middleware"
	usecase2 "src/internal/domain/auth/usecase"
	"src/internal/domain/user/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"strconv"
)

func CheckIsUserRelated(next http.Handler, userUseCase usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "user_id")
		userIdUint, err := strconv.ParseUint(userId, 10, 64)
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

		if userIdUint != userInfo.Id {
			render.JSON(w, r, response.Error(models.ErrAccessDenied.Error()))
			render.Status(r, http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	}
}
