package main

import "fmt"

func cutSliceByIndex(str string, index int) []string {
	res := make([]string, 0)

	res = append(res, str[:index])
	res = append(res, str[index + 1: ])
	return res
}

// 7.16 bubble sort
func BubbleSort(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		for j := i; j < len(arr); j++ {
			if arr[i] > arr[j] {
				tmp := arr[i]
				arr[i] = arr[j]
				arr[j] = tmp
			}
		}
	}

	return arr
}

func mulTen(arr []int) []int {
	for index := range arr {
		arr[index] = arr[index] * 10
	}

	return arr
}



func main() {
	//s := []int{1, 2, 3}
	//
	//factor := 19
	//
	//if len(s) + factor > cap(s) {
	//	newSlice := make([]int, len(s) * factor)
	//	copy(newSlice, s)
	//	s = newSlice
	//}
	//
	//fmt.Println(s)

	//tmp := "123456789101111"
	//a := 0
	//s := "123"
	//s1 := s[1:]
	//var tmp2 string = tmp[0:3]
	//tmpByte := []byte(tmp)
	//
	//fmt.Println((*reflect.StringHeader)(unsafe.Pointer(&tmp)).Data)
	//fmt.Println((*reflect.StringHeader)(unsafe.Pointer(&tmp2)).Data)
	//fmt.Println((*reflect.SliceHeader)(unsafe.Pointer(&tmpByte)).Data)
	//
	//fmt.Printf("tmp %p\n", &tmp)
	//fmt.Printf("tmp2 %p\n", &tmp2)
	//fmt.Printf("%p\n", &tmpByte)
	//fmt.Printf("%p\n", &a)
	//
	//tmpByte[0] = '9'
	//
	//fmt.Println(tmp)
	//fmt.Println(tmp2)
	//fmt.Println(s1)
	//
	//fmt.Println(string(tmpByte))

	// 7.16
	unsortArr := []int{4, 2, 7, 4, 3, 2, 1, 6}
	sortArr := BubbleSort(unsortArr)
	fmt.Println(sortArr)

}
