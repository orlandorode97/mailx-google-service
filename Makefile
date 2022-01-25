SHELL=/usr/bin/env bash
PROJECTNAME=$(shell basename "$(PWD)")
LDFLAGS="-X 'main.buildTime=$(shell date)' -X 'main.lastCommit=$(shell git rev-parse HEAD)' -X 'main.semanticVersion=$(shell git describe --tags --dirty=-dev)'"
GO=go

## help: Get more info on make commands.
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help

## build: Build mailx-google-service binary.
build:
	@echo "--> Building mailx-google-service"
	docker-compose build
.PHONY: build

run:
	@echo "--> Running mailx-google-service"
	docker-compose up --remove-orphans
.PHONY: run

build-run: build run

goose-build:
	@echo "Building goose binary --->"
	${GO} build -o . ./cmd/goose
	@echo "Goose binary built"
.PHONY: gooose-build

up:
	@echo "Up migrations"
	./goose up
	@echo "Migrations up successfully"
.PHONY: up

down:
	@echo "Down migrations"
	./goose down
	@echo "Migrations down successfully"
.PHONY: down

status:
	@echo "Goose status"
	./goose status
.PHONY: status

clear:
	@echo "Removing goose binary"
	rm -rf ./goose
	@echo "goose binary removed successfully"
.PHONY: clear

goose-up: goose-build up
goose-down: goose-build down
goose-status: goose-build status
goose-clear: clear

