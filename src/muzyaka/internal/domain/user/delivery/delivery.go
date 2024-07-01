package delivery

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"src/internal/domain/auth/middleware"
	usecase3 "src/internal/domain/auth/usecase"
	usecase2 "src/internal/domain/musician/usecase"
	"src/internal/domain/user/usecase"
	"src/internal/lib/api/response"
	"src/internal/models"
	"src/internal/models/dto"
	"strconv"
)

// @Summary UpdateUser
// @Security ApiKeyAuth
// @Tags user
// @Description update user
// @ID update-user
// @Accept  json
// @Produce  json
// @Param input body dto.UserInfo true "user info"
// @Param user_id path int true "user ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id} [put]
func UpdateUser(useCase usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.UserInfo
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := chi.URLParam(r, "user_id")
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
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

		err = useCase.UpdateUser(dto.ToModelUserWithRole(&req, userIDUint, userInfo.Role))
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary GetUser
// @Security ApiKeyAuth
// @Tags user
// @Description get user by ID
// @ID get-user
// @Accept  json
// @Produce  json
// @Param user_id path int true "user ID"
// @Success 200 {object} dto.User
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id} [get]
func GetUser(useCase usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := useCase.GetUser(userIDUint)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, user)
	}
}

// @Summary DeleteUser
// @Security ApiKeyAuth
// @Tags user
// @Description delete user
// @ID delete-user
// @Accept  json
// @Produce  json
// @Param user_id path int true "user ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id} [delete]
func DeleteUser(useCase usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = useCase.DeleteUser(userIDUint)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary GetMe
// @Security ApiKeyAuth
// @Tags user
// @Description get me
// @ID get-me
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.GetMeResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/get-me [get]
func GetMe(musicianUseCase usecase2.MusicianUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp dto.GetMeResponse
		userInfo, isOk := r.Context().Value(middleware.ValuesFromContext).(middleware.ContextValues)
		if !isOk {
			render.JSON(w, r, response.Error(models.ErrInvalidContext.Error()))
			render.Status(r, http.StatusBadRequest)
			return
		}
		resp.UserId = userInfo.Id
		resp.Role = userInfo.Role

		if userInfo.Role == usecase3.MusicianRole {
			musicianId, err := musicianUseCase.GetMusicianIdForUser(userInfo.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, response.Error(err.Error()))
				return
			}

			resp.MusicianId = musicianId
		}

		render.JSON(w, r, resp)
	}
}

// @Summary LikeTrack
// @Security ApiKeyAuth
// @Tags user
// @Description like track
// @ID like-track
// @Accept  json
// @Produce  json
// @Param user_id path int true "user ID"
// @Param input body dto.Like true "liked track"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id}/favorite [post]
func Like(useCase usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req dto.Like
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = useCase.LikeTrack(userIDUint, req.TrackId)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary DislikeTrack
// @Security ApiKeyAuth
// @Tags user
// @Description dislike track
// @ID dislike-track
// @Accept  json
// @Produce  json
// @Param user_id path int true "user ID"
// @Param input body dto.Dislike true "disliked track"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id}/favorite [delete]
func Dislike(useCase usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var req dto.Dislike
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = useCase.DislikeTrack(userIDUint, req.TrackId)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, response.OK())
	}
}

// @Summary GetAllLiked
// @Security ApiKeyAuth
// @Tags user
// @Description get all liked tracks
// @ID get-all-liked-tracks
// @Accept  json
// @Produce  json
// @Param user_id path int true "user ID"
// @Success 200 {object} dto.TracksMetaCollection
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/user/{user_id}/favorite [get]
func GetAllLiked(useCase usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		likedTracks, err := useCase.GetAllLikedTracks(userIDUint)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var dtoLikedTracks []*dto.TrackMeta
		for _, v := range likedTracks {
			dtoLikedTracks = append(dtoLikedTracks, dto.ToDtoTrackMeta(v))
		}

		render.JSON(w, r, dto.TracksMetaCollection{Tracks: dtoLikedTracks})
	}
}
