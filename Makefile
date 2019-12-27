# Some make settings
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c  
.ONESHELL:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# Repocket specific
GO = env GO111MODULE=on go

dump:
		${GO} run ./cmd/repocket.go dump

list:
		${GO} run ./cmd/repocket.go list

next:
		${GO} run ./cmd/repocket.go next

.PHONY: build
build:
	${GO} build ./...
