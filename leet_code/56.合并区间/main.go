package main

import (
	"fmt"
	"math"
	"sort"
)

func merge(intervals [][]int) [][]int {
	if len(intervals) == 1 {
		return intervals
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})


	res := make([][]int, 0)
	res = append(res, intervals[0])
	j := 0
	for i := j + 1; i < len(intervals); i++ {
		if intervals[i][0] <= res[j][1] {
			res[j][1] = int(math.Max(float64(res[j][1]), float64(intervals[i][1])))
		} else {
			res = append(res, intervals[i])
			j++
		}
	}

	return res
}

func main() {
	fmt.Println(merge([][]int{{8, 10}, {2, 6}, {1, 3}, {15, 18}}))
}
