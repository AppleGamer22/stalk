SHELL:=/bin/bash
VERSION:=$(shell git describe --tags --abbrev=0 || echo 'main.Version')
HASH:=$(shell git rev-list -1 HEAD)
PACKAGE:=github.com/AppleGamer22/stalk
LDFLAGS:=-ldflags="-X 'main.Version=$(subst v,,$(VERSION))' -X 'main.Hash=$(HASH)'"

test:
	go clean -testcache
	go test -v -race -cover .

debug:
	go build -race $(LDFLAGS) .

completion:
	go run . completion bash > stalk.bash
	go run . completion fish > stalk.fish
	go run . completion zsh > stalk.zsh
	go run . completion powershell > stalk.ps1

manual:
	go run ./utils/replace stalk.1 -b "vVERSION" -a "$(VERSION)"
	go run ./utils/replace stalk.1 -b "DATE" -a "$(shell go run ./utils/date)"

clean:
	rm -rf stalk bin dist stalk.bash stalk.fish stalk.zsh stalk.ps1
	go clean -testcache -cache

.PHONY: debug test clean completion manual