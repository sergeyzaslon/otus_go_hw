package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const string4reverse = "Hello, OTUS!"

func main() {
	fmt.Println(stringutil.Reverse(string4reverse))
}
