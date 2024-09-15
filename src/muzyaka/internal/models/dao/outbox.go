package dao

var (
	TypeDelete = "delete"
	TypeAdd    = "add"
	TypeUpdate = "update"
)

type Outbox struct {
	ID         uint64 `gorm:"column:id"`
	EventId    string `gorm:"column:event_id"`
	TrackId    uint64 `gorm:"column:track_id"`
	Source     string `gorm:"column:source"`
	Name       string `gorm:"column:name"`
	GenreRefer uint64 `gorm:"column:genre"`
	Type       string `gorm:"column:type"`
	Sent       bool   `gorm:"column:sent"`
}

func (Outbox) TableName() string {
	return "outbox"
}
