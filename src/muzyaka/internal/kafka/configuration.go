package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

const (
	DefaultTopic = "sync-default-events"
	AddTopic     = "sync-add-events"
	DeleteTopic  = "sync-delete-events"
	UpdateTopic  = "sync-update-events"
)

func NewProducer(addr string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{fmt.Sprintf(addr)}, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "kafka.NewProducer producer creation error")
	}

	return producer, nil
}
