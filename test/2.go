package main

import (
	"fmt"
)

func main() {
	name := "samrat"
	fmt.Printf("My name is %s\n", name)
	var sum int = 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	if sum < 8 {
		sum /= 2
	} else {
		sum %= 2
	}
	var c string = 'a'
}
