package usecase

import (
	"github.com/pkg/errors"
	usecase2 "src/internal/domain/auth/usecase"
	repository2 "src/internal/domain/track/repository"
	"src/internal/domain/user/repository"
	"src/internal/models"
)

type UserUseCase interface {
	UpdateUser(user *models.User) error
	GetUser(id uint64) (*models.User, error)
	AddUser(user *models.User) (uint64, error)
	DeleteUser(id uint64) error

	LikeTrack(userId uint64, trackId uint64) error
	DislikeTrack(userId uint64, trackId uint64) error
	GetAllLikedTracks(userId uint64) ([]*models.TrackMeta, error)
	IsTrackLiked(userId uint64, trackId uint64) (bool, error)
}

type usecase struct {
	userRep   repository.UserRepository
	trackRep  repository2.TrackRepository
	encryptor usecase2.Encryptor
}

func NewUserUseCase(rep repository.UserRepository, trackRep repository2.TrackRepository, encryptor usecase2.Encryptor) UserUseCase {
	return &usecase{userRep: rep, trackRep: trackRep, encryptor: encryptor}
}

func (u *usecase) IsTrackLiked(userId uint64, trackId uint64) (bool, error) {
	ans, err := u.userRep.IsTrackLiked(userId, trackId)
	if err != nil {
		return false, errors.Wrap(err, "user.usecase.IsTrackLiked error while check")
	}

	return ans, nil
}

func (u *usecase) GetAllLikedTracks(userId uint64) ([]*models.TrackMeta, error) {
	trackIds, err := u.userRep.GetAllLikedTracks(userId)
	if err != nil {
		return nil, errors.Wrap(err, "user.usecase.GetAllLikedTracks error while get")
	}

	var trackMeta []*models.TrackMeta
	for _, v := range trackIds {
		track, err := u.trackRep.GetTrack(v)
		if err != nil {
			return nil, errors.Wrap(err, "user.usecase.GetAllLikedTracks error while get")
		}

		trackMeta = append(trackMeta, track)
	}

	return trackMeta, nil
}

func (u *usecase) UpdateUser(user *models.User) error {
	if user.Password == "" {
		return models.ErrInvalidPassword
	}

	encPassword, err := u.encryptor.EncodePassword([]byte(user.Password))
	if err != nil {
		return errors.Wrap(err, "user.usecase.UpdateUser encode error")
	}

	temp := user
	temp.Password = string(encPassword)

	err = u.userRep.UpdateUser(temp)

	if err != nil {
		return errors.Wrap(err, "user.usecase.UpdateUser error while update")
	}

	user.Password = ""

	return nil
}

func (u *usecase) LikeTrack(userId uint64, trackId uint64) error {
	err := u.userRep.LikeTrack(userId, trackId)

	if err != nil {
		return errors.Wrap(err, "user.usecase.LikeTrack error while add")
	}

	return nil
}

func (u *usecase) DislikeTrack(userId uint64, trackId uint64) error {
	err := u.userRep.DislikeTrack(userId, trackId)

	if err != nil {
		return errors.Wrap(err, "user.usecase.DislikeTrack error while delete")
	}

	return nil
}

func (u *usecase) GetUser(id uint64) (*models.User, error) {
	res, err := u.userRep.GetUser(id)

	if err != nil {
		return nil, errors.Wrap(err, "user.usecase.GetUser error while get")
	}

	return res, nil
}

func (u *usecase) AddUser(user *models.User) (uint64, error) {
	id, err := u.userRep.AddUser(user)

	if err != nil {
		return 0, errors.Wrap(err, "user.usecase.AddUser error while add")
	}

	return id, nil
}

func (u usecase) DeleteUser(id uint64) error {
	err := u.userRep.DeleteUser(id)

	if err != nil {
		return errors.Wrap(err, "user.usecase.DeleteUser error while delete")
	}

	return nil
}
