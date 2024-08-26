package snowflake

import (
	"sync"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	config := SnowFlakeConfig{
		StartTimestamp: time.Now().UnixMilli(),
		TimestampBits:  41,
		MachineIdBits:  10,
		SeqBits:        12,
	}

	_, err := NewWorker(config, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewWorkerInvalidConfig(t *testing.T) {
	config := SnowFlakeConfig{
		StartTimestamp: time.Now().UnixMilli(),
		TimestampBits:  40,
		MachineIdBits:  10,
		SeqBits:        12,
	}

	_, err := NewWorker(config, 1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestNewWorkerInvalidMachineId(t *testing.T) {
	config := SnowFlakeConfig{
		StartTimestamp: time.Now().UnixMilli(),
		TimestampBits:  41,
		MachineIdBits:  10,
		SeqBits:        12,
	}

	_, err := NewWorker(config, -1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	_, err = NewWorker(config, 1024)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGenerateIdConcurrency(t *testing.T) {
	config := SnowFlakeConfig{
		StartTimestamp: time.Now().UnixMilli(),
		TimestampBits:  41,
		MachineIdBits:  10,
		SeqBits:        12,
	}

	worker, err := NewWorker(config, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var wg sync.WaitGroup
	var idSet sync.Map
	numGoroutines := 100
	numIdsPerGoroutine := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIdsPerGoroutine; j++ {
				id, err := worker.GenerateId()
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if _, loaded := idSet.LoadOrStore(id, true); loaded {
					t.Errorf("duplicate id found: %d", id)
				}
			}
		}()
	}

	wg.Wait()
}

func TestGetMachineId(t *testing.T) {
	config := SnowFlakeConfig{
		StartTimestamp: time.Now().UnixMilli(),
		TimestampBits:  41,
		MachineIdBits:  10,
		SeqBits:        12,
	}

	workers := make([]WorkerInterface, 0, 16)
	for i := 0; i < 16; i++ {
		worker, err := NewWorker(config, int64(i))
		if err != nil {
			t.Fatal(err)
		}
		workers = append(workers, worker)
	}

	// 生成id
	for i, worker := range workers {
		for j := 0; j < 100; j++ {
			id, err := worker.GenerateId()
			if err != nil {
				t.Fatal(err)
			}
			machineId, err := GetMachineId(config, id)
			if machineId != int64(i) {
				t.Fatal("机器码错误")
			}
		}
	}
}
