package dto

import "src/internal/models"

type Merch struct {
	Id          uint64   `json:"id"`
	Name        string   `json:"name"`
	PhotoFiles  [][]byte `json:"photo_files"`
	Description string   `json:"description"`
	OrderUrl    string   `json:"order_url"`
}

type MerchWithoutId struct {
	Name        string   `json:"name"`
	PhotoFiles  [][]byte `json:"photo_files"`
	Description string   `json:"description"`
	OrderUrl    string   `json:"order_url"`
}

type MerchWithMusician struct {
	Merch
	MusicianId uint64 `json:"musician_id"`
}

type CreateMerchResponse struct {
	Id uint64 `json:"id"`
}

type MerchCollection struct {
	Items []*Merch `json:"items"`
}

func ToModelMerchWithoutId(m *MerchWithoutId, id uint64) *models.Merch {
	return &models.Merch{
		Id:          id,
		Name:        m.Name,
		PhotoFiles:  m.PhotoFiles,
		Description: m.Description,
		OrderUrl:    m.OrderUrl,
	}
}

func ToModelMerch(m *Merch) *models.Merch {
	return &models.Merch{
		Id:          m.Id,
		Name:        m.Name,
		PhotoFiles:  m.PhotoFiles,
		Description: m.Description,
		OrderUrl:    m.OrderUrl,
	}
}

func ToDtoMerchWithMusician(m *models.Merch, musicianId uint64) *MerchWithMusician {
	return &MerchWithMusician{
		Merch: Merch{
			Id:          m.Id,
			Name:        m.Name,
			PhotoFiles:  m.PhotoFiles,
			Description: m.Description,
			OrderUrl:    m.OrderUrl,
		},
		MusicianId: musicianId,
	}
}

func ToDtoMerch(m *models.Merch) *Merch {
	return &Merch{
		Id:          m.Id,
		Name:        m.Name,
		PhotoFiles:  m.PhotoFiles,
		Description: m.Description,
		OrderUrl:    m.OrderUrl,
	}
}
