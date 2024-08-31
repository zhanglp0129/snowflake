package snowflake

type WorkerInterface interface {
	GenerateId() (int64, error)
}
