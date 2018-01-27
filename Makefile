mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

export GOPATH=$(current_dir)
export GOBIN=$(current_dir)/bin

GOCC := $(current_dir)/bin/gocc

.PHONY: all clean libs test

all: libs bin/lexer

libs: bin/gocc src/github.com/olekukonko/tablewriter

bin/gocc:
	@echo -e "\e[1;34mFetching gocc \e[0m"
	go get -v github.com/goccmack/gocc

src/github.com/olekukonko/tablewriter:
	@echo -e "\e[1;34mFetching tablewriter \e[0m"
	go get -v github.com/olekukonko/tablewriter

bin/lexer: src/tango/main/lexer/lexer.go src/tango/lexer/lexer.go
	@echo -e "\e[1;32mCompiling Lexer \e[0m"
	go install $(current_dir)/src/tango/main/lexer/lexer.go

src/tango/lexer/lexer.go: src/tango/tango.ebnf
	@echo -e "\e[1;33mGenerating Lexer \e[0m"
	cd $(current_dir)/src/tango && $(GOCC) tango.ebnf

test:
	go test tango/lexer

clean:
	@echo -e "\e[1;31mCleaning Files \e[0m"
	@echo -e "\e[1;31m  Clearing pkg and bin \e[0m"
	@rm -rf $(current_dir)/pkg $(current_dir)/bin/**
	@echo -e "\e[1;31m  Clearing generated files \e[0m"
	@rm -rf $(current_dir)/src/tango/util
	@rm -rf $(current_dir)/src/tango/token
	@rm -rf $(current_dir)/src/tango/lexer/lexer.go
	@rm -rf $(current_dir)/src/tango/lexer/acttab.go
	@rm -rf $(current_dir)/src/tango/lexer/transitiontable.go

nuke: clean
	@echo -e "\e[1;31m  Clearing downloaded libraries \e[0m"
	@rm -rf $(current_dir)/src/github.com
