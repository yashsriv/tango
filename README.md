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

## Intermediate Code Representation (IR)

Our IR code is the list of instructions where each instruction is represented as structure 

```
struct {
   type IRType; // Enum values for Type are: {BOP, UOP, CBR, JMP, ASN, KEY} 
   op IROp;
   arg1 IRArg;
   arg2 IRArg;
   dst IRArg;
   target IRTarget;
}
```

### Binary Operations (BOP)

For binary operations (type = BOP), `arg1` and `arg2` are arguments for the operation in the same order and `dst` is the target variable where the result is stored after applying `op`.  

```
dst =  arg1 op arg2
```

`op` can have following values with obvious meanings

```
+, -, *, /, %, <<, >>, &&, ||, &, |, <, <=, >, >=, ==, !=, ^, take (dst = arg1[arg2]), put (arg1[arg2] = dst)
```

`arg1` and `arg2` can be either a variable or a literal. Literals should start with `$` (eg. `$2` means literal 2). Hexadecimals starts with `0x`.  

*Convention*: Temporary variables should start with `r` (r1, r2, etc)

### Unary Operations (UOP)

Unary operations are applied as

```
dst = op arg1
```

Here `op` can have following values

```
neg, !, inc, dec, not, val (Value at address), addr (address of) 
```

### Assignment Operation (ASN)

Assignment is fairly simple

```
dst = arg1
```

### Labels 
 
Instructions can be labelled (with labels being strings) to be referred as target in branch instructions as follows

```
label: dst = arg1 op arg2
```

### Branch Operations (JMP and CBR)

Branch operations can be conditional or unconditional.

Unconditional branches are just a simple call to JMP

```
JMP target
```

For conditional branches, following instructions are provided

```
BREQ arg1 arg2 target  // Branch to target if arg1 == arg2
BRNEQ arg1 arg2 target
BRLT arg1 arg2 target
BRLTE arg1 arg2 target
BRGT arg1 arg2 target
BRGTE arg1 arg2 target
```

### Procedure Call (KEY)

```
PARAM arg1   // pass arg1 as function argument 
CALL target  // calls the function at target
RET          // return to the return address
HALT         // halts the program
```
