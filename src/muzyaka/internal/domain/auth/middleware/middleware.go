package middleware

import (
	"context"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/auth/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"strings"
)

type ContextValues struct {
	Id   uint64
	Role string
}

const ValuesFromContext = "ContextValues"

func JwtParseMiddleware(next http.Handler, useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		tokenModel := models.AuthToken{Secret: []byte(token)}

		id, role, err := useCase.BasicAuthorization(&tokenModel)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}

		val := ContextValues{Id: id, Role: role}

		ctxWithId := context.WithValue(r.Context(), ValuesFromContext, val)
		rWithId := r.WithContext(ctxWithId)

		next.ServeHTTP(w, rWithId)
	}
}

func CheckUserLevelPermissions(next http.Handler, useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		tokenModel := models.AuthToken{Secret: []byte(token)}

		id, err := useCase.Authorization(&tokenModel, usecase.UserRole)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}

		val := ContextValues{Id: id, Role: usecase.UserRole}

		ctxWithId := context.WithValue(r.Context(), ValuesFromContext, val)
		rWithId := r.WithContext(ctxWithId)

		next.ServeHTTP(w, rWithId)
	}
}

func CheckMusicianLevelPermissions(next http.Handler, useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		tokenModel := models.AuthToken{Secret: []byte(token)}

		id, err := useCase.Authorization(&tokenModel, usecase.MusicianRole)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}

		val := ContextValues{Id: id, Role: usecase.MusicianRole}

		ctxWithId := context.WithValue(r.Context(), ValuesFromContext, val)
		rWithId := r.WithContext(ctxWithId)

		next.ServeHTTP(w, rWithId)
	}
}

func CheckAdminLevelPermissions(next http.Handler, useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			render.JSON(w, r, response.Error("Error in parsing token"))
			render.Status(r, http.StatusBadRequest)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		tokenModel := models.AuthToken{Secret: []byte(token)}
		id, err := useCase.Authorization(&tokenModel, usecase.AdminRole)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}

		val := ContextValues{Id: id, Role: usecase.AdminRole}

		ctxWithId := context.WithValue(r.Context(), ValuesFromContext, val)
		rWithId := r.WithContext(ctxWithId)

		next.ServeHTTP(w, rWithId)
	}
}
