package leet_code

import "math"

func numSquares(n int) int {
	queue := make([]int, 0)
	visited := make(map[int]byte, 0)

	queue = append(queue, 0)
	visited[0] = '0'

	level := 0

	for len(queue) != 0 {
		size := len(queue)
		level++

		for i := 0; i < size; i++ {
			digit := queue[0]
			queue = queue[1:]

			for j := 1; j <= n; j++ {
				nodeValue := digit + j * j

				if nodeValue == n {
					return level
				} else if nodeValue > n {
					break
				}

				if _, ok := visited[nodeValue]; !ok {
					queue = append(queue, nodeValue)
					visited[nodeValue] = '0'
				}
			}
		}
	}

	return level
}

func numSquaresDp(n int) int {
	dp := make([]int, n + 1)
	dp[0] = 0

	for i := 1; i < n; i++ {
		dp[i] = i

		for j := 1; j * j <= i; j++ {
			dp[i] = int(math.Min(float64(dp[i]), float64(dp[i - j * j] + 1)))
		}
	}

	return dp[n]
}