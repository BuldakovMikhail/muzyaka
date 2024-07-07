package dto

import "src/internal/models"

type UserInfo struct {
	Name     string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Like struct {
	TrackId uint64 `json:"track_id"`
}

type Dislike struct {
	TrackId uint64 `json:"track_id"`
}

type IsLikedResponse struct {
	IsLiked bool `json:"is_liked"`
}

type User struct {
	UserInfo
	Role string `json:"role"`
}

type CreateUserResponse struct {
	Id uint64 `json:"id"`
}

func ToModelUser(u *User, id uint64) *models.User {
	return &models.User{
		Id:       id,
		Name:     u.Name,
		Password: u.Password,
		Role:     u.Role,
		Email:    u.Email,
	}
}

func ToModelUserWithRole(u *UserInfo, id uint64, role string) *models.User {
	return &models.User{
		Id:       id,
		Name:     u.Name,
		Password: u.Password,
		Role:     role,
		Email:    u.Email,
	}
}

func ToDtoUser(u *models.User) *User {
	return &User{
		UserInfo: UserInfo{
			Name:     u.Name,
			Password: u.Password,
			Email:    u.Email,
		},
		Role: u.Role,
	}
}
