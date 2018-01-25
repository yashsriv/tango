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

bin/lexer: src/tango/lexer/main.go src/tango/lexer/lexer.go
	@echo -e "\e[1;34mCompiling Lexer \e[0m"
	go install $(current_dir)/src/tango/lexer/main.go

src/tango/lexer/lexer.go: src/tango/tango.ebnf
	@echo -e "\e[1;34mGenerating Lexer \e[0m"
	cd $(current_dir)/src/tango && $(GOBIN)/gocc tango.ebnf

clean:
	@echo -e "\e[1;34mCleaning Files \e[0m"
	@rm -rf $(current_dir)/pkg $(current_dir)/bin/**
	@rm -rf $(current_dir)/src/tango/util
	@rm -rf $(current_dir)/src/tango/token
	@rm $(current_dir)/src/tango/lexer/lexer.go
	@rm $(current_dir)/src/tango/lexer/acttab.go
	@rm $(current_dir)/src/tango/lexer/transitiontable.go
