package outbox_producer

import (
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/models"
	"src/internal/models/dao"
)

// Идея в том, что сбои после загрузки в кафку уже обрабатываются на стороне кафке, те вроде можно
// обойтись без обратной связи, если что по-мужски поставить время жизни сообщения неделю

type OutboxProducer struct {
	producer    sarama.SyncProducer
	db          *gorm.DB
	outboxTable string
}

func New(producer sarama.SyncProducer, db *gorm.DB, outboxTable string) *OutboxProducer {
	return &OutboxProducer{
		producer:    producer,
		db:          db,
		outboxTable: outboxTable,
	}
}

func (op *OutboxProducer) ProduceMessages() error {
	err := op.db.Transaction(func(tx *gorm.DB) error {
		var events []dao.Outbox
		err := op.db.Limit(models.MaxLimit).Find(&events, "sent = false").Error
		if err != nil {
			return errors.Wrap(tx.Error, "error in Kafka producer")
		}

		// TODO Дописать.
		return nil
	})
	return err
}
