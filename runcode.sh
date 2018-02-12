#!/bin/bash

make
bin/codegen "test/test$1.ir" > "$1.S"
clang -m32 "$1.S"
./a.out
