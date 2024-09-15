package delivery

import (
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/auth/usecase"
	"src/internal/lib/api/response"
	"src/internal/models/dto"
)

// @Summary SignIn
// @Tags auth
// @Description sign in
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param input body dto.SignIn true "login and password"
// @Success 200 {object} dto.SignInResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/auth/sign-in [post]
func SignIn(useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.SignIn
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		token, err := useCase.SignIn(req.Email, req.Password)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.SignInResponse{
			Token: string(token.Secret),
		})
	}
}

// @Summary SignUp
// @Tags auth
// @Description sign up
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param input body dto.SignUp true "user info"
// @Success 200 {object} dto.SignUpResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/auth/sign-up/user [post]
func SignUp(useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.SignUp
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		token, err := useCase.SignUp(dto.ToModelUserWithRole(&req.UserInfo, 0, usecase.UserRole))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.SignUpResponse{
			Token: string(token.Secret),
		})
	}
}

// @Summary SignUpAdmin
// @Tags auth
// @Description sign up admin
// @ID sign-up-admin
// @Accept  json
// @Produce  json
// @Param input body dto.SignUp true "user info"
// @Success 200 {object} dto.SignUpResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/auth/sign-up/admin [post]
func SignUpAdmin(useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.SignUp
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		token, err := useCase.SignUp(dto.ToModelUserWithRole(&req.UserInfo, 0, usecase.AdminRole))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.SignUpResponse{
			Token: string(token.Secret),
		})
	}
}

// @Summary SignUpMusician
// @Tags auth
// @Description sign up musician
// @ID sign-up-musician
// @Accept  json
// @Produce  json
// @Param input body dto.SignUpMusician true "user and musician info"
// @Success 200 {object} dto.SignUpResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/auth/sign-up/musician [post]
func SignUpMusician(useCase usecase.AuthUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.SignUpMusician
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		token, err := useCase.SignUpMusician(
			dto.ToModelUserWithRole(&req.UserInfo, 0, usecase.MusicianRole),
			dto.ToModelMusicianWithoutId(&req.MusicianWithoutId, 0))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		render.JSON(w, r, dto.SignUpResponse{
			Token: string(token.Secret),
		})
	}
}
