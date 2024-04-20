package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/track/repository"
	"src/internal/models"
)

type TrackUseCase interface {
	UpdatedTrack(track *models.TrackMeta) error
	GetTrack(id uint64) (*models.TrackMeta, error)
}

type usecase struct {
	trackRep repository.TrackRepository
}

func NewTrackUseCase(rep repository.TrackRepository) TrackUseCase {
	return &usecase{trackRep: rep}
}

func (u *usecase) GetTrack(id uint64) (*models.TrackMeta, error) {
	res, err := u.trackRep.GetTrack(id)

	if err != nil {
		return nil, errors.Wrap(err, "track.usecase.GetTrack error while get")
	}

	return res, nil
}

func (u *usecase) UpdatedTrack(track *models.TrackMeta) error {
	err := u.trackRep.UpdateTrack(track)

	if err != nil {
		return errors.Wrap(err, "track.usecase.UpdatedTrack error while update")
	}

	return nil
}
