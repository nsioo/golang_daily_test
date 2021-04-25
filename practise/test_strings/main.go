package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Rope string

func testHasPrefix(s, prefix string) {
	fmt.Printf("Does '%s' has %s prefix: %v\n", s, prefix, strings.HasPrefix(s, prefix))
}

func testLastIndex(s, substr string) {
	fmt.Println(strings.LastIndex(s, substr))
}

func main()  {
	group := "a.group_id"
	ad := "ad_id"

	fmt.Println(group[strings.Index(group, ".") + 1:])
	fmt.Println(ad[:strings.Index(ad, "_")])

	var align Rope = "test"
	fmt.Println(align)

	testHasPrefix("This is QQMAN", "h")

	testLastIndex("app_id:uid", ":")

	println(strconv.IntSize)
}
