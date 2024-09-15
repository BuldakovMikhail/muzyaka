package repository

import "src/internal/models/dao"

type OutboxRepository interface {
	GetWaitingEvents() ([]*dao.Outbox, error)
	MarkEventsSent(events []*dao.Outbox) error
}
