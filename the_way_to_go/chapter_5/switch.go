package main

import (
	"fmt"
)

func getSeason(month int) {
	switch month {
	case 3, 4, 5:
		fmt.Println("It's Spring!")
	case 6, 7, 8:
		fmt.Println("It's Summer!")
	case 9, 10, 11:
		fmt.Println("It's Autumn!")
	case 12, 1, 2:
		fmt.Println("It's Winter！")
	}
}

func main() {
	month := 0
	_, _ = fmt.Scan(&month)

	getSeason(month)

}