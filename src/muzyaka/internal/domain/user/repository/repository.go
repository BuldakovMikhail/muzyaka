package repository

import "src/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type UserRepository interface {
	GetUser(id uint64) (*models.User, error)
	UpdateUser(user *models.User) error
	AddUser(user *models.User) (uint64, error)
	DeleteUser(id uint64) error
	GetUserByEmail(email string) (*models.User, error)

	AddUserWithMusician(musician *models.Musician, user *models.User) (uint64, error)

	LikeTrack(userId uint64, trackId uint64) error
	DislikeTrack(userId uint64, trackId uint64) error
	GetAllLikedTracks(userId uint64) ([]uint64, error)
	IsTrackLiked(userId uint64, trackId uint64) (bool, error)
}
