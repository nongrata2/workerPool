package handlers

import (
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"
)

type Worker struct {
	ID       int
	JobQueue chan string
	Quit     chan bool
	log      *slog.Logger
}

func NewWorker(id int, jobQueue chan string, log *slog.Logger) *Worker {
	return &Worker{
		ID:       id,
		JobQueue: jobQueue,
		Quit:     make(chan bool),
		log:      log,
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job := <-w.JobQueue:
			w.log.Info("Worker started job",
				slog.Int("worker_id", w.ID),
				slog.String("job", job),
			)
			time.Sleep(2 * time.Second)
			w.log.Info("Worker completed job",
				slog.Int("worker_id", w.ID),
				slog.String("job", job),
			)
		case <-w.Quit:
			w.log.Info("Worker quitting", slog.Int("worker_id", w.ID))
			return
		}
	}
}

type Pool struct {
	Workers  map[int]*Worker
	JobQueue chan string
	wg       sync.WaitGroup
	nextID   int
	mu       sync.Mutex
	log      *slog.Logger
}

func NewPool(log *slog.Logger) *Pool {
	return &Pool{
		Workers:  make(map[int]*Worker),
		JobQueue: make(chan string, 10),
		log:      log,
	}
}

func (p *Pool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	maxID := 1
	for id := range p.Workers {
		if id+1 > maxID {
			maxID = id + 1
		}
	}

	worker := NewWorker(maxID, p.JobQueue, p.log)
	p.Workers[maxID] = worker

	p.wg.Add(1)
	go worker.Start(&p.wg)

	p.log.Info("Worker added", slog.Int("worker_id", maxID))
}

func (p *Pool) AddWorkerByID(id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.Workers[id]; exists {
		p.log.Warn("Failed to add worker: ID already exists", slog.Int("worker_id", id))
		return fmt.Errorf("воркер с ID %d уже существует", id)
	}

	worker := NewWorker(id, p.JobQueue, p.log)
	p.Workers[id] = worker
	p.wg.Add(1)
	go worker.Start(&p.wg)

	if id > p.nextID {
		p.nextID = id
	}

	p.log.Info("Worker added by ID", slog.Int("worker_id", id))
	return nil
}

func (p *Pool) RemoveLastWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Workers) == 0 {
		p.log.Info("No active workers to remove")
		return
	}

	var maxID int
	for id := range p.Workers {
		if id > maxID {
			maxID = id
		}
	}

	p.Workers[maxID].Quit <- true
	delete(p.Workers, maxID)
	p.log.Info("Worker removed", slog.Int("worker_id", maxID))
}

func (p *Pool) RemoveWorkerByID(id int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	worker, exists := p.Workers[id]
	if !exists {
		p.log.Warn("Worker not found", slog.Int("worker_id", id))
		return
	}

	worker.Quit <- true
	delete(p.Workers, id)
	p.log.Info("Worker removed", slog.Int("worker_id", id))
}

func (p *Pool) Stop() {
	p.log.Info("Stopping all workers in the pool")
	for id := range p.Workers {
		p.RemoveWorkerByID(id)
	}
	p.wg.Wait()
	p.log.Info("All workers stopped and pool is clean")
}

func (p *Pool) PrintWorkers() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Workers) == 0 {
		fmt.Println("Нет активных воркеров")
		return
	}

	ids := make([]int, 0, len(p.Workers))
	for id := range p.Workers {
		ids = append(ids, id)
	}

	sort.Ints(ids)

	fmt.Print("Активные воркеры: ")
	for _, id := range ids {
		fmt.Printf("%d ", id)
	}
	fmt.Println()
}

func PrintHelp() {
	fmt.Println(`Доступные команды:
  add              - добавить воркера
  addid [id]       - добавить воркера по ID
  ls               - список воркеров
  remove           - удалить последнего воркера
  removeid [id]    - удалить по ID
  process [task]   - отправить задачу
  help             - вывести информацию о командах
  exit             - выход`)
}
