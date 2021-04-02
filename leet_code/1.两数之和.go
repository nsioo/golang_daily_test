package leet_code
//给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 的那 两个 整数，并返回它们的数组下标。
//
// 你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。
//
// 你可以按任意顺序返回答案。
//
//
//
// 示例 1：
//
//
//输入：nums = [2,7,11,15], target = 9
//输出：[0,1]
//解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。
//
//
// 示例 2：
//
//
//输入：nums = [3,2,4], target = 6
//输出：[1,2]
//
//
// 示例 3：
//
//
//输入：nums = [3,3], target = 6
//输出：[0,1]
//
//
//
//
// 提示：
//
//
// 2 <= nums.length <= 103
// -109 <= nums[i] <= 109
// -109 <= target <= 109
// 只会存在一个有效答案
//
// Related Topics 数组 哈希表
// 👍 10692 👎 0





// 使用暴力解法，时间复杂度 n^2
func twoSum(nums []int, target int) []int {
	res := make([]int, 2)
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] + nums[j] == target {
				res[0] = i
				res[1] = j
			}
		}
	}
	return res
}

// 使用哈希方法解决，时间复杂度 n*logn
// 初始化一个 map，遍历 nums
// 每次使用 target - nums[i]，得到一个值在 map 中找这个 key
// 如果没找到就把 nums[i]-i 放入 map
// 如果找到就返回 [i, map[target - nums[i]]]
func TwoSumHash(nums []int, target int) []int {
	hashMap := make(map[int]int, 0)
	res := make([]int, 2)

	for i := 0; i < len(nums); i++ {
		if index, ok := hashMap[target - nums[i]]; ok {
			res[0] = index
			res[1] = i
			break
		} else {
			hashMap[nums[i]] = i
		}
	}

	return res
}