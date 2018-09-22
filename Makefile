# This Makefile is meant to be used by people that do not usually work with Go source code.
# If you know what GOPATH is then you probably don't need to bother with make.

default: help

help:
	go run ci/main.go help

deps:
	go run ci/main.go deps

test: deps
	go run ci/main.go test