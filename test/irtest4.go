package main

import "fmt"

func fact(i int) int {
	if i == 0 {
		return 1
	}
	return i * fact(i-1)
}

func main() {
	var arr [3]int
	arr[1] = 5
	var j int = arr[1] + 2
	fact(j)
}
