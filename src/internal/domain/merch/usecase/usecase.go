package usecase

import (
	"github.com/pkg/errors"
	"src/internal/domain/merch/repository"
	"src/internal/models"
)

type MerchUseCase interface {
	GetMerch(id uint64) (*models.Merch, error)
	UpdateMerch(merch *models.Merch) error
	AddMerch(merch *models.Merch) (uint64, error)
	DeleteMerch(id uint64) error
}

type usecase struct {
	merchRep repository.MerchRepository
}

func (u *usecase) GetMerch(id uint64) (*models.Merch, error) {
	res, err := u.merchRep.GetMerch(id)

	if err != nil {
		return nil, errors.Wrap(err, "merch.usecase.GetMerch error while get")
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

func (u *usecase) AddMerch(merch *models.Merch) (uint64, error) {
	id, err := u.merchRep.AddMerch(merch)

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
