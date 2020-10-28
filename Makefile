#include .env

PROJECTNAME=$(shell basename "$(PWD)")

GOBASE=$(shell pwd)
# GOBIN=$(GOBASE)/bin
# GOFILES=$(wildcard *.go)

# MAKEFLAGS += --silent

# variables:
# 	@echo "$(PROJECTNAME)"
# 	@echo "$(GOBASE)"
# 	@echo "$(GOBIN)"
# 	@echo "$(GOFILES)"

start:
	docker-compose start

stop:
	docker-compose stop

restart:
	docker-compose restart

build:
	go build -o bin/$(PROJECTNAME) main.go

run:
	go run $(PROJECTNAME)/bin/fReader

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go