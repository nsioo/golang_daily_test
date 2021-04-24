package leet_code

func numIslands(grid [][]byte) int {
	if grid == nil || len(grid) == 0 {
		return -1
	}

	m := len(grid)
	n := len(grid[0])
	islandCount := 0

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				islandCount++

				dfsIsland(i, j, &grid)
			}
		}
	}

	return islandCount
}


func dfsIsland(i, j int, grid *[][]byte) {
	if i >= len(*grid) || j >= len((*grid)[0]) ||
		i < 0 || j < 0 ||
		(*grid)[i][j] == '0' {
		return
	}

		(*grid)[i][j] = '0'

		dfsIsland(i + 1, j, grid)
		dfsIsland(i - 1, j, grid)
		dfsIsland(i, j + 1, grid)
		dfsIsland(i, j - 1, grid)
}