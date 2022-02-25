package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	testString := "Hello, OTUS!"

	fmt.Println(stringutil.Reverse(testString))
}
