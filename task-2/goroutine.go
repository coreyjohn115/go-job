package task2

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Goroutine11() {
	go func() {
		for v := range 10 {
			if v%2 == 1 {
				fmt.Println("奇数", v)
			}
		}
	}()
	go func() {
		for v := range 10 {
			if v%2 == 0 {
				fmt.Println("偶数", v)
			}
		}
	}()

	time.Sleep(time.Second * 1)
}

const (
	StatePending = iota
	StateRunning
	StateCompleted
	StateFailed
)

type Task struct {
	ID     int
	Status int
	Data   any
}

func NewTask(d any) *Task {
	return &Task{
		ID:     rand.Intn(10000),
		Status: StatePending,
		Data:   d,
	}
}

func (t *Task) Run() error {
	t.Status = StateRunning
	delay := time.Second * time.Duration(rand.Intn(5))
	time.Sleep(delay)
	t.Status = StateCompleted
	fmt.Println("Task: ", t.ID, "completed", t.Data)
	return nil
}

// 执行任务
func ExecuteTask(t []*Task) {
	wg := sync.WaitGroup{}
	for _, v := range t {
		wg.Add(1)
		go func(t *Task) {
			defer wg.Done()
			err := t.Run()
			if err != nil {
				t.Status = StateFailed
				fmt.Println("Task: ", t.ID, "failed", t.Data)
			}
		}(v)
	}
	wg.Wait()
}
