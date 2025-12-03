package task1

// 加一
func PlusOne(digits []int) []int {
	ll := len(digits)
	digits[ll-1]++
	if digits[ll-1] < 10 {
		return digits
	}
	for i := ll - 1; i >= 0; i-- {
		if digits[i] == 10 {
			digits[i] = 0
		} else {
			break
		}
		if i == 0 {
			digits = append([]int{1}, digits...)
		} else {
			digits[i-1]++
		}
	}
	return digits
}
