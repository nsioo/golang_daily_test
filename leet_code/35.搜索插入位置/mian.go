package main

import "fmt"

func searchInsert(nums []int, target int) int {
	if len(nums) == 0 {
		return 0
	}

	left := 0
	right := len(nums)

	for right > left {
		mid := (right + left) / 2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			right = mid

		} else {
			left = mid + 1

		}
	}
	return left
}

func main() {
	fmt.Println(searchInsert([]int{1, 3, 5, 6}, 4))
}
