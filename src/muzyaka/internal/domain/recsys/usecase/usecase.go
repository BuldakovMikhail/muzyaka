package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/recsys/recsys_client"
	"src/internal/domain/track/repository"
	"src/internal/models"
)

type RecSysUseCase interface {
	GetSameTracks(id uint64, page int, pageSize int) ([]*models.TrackMeta, error)
}

type usecase struct {
	recsProvider recsys_client.RecSysProvider
	trackRep     repository.TrackRepository
}

func NewRecSysUseCase(recs recsys_client.RecSysProvider, repo repository.TrackRepository) RecSysUseCase {
	return &usecase{
		recsProvider: recs,
		trackRep:     repo,
	}
}

func (u *usecase) GetSameTracks(id uint64, page int, pageSize int) ([]*models.TrackMeta, error) {
	trackIds, err := u.recsProvider.GetRecs(id, page, pageSize)
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
