package outbox_producer

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"src/internal/kafka"
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
	saramaMsgs := make([]*sarama.ProducerMessage, 0, dao.MaxLimit)
	eventIds := []string{}

	var events []*dao.Outbox

	err := op.db.Transaction(func(tx *gorm.DB) error {
		err := op.db.Limit(dao.MaxLimit).Find(&events, "sent = false").Error
		if err != nil {
			return errors.Wrap(tx.Error, "error in Kafka producer")
		}

		for _, v := range events {
			// TODO: Нужны ли разделения по топикам? Может ли привести к ошибкам когда сообщения не в том порядке?

			saramaMsgs = append(saramaMsgs, &sarama.ProducerMessage{
				Topic: kafka.DefaultTopic,
				Value: sarama.StringEncoder(
					fmt.Sprintf(
						"{\"event_id\": \"%s\","+
							" \"track_id\": %d,"+
							" \"operation\": %s,"+
							" \"source\": %s,"+
							" \"name\": %s,"+
							" \"genre_id\": %d}",
						v.EventId,
						v.TrackId,
						v.Type,
						v.Source,
						v.Name,
						v.GenreRefer)),
			})

			eventIds = append(eventIds, v.EventId)
		}

		if err := op.producer.SendMessages(saramaMsgs); err != nil {
			return err
		}

		if err := op.db.Model(&events).Update("sent", "true").Error; err != nil {
			return err
		}

		return nil
	})
	return err
}
