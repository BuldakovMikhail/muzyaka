package remote

type RecSysProvider interface {
	GetRecs(id uint64) ([]uint64, error)
}

// TODO: Implement GetRecs, maybe add Kafka
// TODO: transact outbox, чтобы объед две записи в кафку и бд, таблица ивентов. Подумать про записи
