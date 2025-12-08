package task2

import "fmt"

func PointAdd(a *int) {
	*a = *a + 10
}

func Double(arr []int) {
	for i := range arr {
		arr[i] = arr[i] * 2
	}
}

func RunPoint() {
	var a = 10
	PointAdd(&a)
	fmt.Println("a: ", a)

	var arr = []int{1, 2, 3, 4, 5}
	Double(arr)
	fmt.Println("arr: ", arr)
}
