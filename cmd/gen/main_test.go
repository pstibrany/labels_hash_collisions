package main

import (
	"fmt"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		s, _, _ := generateRandomString(5, nil, nil)
		fmt.Println(s)
	}
}
