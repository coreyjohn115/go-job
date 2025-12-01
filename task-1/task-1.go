package task1

// 两数之和
func TwoSum(nums []int, target int) []int {
	mm := make(map[int]int)
	for i := range nums {
		need := target - nums[i]
		if _, ok := mm[need]; ok {
			return []int{mm[need], i}
		}
		mm[nums[i]] = i
	}
	return nil
}

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
