package remote

type RecSysProvider interface {
	GetRecs(id uint64) ([]uint64, error)
}

// TODO: Implement GetRecs, maybe add Kafka
