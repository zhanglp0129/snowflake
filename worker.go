package snowflake

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Worker 工作节点
type Worker struct {
	mtx             sync.Mutex
	timestamp       int64 // 生成id的时间戳，从startTimestamp开始
	timestampMax    int64
	timestampOffset uint8
	machineId       int64
	machineIdOffset uint8
	seq             int64
	seqMax          int64
	seqOffset       uint8
}

// NewWorker 创建一个雪花算法的工作节点
func NewWorker(c SnowFlakeConfig, machineId int64) (*Worker, error) {
	// 检查配置
	sumBits := c.TimestampBits + c.MachineIdBits + c.SeqBits
	if sumBits != 63 {
		return nil, errors.New(fmt.Sprintf("the sum of bits is %d, not 63", sumBits))
	}

	// 检查机器码
	var machineMax int64 = (1 << c.MachineIdBits) - 1
	if machineId < 0 || machineId > machineMax {
		return nil, errors.New(fmt.Sprintf("machine id %d is illegal", machineId))
	}

	return &Worker{
		mtx:             sync.Mutex{},
		timestamp:       time.Now().UnixMilli() - c.StartTimestamp,
		timestampMax:    (1 << c.TimestampBits) - 1,
		timestampOffset: c.SeqBits + c.MachineIdBits,
		machineId:       machineId,
		machineIdOffset: c.SeqBits,
		seq:             0,
		seqMax:          (1 << c.SeqBits) - 1,
		seqOffset:       0,
	}, nil
}

func (w *Worker) getId() (int64, error) {
	// 先校验参数
	if w.timestamp < 0 || w.timestamp > w.timestampMax {
		return 0, errors.New(fmt.Sprintf("timestamp %d is illegal", w.timestamp))
	}
	if w.seq < 0 || w.seq > w.seqMax {
		return 0, errors.New(fmt.Sprintf("sequence %d is illegal", w.seq))
	}

	// 生成id
	var id int64 = 0
	id |= w.timestamp << w.timestampOffset
	id |= w.machineId << w.machineIdOffset
	id |= w.seq << w.seqOffset

	return id, nil
}

func (w *Worker) generateId() (int64, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.seq == w.seqMax {
		w.seq = 0
		w.timestamp++
	} else {
		w.seq++
	}

	return w.getId()
}

// GenerateId 生成id
func (w *Worker) GenerateId() (int64, error) {
	for i := 0; i < 3; i++ {
		id, err := w.generateId()
		if err == nil {
			return id, nil
		}
		// 出现错误，等待1-5毫秒后重试
		wait(1, 5)
	}
	return w.generateId()
}

// 随机等待一个时间，单位为毫秒
func wait(from, to int64) {
	waitMilli := rand.Int63()%(to-from+1) + from
	time.Sleep(time.Duration(waitMilli) * time.Millisecond)
}
