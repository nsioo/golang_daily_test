package main

import "fmt"

func decodeString(s string) string {
	nums := make([]int, 0)
	strs := make([]string, 0)

	num := 0
	str := ""

	for _, curByte := range s {
		// 如果 curByte 为数字，将num * 10 再加 curByte
		if curByte >= '0' && curByte <= '9' {
			num = num * 10 + int(curByte - '0')
		} else if (curByte >= 'a' && curByte <= 'z') || (curByte >= 'A' &&  curByte <= 'Z') {
			str += string(curByte)
		} else if curByte == '[' {
			nums = append(nums, num)
			num = 0
			strs = append(strs, str)
			str = ""
		} else {
			times := nums[len(nums) - 1]
			nums = nums[:len(nums) - 1]

			for i := 0; i < times; i++ {
				strs[len(strs) - 1] += str
			}

			str = strs[len(strs) - 1]
			strs = strs[:len(strs) - 1]
		}
	}

	return str
}

func main() {
	fmt.Println(decodeString("2[a3[c]]"))
}