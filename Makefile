# This Makefile is meant to be used by people that do not usually work with Go source code.
# If you know what GOPATH is then you probably don't need to bother with make.

MAGE_PATH=${GOPATH}/bin/mage
MAGE=go run ./ci/mage.go -d ./ci

default: help

mage: 
ifeq ("$(wildcard $(MAGE_PATH))","")
	go get -u -d github.com/magefile/mage
endif

help: mage
	${MAGE}

deps: mage
	${MAGE} deps

test: mage deps
	${MAGE} test
