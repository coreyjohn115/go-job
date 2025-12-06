package task1

import (
	"math"
	"sort"
	"strconv"
)

// 只出现一次的数字
func SingleNumber(nums []int) int {
	// 异或运算，相同为0，不同为1
	// result := 0
	// for _, v := range nums {
	// 	result ^= v
	// }
	// return result
	var m = make(map[int]int)
	for _, v := range nums {
		m[v]++
	}
	for k, v := range m {
		if v == 1 {
			return k
		}
	}
	return -1
}

// 回文数
func IsPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	var l = len(strconv.Itoa(x))
	for i := 0; i < l/2; i++ {
		p1, p2 := int(math.Pow10(i+1)), int(math.Pow10(l-i))
		d1, d2 := int(math.Pow10(i)), int(math.Pow10(l-i-1))
		if x%p1/d1 != x%p2/d2 {
			return false
		}
	}
	return true
}

// 有效的括号
func StrValid(s string) bool {
	m := map[string]string{
		"(": ")",
		"[": "]",
		"{": "}",
	}
	stack := []string{}
	for _, v := range s {
		if _, ok := m[string(v)]; ok {
			stack = append(stack, m[string(v)])
		} else {
			if len(stack) == 0 || stack[len(stack)-1] != string(v) {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

// 最长公共前缀
func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, str := range strs {
		for i := 0; i < len(prefix); i++ {
			if i < len(str) && str[i] != prefix[i] {
				prefix = prefix[:i]
				break
			}
		}
	}
	return prefix
}

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

// 删除有序数组中的重复项
func RemoveDuplicates(nums []int) int {
	ll := len(nums)
	unique := 1
	for i := 1; i < ll; i++ {
		if nums[i] != nums[i-1] {
			nums[unique-1] = nums[i]
			unique++
		}
	}
	return unique
}

// 合并区间
func MergeIntervals(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	merged := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] <= merged[len(merged)-1][1] {
			merged[len(merged)-1][1] = max(merged[len(merged)-1][1], intervals[i][1])
		} else {
			merged = append(merged, intervals[i])
		}
	}
	return merged
}
