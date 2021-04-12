package main

import "fmt"

func main()  {
	x := []int{}
	x = append(x, 0)
	x = append(x, 1)
	x = append(x, 2)
	y := append(x, 3)
	z := append(x, 4)
	fmt.Println(x, y, z)
	fmt.Printf("%p\n%p\n%p\n", x, y, z)
}
