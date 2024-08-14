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
