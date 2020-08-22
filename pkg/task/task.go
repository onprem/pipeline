package task

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Status represents current status of a task.
type Status string

// Various possible task status.
const (
	TaskNotStarted Status = "not-started"
	TaskRunning    Status = "running"
	TaskPaused     Status = "paused"
	TaskTerminated Status = "terminated"
	TaskGotError   Status = "got-error"
	TaskFinished   Status = "finished"
)

// Task represents a processing task in our system.
type Task struct {
	ID       string
	FilePath string
	State    Status
	Err      error

	pause     chan struct{}
	resume    chan struct{}
	terminate chan struct{}
	mutex     sync.Mutex
}

// NewTask returns an initialized instance of task.
func NewTask(id, path string) *Task {
	return &Task{
		ID:        id,
		FilePath:  path,
		State:     TaskNotStarted,
		pause:     make(chan struct{}),
		resume:    make(chan struct{}),
		terminate: make(chan struct{}),
	}
}

// Run is used to start the task.
// No effect if task is not in not started status.
func (t *Task) Run() {
	if t.State != TaskNotStarted {
		return
	}

	t.update(TaskRunning)
	go t.process()
	log.Printf("[%s] running\n", t.ID)
}

// Pause function pauses a running task.
// If task is not running it doesn't have any effect.
func (t *Task) Pause() {
	if t.State != TaskRunning {
		return
	}

	t.update(TaskPaused)
	t.pause <- struct{}{}
	log.Printf("[%s] paused\n", t.ID)
}

// Resume function resumes a paused task.
// If task is not paused it doesn't have any effect.
func (t *Task) Resume() {
	if t.State != TaskPaused {
		return
	}

	t.update(TaskRunning)
	t.resume <- struct{}{}
	log.Printf("[%s] resumed\n", t.ID)
}

// Terminate will kill the running/paused task.
// Doesn't have any effect on already finished/terminated tasks.
func (t *Task) Terminate() {
	if !(t.State == TaskRunning || t.State == TaskPaused) {
		return
	}

	t.terminate <- struct{}{}
}

func (t *Task) finish() {
	t.update(TaskFinished)
	t.cleanup()
	log.Printf("[%s] finished\n", t.ID)
}

func (t *Task) kill() {
	t.update(TaskTerminated)
	t.cleanup()
	log.Printf("[%s] terminated\n", t.ID)
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
		case <-t.terminate:
			t.kill()
			return
		case <-t.pause:
			select {
			case <-t.resume:
			case <-t.terminate:
				t.kill()
				return
			}
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

func (t *Task) update(status Status) {
	t.mutex.Lock()
	t.State = status
	t.mutex.Unlock()
}

func (t *Task) cleanup() {
	close(t.pause)
	close(t.resume)
	close(t.terminate)
}
