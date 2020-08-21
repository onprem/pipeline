package main

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents current status of a task.
type TaskStatus string

// Various possible task status.
const (
	TaskNotStarted TaskStatus = "notstarted"
	TaskRunning    TaskStatus = "running"
	TaskPaused     TaskStatus = "paused"
	TaskTerminated TaskStatus = "terminated"
	TaskGotError   TaskStatus = "goterror"
	TaskFinished   TaskStatus = "finished"
)

// Task represents a processing task in our system.
type Task struct {
	ID       string
	FilePath string
	State    TaskStatus
	Err      error

	pause     chan struct{}
	resume    chan struct{}
	terminate chan struct{}
	mutex     sync.Mutex
}

// NewTask returns an initialized instance of task.
func NewTask(path string) *Task {
	return &Task{
		ID:        uuid.New().String(),
		FilePath:  path,
		State:     TaskNotStarted,
		pause:     make(chan struct{}),
		resume:    make(chan struct{}),
		terminate: make(chan struct{}),
	}
}

// Run is used to start the task.
func (t *Task) Run() {
	t.update(TaskRunning)
	go t.process()
}

// Pause function pauses a running task.
func (t *Task) Pause() {
	t.update(TaskPaused)
	t.pause <- struct{}{}
}

// Resume function resumes a paused task.
// If task is not paused it doesn't have any effect.
func (t *Task) Resume() {
	t.update(TaskRunning)
	t.resume <- struct{}{}
}

func (t *Task) finish() {
	t.update(TaskFinished)
	log.Printf("[%s] finished\n", t.ID)

	close(t.pause)
	close(t.resume)
	close(t.terminate)
}

func (t *Task) error(err error) {
	t.update(TaskGotError)
	t.Err = err
}

func (t *Task) process() {
	file, err := os.Open(t.FilePath)
	if err != nil {
		t.error(err)
		return
	}
	defer file.Close()

	csvR := csv.NewReader(file)

Out:
	for {
		select {
		case <-t.pause:
			<-t.resume
		default:
			record, err := csvR.Read()
			if err == io.EOF {
				break Out
			}
			processRecord(record)
			log.Printf("[%s] processed: %v\n", t.ID, record)
		}
	}

	t.finish()
}

func processRecord(record []string) {
	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)
}

func (t *Task) update(status TaskStatus) {
	t.mutex.Lock()
	t.State = status
	t.mutex.Unlock()
}
