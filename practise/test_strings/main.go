package main

import (
	"fmt"
	"strings"
)

func main()  {
	group := "a.group_id"
	ad := "ad_id"

	fmt.Println(group[strings.Index(group, ".") + 1:])
	fmt.Println(ad[:strings.Index(ad, "_")])

}
