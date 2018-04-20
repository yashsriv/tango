package main

import "fmt"

func fact(i int) int {
	if i == 0 {
		return 1
	}
	return i * fact(i-1)
}

func main() {
	var arr []int = []int[(1, 5, 6)]
	var j int = arr[0] + 2
	printf("%d\n", fact(j))
}
