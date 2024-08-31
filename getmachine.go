package snowflake

// GetMachineId 获取机器码
func GetMachineId(c SnowFlakeConfig, id int64) (int64, error) {
	// 检查配置
	sumBits := c.TimestampBits + c.MachineIdBits + c.SeqBits
	if sumBits != 63 {
		return 0, BitsSumError
	}

	return ((((1 << c.MachineIdBits) - 1) << c.SeqBits) & id) >> c.SeqBits, nil
}
