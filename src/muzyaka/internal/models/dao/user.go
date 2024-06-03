package dao

import "src/internal/models"

type User struct {
	ID       uint64 `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Role     string `gorm:"column:role"`
}

func (User) TableName() string {
	return "users"
}

type UserTrack struct {
	TrackId uint64 `gorm:"column:track_id"`
	UserId  uint64 `gorm:"column:user_id"`
}

func (UserTrack) TableName() string {
	return "user_track"
}

type UserMusician struct {
	UserId     uint64 `gorm:"user_id"`
	MusicianId uint64 `gorm:"musician_id"`
}

func (UserMusician) TableName() string {
	return "users_musicians"
}

func ToModelUser(user *User) *models.User {
	return &models.User{
		Id:       user.ID,
		Name:     user.Name,
		Password: user.Password,
		Role:     user.Role,
		Email:    user.Email,
	}
}

func ToPostgresUser(user *models.User) *User {
	return &User{
		ID:       user.Id,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}
