package task2

func PointAdd(a *int) {
	*a = *a + 10
}

func Double(arr []int) {
	for i := range arr {
		arr[i] = arr[i] * 2
	}
}
