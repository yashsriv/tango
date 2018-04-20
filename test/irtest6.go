package main

import "fmt"

type mystruct struct {
	val int
}

func (| s mystruct |) Val() int {
	return s.val
}

func fact(i int) int {
	if i == 0 {
		return 1
	}
	return i * fact(i-1)
}

func main() {
	var s mystruct
	var i int
	scanf("%d", &i)
	s.val = i
	printf("Hello %d\n", fact(s.Val()))
}
