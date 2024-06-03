package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/merch/repository"
	"src/internal/models"
)

type MerchUseCase interface {
	GetMerch(id uint64) (*models.Merch, error)
	GetAllMerchForMusician(musicianId uint64) ([]*models.Merch, error)
	UpdateMerch(merch *models.Merch) error
	AddMerch(merch *models.Merch, musicianId uint64) (uint64, error)
	DeleteMerch(id uint64) error
	GetMusicianForMerch(merchId uint64) (uint64, error)

	IsMerchOwned(merchId uint64, musicianId uint64) (bool, error)
}

type usecase struct {
	merchRep repository.MerchRepository
}

func NewMerchUseCase(merchRepository repository.MerchRepository) MerchUseCase {
	return &usecase{merchRep: merchRepository}
}

func (u *usecase) IsMerchOwned(merchId uint64, musicianId uint64) (bool, error) {
	res, err := u.merchRep.IsMerchOwned(merchId, musicianId)

	if err != nil {
		return false, errors.Wrap(err, "album.usecase.IsAlbumOwned error while get")
	}

	return res, nil
}

func (u *usecase) GetMusicianForMerch(merchId uint64) (uint64, error) {
	res, err := u.merchRep.GetMusicianForMerch(merchId)

	if err != nil {
		return 0, errors.Wrap(err, "merch.usecase.GetMusicianForMerch error while get")
	}

	return res, nil
}

func (u *usecase) GetMerch(id uint64) (*models.Merch, error) {
	res, err := u.merchRep.GetMerch(id)

	if err != nil {
		return nil, errors.Wrap(err, "merch.usecase.GetMerch error while get")
	}

	return res, nil
}

func (u *usecase) GetAllMerchForMusician(musicianId uint64) ([]*models.Merch, error) {
	res, err := u.merchRep.GetAllMerchForMusician(musicianId)

	if err != nil {
		return nil, errors.Wrap(err, "merch.usecase.GetAllMerchForMusician error while get")
	}

	return res, nil
}

func (u *usecase) UpdateMerch(merch *models.Merch) error {
	err := u.merchRep.UpdateMerch(merch)

	if err != nil {
		return errors.Wrap(err, "merch.usecase.UpdateMerch error while update")
	}

	return nil
}

func (u *usecase) AddMerch(merch *models.Merch, musicianId uint64) (uint64, error) {
	id, err := u.merchRep.AddMerch(merch, musicianId)

	if err != nil {
		return 0, errors.Wrap(err, "merch.usecase.AddMerch error while add")
	}

	return id, nil
}

func (u *usecase) DeleteMerch(id uint64) error {
	err := u.merchRep.DeleteMerch(id)

	if err != nil {
		return errors.Wrap(err, "merch.usecase.DeleteMerch error while delete")
	}

	return nil
}
