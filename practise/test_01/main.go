package main

import "fmt"

func main() {
	//body, _ := uuid.GenerateUUID()
	//fmt.Println(strings.ToUpper(body))
	//fmt.Println(body)
	//
	//test := map[string]int {
	//	"1": 1,
	//	"2": 2,
	//}
	//
	//k, ok := test["1"]

	tmp := "123456"
	tmp1 := tmp[1:]
	fmt.Printf("%p\n", &tmp)
	fmt.Printf("%p", &tmp1)


}
