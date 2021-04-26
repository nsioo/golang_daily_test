package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorderTraversal(root *TreeNode) []int {
	res := make([]int, 0)
	stack := make([]*TreeNode, 0)

	for root != nil || len(stack) != 0 {
		// 先走到最左边
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}

		top := stack[len(stack) - 1]
		if len(stack) != 0 {
			stack = stack[:len(stack) - 1]
		}

		res = append(res, top.Val)

		root = top.Right

	}

	return res
}
