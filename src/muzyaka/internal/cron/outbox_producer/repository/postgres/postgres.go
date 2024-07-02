package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/cron/outbox_producer/repository"
	"src/internal/models/dao"
)

type outboxRepo struct {
	db *gorm.DB
}

func NewOutboxRepo(db *gorm.DB) repository.OutboxRepository {
	return &outboxRepo{db: db}
}

func (o outboxRepo) GetWaitingEvents() ([]*dao.Outbox, error) {
	var events []*dao.Outbox

	tx := o.db.Limit(dao.MaxLimit).Find(&events, "sent = false").Order("id")

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "GetWaitingEvents database error (table outbox)")
	}

	return events, nil
}

func (o outboxRepo) MarkEventsSent(events []*dao.Outbox) error {
	if len(events) > 0 {
		tx := o.db.Model(&events).Update("sent", "true")
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "MarkEventsSent database error (table outbox)")
		}
	}
	return nil
}
