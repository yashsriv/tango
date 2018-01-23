mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

GOCC := ./bin/gocc

.Phony: all clean

all: export GOPATH=$(current_dir)
all: export GOBIN=$(current_dir)/bin
all: bin/gocc bin/lexer

bin/gocc:
	@echo -e "\e[1;34mFetching gocc \e[0m"
	go get -v github.com/goccmack/gocc

bin/lexer:
	@echo -e "\e[1;34mCompiling Lexer \e[0m"
	go install src/tango/lexer/lexer.go

clean:
	@echo -e "\e[1;34mCleaning Files \e[0m"
	@rm -rf pkg bin/**
