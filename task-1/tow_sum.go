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
