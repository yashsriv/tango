# tango
The Abominably Naive Go Compiler

## Features

Apart from the usual golang imperative features
(except goroutines and chans), we support the following
extra features:

### Currying

```
func myCurry(a int, d int)(b int)(c int) int {
     if (d % 2 == 0) {
         return a + b + c;
     } else {
         return a;
     }
}

x1 = myCurry(4, 2)
x2 = myCurry(4, 2)(5)
x3 = myCurry(4, 1)(5)(6)

// Type of x1:
func(int)(int) int

// Type of x2:
func(int) int

// Type of x3:
int

// Type of myCurry
func(int)(int)(int) int
```

### For Comprehensions

```
xs := []int{1, 2, 3}
doublexs := [x * 2 | _, x := range xs]
```
