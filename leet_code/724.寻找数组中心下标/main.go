package main

import "fmt"

// 1. 暴力解法
func pivotIndex(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	for index := 0; index < len(nums); index++ {
		left := getArraySum(nums, 0, index)
		right := getArraySum(nums, index + 1, len(nums))

		if left - right == 0 {
			return index
		}
	}

	return -1
}

func getArraySum(nums []int, left, right int) int {
	if left == right || left > right {
		return 0
	}

	res := 0
	for index := left; index < right; index++ {
		res += nums[index]
	}

	return res
}

// 该方法可以试着画一画
func pivotIndexBetter(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}

	if sum - nums[0] == 0 {
		return 0
	}

	curSum := 0
	leftWithoutIndex := 0
	for index := 0; index < len(nums); index++ {
		curSum += nums[index]

		if index != 0 {
			leftWithoutIndex = curSum - nums[index]
		}

		if leftWithoutIndex == sum - curSum {
			return index
		}
	}

	return -1
}

func main() {
	fmt.Println(pivotIndexBetter([]int{1, 7, 3, 6, 5, 6}))
}