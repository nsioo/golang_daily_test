package leet_code

//给你两个 非空 的链表，表示两个非负的整数。它们每位数字都是按照 逆序 的方式存储的，并且每个节点只能存储 一位 数字。
//
// 请你将两个数相加，并以相同形式返回一个表示和的链表。
//
// 你可以假设除了数字 0 之外，这两个数都不会以 0 开头。
//
//
//
// 示例 1：
//
//
//输入：l1 = [2,4,3], l2 = [5,6,4]
//输出：[7,0,8]
//解释：342 + 465 = 807.
//
//
// 示例 2：
//
//
//输入：l1 = [0], l2 = [0]
//输出：[0]
//
//
// 示例 3：
//
//
//输入：l1 = [9,9,9,9,9,9,9], l2 = [9,9,9,9]
//输出：[8,9,9,9,0,0,0,1]
//
//
//
//
// 提示：
//
//
// 每个链表中的节点数在范围 [1, 100] 内
// 0 <= Node.val <= 9
// 题目数据保证列表表示的数字不含前导零
//
// Related Topics 递归 链表 数学
// 👍 5942 👎 0

// Definition for singly-linked list.
type ListNode struct {
	Val  int
	Next *ListNode
}

// 新建一个 List 用来存每一位上之和的结果
// 一个数字表示进位
// 遍历长的 List
// 解答失败: 测试用例:[9,9,9,9,9,9,9] [9,9,9,9] 测试结果:[8,9,9,9,0,0,0] 期望结果:[8,9,9,9,0,0,0,1] stdout:
func AddTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
    if l1 == nil {
        return l2
    } else if l2 == nil {
        return l1
    }

    l1Len := getListNodeLen(l1)
    l2Len := getListNodeLen(l2)

    longerList := l1
    shorterList := l2
    if l1Len < l2Len {
        longerList = l2
        shorterList = l1
    }

    res := ListNode{
        -1, nil,
    }
    list := &res
    carry := 0
    for ; shorterList != nil && longerList != nil;  {
        tmpSum := longerList.Val + shorterList.Val + carry
        carry = 0
        if tmpSum > 9 {
            carry = tmpSum / 10
            tmpSum %= 10
        }

        if list.Val == -1 {
            list.Val = tmpSum
        } else {
            list.Next = new(ListNode)
            list = list.Next
            list.Val = tmpSum
        }

        longerList = (*longerList).Next
        shorterList = (*shorterList).Next
    }

    if longerList != nil && carry != 0 {
        for ; longerList != nil; longerList = (*longerList).Next {
            if carry == 0 {
                list.Next = longerList
                break
            }

            tmpSum := longerList.Val + carry
            carry = 0
            if tmpSum > 9 {
                carry = tmpSum / 10
                tmpSum %= 10
            }

            list.Next = new(ListNode)
            list = list.Next
            list.Val = tmpSum
        }
    } else if longerList != nil && carry == 0 {
        list.Next = longerList
    }
    if carry != 0 {
        list.Next = new(ListNode)
        list = list.Next
        list.Val = carry
    }

    return &res
}

func getListNodeLen(list *ListNode) int {
    count := 0
    for node := list; node != nil; node = (*node).Next {
        count++
    }

    return count
}