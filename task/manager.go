package task

import (
	"context"
	"sync"
)

type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
)

type Result struct {
	ID     string      `json:"id"`
	Status Status      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type Manager struct {
	tasks     map[string]*Result
	taskQueue chan Task
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		tasks:     make(map[string]*Result),
		taskQueue: make(chan Task, 100),
	}
}

func (m *Manager) EnqueueTask(t Task) {
	m.mu.Lock()
	m.tasks[t.ID()] = &Result{
		ID:     t.ID(),
		Status: StatusQueued,
	}
	m.mu.Unlock()

	m.taskQueue <- t
}

func (m *Manager) GetTask(id string) (*Result, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res, ok := m.tasks[id]
	return res, ok
}

func (m *Manager) StartWorkerPool(n int) {
	for i := 0; i < n; i++ {
		go m.worker()
	}
}

func (m *Manager) processTask(t Task) {
	m.mu.Lock()
	m.tasks[t.ID()].Status = StatusRunning
	m.mu.Unlock()

	result, err := t.Run(context.Background())

	m.mu.Lock()
	defer m.mu.Unlock()
	if err != nil {
		m.tasks[t.ID()].Status = StatusFailed
		m.tasks[t.ID()].Error = err.Error()
	} else {
		m.tasks[t.ID()].Status = StatusDone
		m.tasks[t.ID()].Result = result
	}
}

func (m *Manager) worker() {
	for task := range m.taskQueue {
		m.processTask(task)
	}
}
