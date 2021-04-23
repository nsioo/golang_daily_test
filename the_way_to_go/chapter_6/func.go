package main

import "fmt"

func printOneToTen(n int) {
	if n <= 0 {
		return
	}

	fmt.Println(n)
	printOneToTen(n - 1)
}

func main() {
	printOneToTen(10)
}
