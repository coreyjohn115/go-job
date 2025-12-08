package task2

import "fmt"

func spawn_num(ch chan int) {
	for v := range 10 {
		ch <- v
	}
	defer close(ch)
}

func receive_num(ch chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

func RunChannel() {
	ch := make(chan int)
	go spawn_num(ch)
	receive_num(ch)
}
