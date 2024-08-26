package snowflake

const (
	BitsSumError     SnowflakeError = "the sum of bits error, must be 63"
	MachineIdIllegal SnowflakeError = "machine id is illegal"
	TimestampIllegal SnowflakeError = "timestamp is illegal"
	SeqIllegal       SnowflakeError = "sequence is illegal"
)

type SnowflakeError string

func (e SnowflakeError) Error() string {
	return string(e)
}
