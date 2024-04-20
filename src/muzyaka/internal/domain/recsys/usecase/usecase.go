package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/recsys/remote"
	"src/internal/domain/track/repository"
	"src/internal/models"
)

type RecSysUseCase interface {
	GetSameTracks(id uint64) ([]*models.TrackMeta, error)
}

type usecase struct {
	recsProvider remote.RecSysProvider
	trackRep     repository.TrackRepository
}

func NewRecSysUseCase(recs remote.RecSysProvider, repo repository.TrackRepository) RecSysUseCase {
	return &usecase{
		recsProvider: recs,
		trackRep:     repo,
	}
}

func (u *usecase) GetSameTracks(id uint64) ([]*models.TrackMeta, error) {
	trackIds, err := u.recsProvider.GetRecs(id)
	if err != nil {
		return nil, errors.Wrap(err, "recsys.usecase.GetSameTracks error while GetRecs call")
	}

	var tracks []*models.TrackMeta

	for _, v := range trackIds {
		track, err := u.trackRep.GetTrack(v)
		if err != nil {
			return nil, errors.Wrap(err, "recsys.usecase.GetSameTracks error while trackRep call")
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}
