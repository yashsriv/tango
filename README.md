# tango
The Abominably Naive Go Compiler

## Building and Running

In order to build this code, you need golang to be installed
on your system.

To install go, see: https://golang.org/doc/install#install

After that, we need your `$GOPATH` to be set appropriately
and `$GOPATH/bin` should be in your `$PATH`:

```sh
export GOPATH=/tmp
export PATH=$PATH:$GOPATH/bin
```

Now install `dep` and `gorunpkg`:
```sh
go get github.com/golang/dep/cmd/dep
go get github.com/Vektah/gorunpkg
```

After that ensure that the files are placed in `$GOPATH/src/tango`.

For example, if we have submitted `final.zip`, unzip it and then:
```
unzip final.zip
mv final $GOPATH/src/
cd $GOPATH/src
mv final tango
```

After that, just `cd` into the folder and run:

```
make
```

to fetch the necessary libs and tools and to generate the lexer and codegen.

To run, follow:

```
Usage: bin/tango [<options>] <file name>
Usage of bin/tango:
  -o string
        Filename of output executable. Only valid with language None (default "./a.out")
  -x string
        Language to output. Possible values are IR, ASM and None (default "None")
```

## Changes from Asgn1 in Asgn3

### Lexer Changes

Some lexical tokens were made independent if they were used in more than one place:
* `*` - Unary, Binary Operator as well as used in types
* `+` - Both  unary and binary operator
* `-` - Both unary and binary operator
* `^` - Both unary and binary operator
* `|` - Last remaining addition-type Binary Operator 
* `mul_op` - Multiply-type binary operators
* `rel_op` - Relational binary operators
* `unary_op` - Remaining unary operators

Different types of assignment statements were made independent lexical tokens

Few new lexical tokens were added:
* `<<<` and `>>>`
* `(|` and `|)`
* `[(` and `)]`
* `|||`

### Syntactical Changes

#### Type Casting

Earlier syntax: `int64(0)`
New syntax: `int64<<<0>>>`

#### Initializing Complex Types

Earlier Syntax: `rect{width: 10, height: 5}`
New Syntax: `rect[(width: 10, height: 5)]`

#### Function Multiple Return Syntax

Earlier Syntax: `func test() (int, int) { //blah }`
New Syntax: `func test() (|int, int|) { //blah }`

#### Method Declaration Syntax

Earlier Syntax: `func (g geometry) test() (int, int) { //blah }`
New Syntax: `func (|g geometry|) test() (int, int) { //blah }`

#### For Comprehension Syntax

Earlier Syntax: `[i * i | i, _ := range arr]`
New Syntax: `[i * i ||| i, _ := range arr]`

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
   type IRType;            // Enum values for Type are: {LBL, BOP, LOP, SOP, DOP, UOP, CBR, JMP, ASN, KEY} 
   op IROp;                // Specific Values of op
   arg1 SymbolTableEntry;  // 
   arg2 SymbolTableEntry;  //
   dst  SymbolTableEntry;  //
}
```

### SymbolTableEntry

A SymbolTableEntry can be of 3 types:
* SymbolTableLiteralEntry: A literal. Just contains value of the literal
* SymbolTableVariableEntry: A variable (virtual register). Contains the memory location of variable in data segment.
* SymbolTableTargetEntry: Contains the actual target for a jmp/branch instruction

The symbol table entries have to be encoded in a specific format in the IR:
* Literals: Must begin with a $
* Variable: Must begin with a r
* Target: Must begin with a #

### Conventions

* All programs have a `_func_main` label.
* All functions have a label beginning with `_func`.

### Operation Types

#### Label Operations (LBL)
A beginning of a Label. Contains only a dst field.

Instructions can be labelled (with labels being strings) to be referred as target in branch instructions as follows

```
label:
```

#### Binary Operations (BOP)

For binary operations (type = BOP), `arg1` and `arg2` are arguments for the operation in the same order and `dst` is the target variable where the result is stored after applying `op`.

```
dst =  arg1 op arg2
```

`op` can have following values:
* `+`: Add
* `-`: Subtract
* `*`: Multiply
* `&`: Bitwise AND
* `|`: Bitwise OR
* `^`: Bitwise XOR
* `&&`: Logical AND
* `||`: Logical OR


`arg1` and `arg2` can be either a variable or a literal.

#### Logical Operations (LOP)

For logical operations, `arg1` and `arg2` are arguments for the operation in the same order and `dst` is the target variable where the result is stored after applying `op`.

```
dst =  arg1 op arg2
```

`op` can have following values:
* `<`: Less Than
* `>`: Greater Than
* `<=`: Less Than Equal
* `>=`: Greater Than Equal
* `==`: Equals
* `!=`: Not Equals

`arg1` and `arg2` can be either a variable or a literal (having value either 0 or 1).

#### Shift Operations (SOP)
For shift operations, `arg1` and `arg2` are arguments for the operation in the same order and `dst` is the target variable where the result is stored after applying `op`.

```
dst =  arg1 op arg2
```

`op` can have following values:
* `<<`: Bitwise Shift Left
* `>>`: Bitwise Shift Right

`arg1` and `arg2` can be either a variable or a literal (having value either 0 or 1).

#### Division Operations (DOP)
For division operations, `arg1` and `arg2` are arguments for the operation in the same order and `dst` is the target variable where the result is stored after applying `op`.

```
dst =  arg1 op arg2
```

`op` can have following values:
* `/`: Divide
* `%`: Remainder

`arg1` and `arg2` can be either a variable or a literal (having value either 0 or 1).

#### Unary Operations (UOP)

Unary operations are applied as

```
dst = op arg1
```

Here `op` can have following values:
* `neg`: Negate a value
* `not`: Bitwise Not
* `!`: Logical Not

#### Assignment Operation (ASN)

Assignment is fairly simple

```
dst = arg1
```

#### Branch Operations (JMP and CBR)

Branch operations can be conditional or unconditional.

Unconditional branches are just a simple call to JMP

```
jmp target
```

For conditional branches, following instructions are provided

```
breq target arg1 arg2  // Branch to target if arg1 == arg2
brlt target arg1 arg2  // Branch to target if arg2 < arg1 
brgt target arg1 arg2  // Branch to targtet if arg2 > arg1
brlte arg1 arg2 target // Branch to targtet if arg2 <= arg1
brgte arg1 arg2 target // Branch to targtet if arg2 >= arg1
brneq target arg1 arg2 // Branch to targtet if arg2 != arg1
```

#### Special Instructions (KEY)

```
inc arg1     // Increments the virtual register
dec arg1     // Decrements the virtual register
call dst     // calls the function at dst target
param arg1   // Push to stack
ret          // return from a function
reti         // return a value from a function
setret arg1  // set arg1 to return value of a function 
halt         // halts the program. We can pass a status code by using param
printi arg1  // Prints a
printc arg1  // Prints a
prints arg1  // Prints a
scani  arg1  // Scan into a
scanc  arg1  // Scan into a
scans  arg1  // Scan into a
```
