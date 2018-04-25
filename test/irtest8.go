package main

import "fmt"

func main() {
	var f func(i int, j int) int = func(i int, j int) int {
		return i + j
	}
	var k int = f(2, 4)
	printf("%d\n", k)
}
