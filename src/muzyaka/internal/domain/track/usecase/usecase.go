package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/track/repository"
	"src/internal/models"
)

const MinPageSize = 10
const MaxPageSize = 100

type TrackUseCase interface {
	UpdateTrack(track *models.TrackObject) error
	GetTrack(id uint64) (*models.TrackObject, error)
	// TODO: нужен ли метод ниже? как будто его не используют
	DeleteTrack(trackId uint64) error

	GetTracksByPartName(name string, page int, pageSize int) ([]*models.TrackMeta, error)
}

type usecase struct {
	trackRep   repository.TrackRepository
	storageRep repository.TrackStorage
}

func NewTrackUseCase(rep repository.TrackRepository, storage repository.TrackStorage) TrackUseCase {
	return &usecase{trackRep: rep, storageRep: storage}
}

func (u *usecase) GetTracksByPartName(name string, page int, pageSize int) ([]*models.TrackMeta, error) {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > MaxPageSize:
		pageSize = MaxPageSize
	case pageSize < MinPageSize:
		pageSize = MinPageSize
	}

	offset := (page - 1) * pageSize
	tracks, err := u.trackRep.GetTracksByPartName(name, offset, pageSize)
	if err != nil {
		return nil, errors.Wrap(err, "track.usecase.GetTracksByPartName error while get")
	}

	return tracks, nil
}

func (u *usecase) DeleteTrack(trackId uint64) error {

	trackMeta, err := u.trackRep.GetTrack(trackId)
	if err != nil {
		return errors.Wrap(err, "track.usecase.DeleteTrack error while get")
	}

	err = u.trackRep.DeleteTrack(trackId)
	if err != nil {
		return errors.Wrap(err, "track.usecase.DeleteTrack error while delete")
	}

	err = u.storageRep.DeleteObject(trackMeta)
	if err != nil {
		return errors.Wrap(err, "track.usecase.DeleteTrack error while delete")
	}

	return nil
}

func (u *usecase) GetTrack(id uint64) (*models.TrackObject, error) {
	meta, err := u.trackRep.GetTrack(id)
	if err != nil {
		return nil, errors.Wrap(err, "track.usecase.GetTrack error while get")
	}

	res, err := u.storageRep.LoadObject(meta)
	if err != nil {
		return nil, errors.Wrap(err, "track.usecase.GetTrack error while get")
	}

	return res, nil
}

func (u *usecase) UpdateTrack(track *models.TrackObject) error {
	err := u.storageRep.UploadObject(track)
	if err != nil {
		return errors.Wrap(err, "track.usecase.UpdateTrack error while update")
	}

	err = u.trackRep.UpdateTrack(track.ExtractMeta())
	if err != nil {
		return errors.Wrap(err, "track.usecase.UpdateTrack error while update")
	}

	return nil
}
