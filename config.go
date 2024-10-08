package snowflake

import "time"

// SnowFlakeConfig 雪花
type SnowFlakeConfig struct {
	// 起始时间戳，单位为毫秒，默认为0，即1970-01-01 00:00:00.000
	StartTimestamp int64
	// 时间戳位数，默认为41
	TimestampBits uint8
	// 机器码位数，默认为10
	MachineIdBits uint8
	// 序列号位数，默认为12
	SeqBits uint8
}

var (
	// DefaultConfig 雪花算法默认配置
	DefaultConfig = SnowFlakeConfig{
		StartTimestamp: 0,
		TimestampBits:  41,
		MachineIdBits:  10,
		SeqBits:        12,
	}
)

// NewDefaultConfigWithStartTime 创建一个雪花算法默认配置，并指定起始时间
func NewDefaultConfigWithStartTime(startTime time.Time) SnowFlakeConfig {
	cfg := DefaultConfig
	cfg.SetStartTime(startTime)
	return cfg
}

// SetStartTime 设置起始时间
func (c *SnowFlakeConfig) SetStartTime(startTime time.Time) {
	c.StartTimestamp = startTime.UnixMilli()
}
