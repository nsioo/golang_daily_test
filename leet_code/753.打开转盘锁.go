package leet_code

func OpenLock(deadends []string, target string) int {
	deadMap := make(map[string]byte, 0)
	visitedMap := make(map[string]byte, 0)
	for _, deadNum := range deadends {
		deadMap[deadNum] = '0'
	}
	if _, ok := deadMap["0000"]; ok {
		return -1
	}

	visitSlice := make([]string, 0)
	visitSlice = append(visitSlice, "0000")
	visitedMap["0000"] = '0'
	step := 0

	for len(visitSlice) != 0 {
		size := len(visitSlice)
		for i := 0; i < size; i++ {
			cur := visitSlice[0]
			visitSlice = visitSlice[1:]
			if cur == target {
				return step
			}

			for j := 0; j < 4; j++ {
				for k := -1; k <= 1; k += 2 {
					tmp := cur[0:j] + string((cur[j]-'0'+uint8(10+k))%10+'0') + cur[j+1:]
					if _, ok := deadMap[tmp]; !ok {
						if _, ok := visitedMap[tmp]; ok {
							visitedMap[tmp] = '0'
							visitSlice = append(visitSlice, tmp)
						}
					}
				}
			}
		}
		step++
	}

	return -1
}
