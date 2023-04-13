package main

import (
	"fmt"
)

func main() {
	for _, col := range collisions {
		fmt.Println(col[0].String(), col[0].Hash())
		fmt.Println(col[1].String(), col[1].Hash())

		if col[0].Hash() != col[1].Hash() || col[0].String() == col[1].String() {
			fmt.Println("not a collision!!!")
		}

		fmt.Println()
	}
}
