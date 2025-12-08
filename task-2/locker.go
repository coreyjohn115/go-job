package task2

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func mutexAdd(m *sync.Mutex, count *int) {
	m.Lock()
	defer m.Unlock()
	for range 1000 {
		*count++
	}
}

func atomicAdd(count *int32) {
	for range 1000 {
		atomic.AddInt32(count, 1)
	}
}

func RunMutex() {
	var m = sync.Mutex{}
	var count = 0
	for range 10 {
		go mutexAdd(&m, &count)
	}
	time.Sleep(time.Millisecond * 500)
	fmt.Println("mutex count: ", count)

	var count32 = int32(0)
	for range 10 {
		go atomicAdd(&count32)
	}
	time.Sleep(time.Millisecond * 500)
	fmt.Println("atomic count: ", count32)
}
