package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())
	for {
		tmp := rand.Intn(2)
		fmt.Println(tmp)
		if tmp == 1{
			break
		}

	}
}
