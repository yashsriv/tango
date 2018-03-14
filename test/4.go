package main

import "fmt"

type test struct {
	foo int
	bar string
}

type tester interface {
	Test() bool
}

const five = 5

func main() {
	var a = 'a'
	var b int = 0
	myMap := map[string]int[(
		"a": 5,
		"b": 7,
	)]
	for key, count := range myMap {
		fmt.Println(key, count)
	}
	defer fmt.Println(a)
	for {
		a ^= a
		b++
	}
}
