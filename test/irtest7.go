package main

import "fmt"

type mystruct struct {
	val int
}

type mystruct1 struct {
	val mystruct
}

func fact(i int) int {
	if i == 0 {
		return 1
	}
	return i * fact(i-1)
}

func main() {
	var s mystruct1
	s.val.val = 5
	fact(s.val.val)
}
