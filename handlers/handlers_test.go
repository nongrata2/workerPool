package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool_AddWorker(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	testTable := []struct {
		name              string
		initialWorkers    int
		workersToAdd      int
		expectedWorkerIDs []int
	}{
		{
			name:              "Add one worker to empty pool",
			initialWorkers:    0,
			workersToAdd:      1,
			expectedWorkerIDs: []int{1},
		},
		{
			name:              "Add two workers to empty pool",
			initialWorkers:    0,
			workersToAdd:      2,
			expectedWorkerIDs: []int{1, 2},
		},
		{
			name:              "Add one worker to a pool with existing workers",
			initialWorkers:    1,
			workersToAdd:      1,
			expectedWorkerIDs: []int{1, 2},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPool(log)
			for range tc.initialWorkers {
				p.AddWorker()
			}
			for range tc.workersToAdd {
				p.AddWorker()
			}

			assert.Len(t, p.Workers, tc.initialWorkers+tc.workersToAdd)

			actualIDs := make([]int, 0, len(p.Workers))
			for id := range p.Workers {
				actualIDs = append(actualIDs, id)
			}
			assert.ElementsMatch(t, tc.expectedWorkerIDs, actualIDs)

			p.Stop()
		})
	}
}

func TestPool_AddWorkerByID(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	testTable := []struct {
		name                string
		workerID            int
		addAgainID          int
		expectedErr         error
		expectedWorkerCount int
		expectedWorkerIDs   []int
	}{
		{
			name:                "Add new worker by ID",
			workerID:            5,
			expectedErr:         nil,
			expectedWorkerCount: 1,
			expectedWorkerIDs:   []int{5},
		},
		{
			name:                "Add another new worker by ID",
			workerID:            10,
			expectedErr:         nil,
			expectedWorkerCount: 1,
			expectedWorkerIDs:   []int{10},
		},
		{
			name:                "Add existing worker by ID",
			workerID:            3,
			addAgainID:          3,
			expectedErr:         fmt.Errorf("воркер с ID %d уже существует", 3),
			expectedWorkerCount: 1,
			expectedWorkerIDs:   []int{3},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPool(log)
			defer p.Stop()

			err := p.AddWorkerByID(tc.workerID)
			assert.NoError(t, err)

			if tc.addAgainID != 0 {
				err = p.AddWorkerByID(tc.addAgainID)
				if tc.expectedErr != nil {
					assert.EqualError(t, err, tc.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			}

			assert.Len(t, p.Workers, tc.expectedWorkerCount)
			actualIDs := make([]int, 0, len(p.Workers))
			for id := range p.Workers {
				actualIDs = append(actualIDs, id)
			}
			assert.ElementsMatch(t, tc.expectedWorkerIDs, actualIDs)
		})
	}
}

func TestPool_RemoveLastWorker(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	testTable := []struct {
		name              string
		initialWorkers    []int
		workersToRemove   int
		expectedWorkerIDs []int
	}{
		{
			name:              "Remove one worker from two",
			initialWorkers:    []int{1, 2},
			workersToRemove:   1,
			expectedWorkerIDs: []int{1},
		},
		{
			name:              "Remove all workers",
			initialWorkers:    []int{1, 2, 3},
			workersToRemove:   3,
			expectedWorkerIDs: []int{},
		},
		{
			name:              "Remove from empty pool",
			initialWorkers:    []int{},
			workersToRemove:   1,
			expectedWorkerIDs: []int{},
		},
		{
			name:              "Remove one worker from three (non-sequential IDs)",
			initialWorkers:    []int{1, 5, 10},
			workersToRemove:   1,
			expectedWorkerIDs: []int{1, 5},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPool(log)
			for _, id := range tc.initialWorkers {
				_ = p.AddWorkerByID(id)
			}
			defer p.Stop()

			for i := 0; i < tc.workersToRemove; i++ {
				p.RemoveLastWorker()
				time.Sleep(10 * time.Millisecond)
			}

			assert.Len(t, p.Workers, len(tc.expectedWorkerIDs))
			actualIDs := make([]int, 0, len(p.Workers))
			for id := range p.Workers {
				actualIDs = append(actualIDs, id)
			}
			assert.ElementsMatch(t, tc.expectedWorkerIDs, actualIDs)
		})
	}
}

func TestPool_RemoveWorkerByID(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	testTable := []struct {
		name                string
		initialWorkers      []int
		removeID            int
		expectedWorkerCount int
		expectedWorkerIDs   []int
	}{
		{
			name:                "Remove existing worker",
			initialWorkers:      []int{1, 2, 3},
			removeID:            2,
			expectedWorkerCount: 2,
			expectedWorkerIDs:   []int{1, 3},
		},
		{
			name:                "Remove non-existent worker",
			initialWorkers:      []int{1, 2},
			removeID:            99,
			expectedWorkerCount: 2,
			expectedWorkerIDs:   []int{1, 2},
		},
		{
			name:                "Remove the only worker",
			initialWorkers:      []int{5},
			removeID:            5,
			expectedWorkerCount: 0,
			expectedWorkerIDs:   []int{},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPool(log)
			for _, id := range tc.initialWorkers {
				_ = p.AddWorkerByID(id)
			}
			defer p.Stop()

			p.RemoveWorkerByID(tc.removeID)
			time.Sleep(10 * time.Millisecond)

			assert.Len(t, p.Workers, tc.expectedWorkerCount)
			actualIDs := make([]int, 0, len(p.Workers))
			for id := range p.Workers {
				actualIDs = append(actualIDs, id)
			}
			assert.ElementsMatch(t, tc.expectedWorkerIDs, actualIDs)
		})
	}
}

func TestPool_Stop(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Stop an active pool", func(t *testing.T) {
		p := NewPool(log)
		p.AddWorker()
		p.AddWorker()
		assert.Len(t, p.Workers, 2)

		p.Stop()

		assert.Len(t, p.Workers, 0)
		time.Sleep(50 * time.Millisecond)
		select {
		case <-p.JobQueue:
			t.Fatal("JobQueue should be empty after Stop")
		default:
		}
	})

	t.Run("Stop an empty pool", func(t *testing.T) {
		p := NewPool(log)
		assert.Len(t, p.Workers, 0)

		p.Stop()

		assert.Len(t, p.Workers, 0)
	})
}

func TestPool_JobProcessing(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("Jobs are processed by workers", func(t *testing.T) {
		p := NewPool(log)
		p.AddWorker()
		defer p.Stop()

		jobCount := 3
		var wg sync.WaitGroup
		wg.Add(jobCount)
		for i := 0; i < jobCount; i++ {
			job := fmt.Sprintf("job-%d", i)
			go func(job string) {
				p.JobQueue <- job
				wg.Done()
			}(job)
		}

		wg.Wait()
		time.Sleep(time.Duration(jobCount)*2*time.Second + 500*time.Millisecond)
	})
}
