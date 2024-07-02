package usecase

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"src/internal/cron/outbox_producer/repository"
	"src/internal/lib/kafka"
	"src/internal/models/dao"
)

type OutboxProducer struct {
	producer   sarama.SyncProducer
	repository repository.OutboxRepository
}

func NewOutboxUseCase(producer sarama.SyncProducer, outboxRepository repository.OutboxRepository) *OutboxProducer {
	return &OutboxProducer{
		producer:   producer,
		repository: outboxRepository,
	}
}

func (op *OutboxProducer) ProduceMessages() error {
	saramaMsgs := make([]*sarama.ProducerMessage, 0, dao.MaxLimit)

	events, err := op.repository.GetWaitingEvents()
	if err != nil {
		return errors.Wrap(err, "outbox_producer.ProduceMessages error from repository")
	}

	for _, v := range events {
		saramaMsgs = append(saramaMsgs, &sarama.ProducerMessage{
			Topic: kafka.DefaultTopic,
			Value: sarama.StringEncoder(
				fmt.Sprintf(
					"{\"event_id\": \"%s\","+
						" \"track_id\": %d,"+
						" \"operation\": \"%s\","+
						" \"source\": \"%s\","+
						" \"name\": \"%s\","+
						" \"genre_id\": %d}",
					v.EventId,
					v.TrackId,
					v.Type,
					v.Source,
					v.Name,
					v.GenreRefer)),
		})
	}

	if err := op.producer.SendMessages(saramaMsgs); err != nil {
		return errors.Wrap(err, "outbox_producer.ProduceMessages error from producer")
	}

	if err := op.repository.MarkEventsSent(events); err != nil {
		return errors.Wrap(err, "outbox_producer.ProduceMessages error from repository")
	}

	return nil
}
