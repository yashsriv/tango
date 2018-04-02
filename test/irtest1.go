package main

const (
	a int = 5
	b int = 6
	c int = 7
	d int = 8
)

var (
	x, y int = a + c, 6
)

var z int

func main() {
	var g int = 5
	g += 6

	if 1 < 5 {
		g += 5
	} else if 5 < 7 && g == 11 {
		g = 15
	} else if 80 {
		g = 4
	} else {
		g = 7
	}

	var i int
	for i = 0; i < 5; i++ {
		g += 2
	}
}
