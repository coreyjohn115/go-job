package main

import (
	task2 "job/task-2"
)

func main() {
	tasks := []*task2.Task{
		task2.NewTask("Task 1"),
		task2.NewTask("Task 2"),
		task2.NewTask("Task 3"),
	}
	task2.ExecuteTask(tasks)
}
