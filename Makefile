mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

GOCC := ./bin/gocc

.PHONY: all clean libs test

all: export GOPATH=$(current_dir)
all: export GOBIN=$(current_dir)/bin
all: libs bin/lexer

libs: bin/gocc src/github.com/ryanuber/columnize

bin/gocc:
	@echo -e "\e[1;34mFetching gocc \e[0m"
	go get -v github.com/goccmack/gocc

src/github.com/ryanuber/columnize:
	@echo -e "\e[1;34mFetching columnize \e[0m"
	go get -v github.com/ryanuber/columnize

bin/lexer: src/tango/main/lexer/lexer.go src/tango/lexer/lexer.go
	@echo -e "\e[1;34mCompiling Lexer \e[0m"
	go install $(current_dir)/src/tango/main/lexer/lexer.go

src/tango/lexer/lexer.go: src/tango/tango.ebnf
	@echo -e "\e[1;34mGenerating Lexer \e[0m"
	cd $(current_dir)/src/tango && $(GOBIN)/gocc tango.ebnf

test: export GOPATH=$(current_dir)
test: export GOBIN=$(current_dir)/bin
test:
	go test tango/lexer

clean:
	@echo -e "\e[1;34mCleaning Files \e[0m"
	@rm -rf $(current_dir)/pkg $(current_dir)/bin/**
	@rm -rf $(current_dir)/src/tango/util
	@rm -rf $(current_dir)/src/tango/token
	@rm $(current_dir)/src/tango/lexer/lexer.go
	@rm $(current_dir)/src/tango/lexer/acttab.go
	@rm $(current_dir)/src/tango/lexer/transitiontable.go
