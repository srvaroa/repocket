OUTPUT_DIR ?= ./repocket
CONSUMER_KEY ?= 85480-9793dd8ed508561cb941d987

GO = env GO111MODULE=on go

run:
	REPOCKET_OUTPUT_DIR=${OUTPUT_DIR} \
		REPOCKET_CONSUMER_KEY=${CONSUMER_KEY} \
		${GO} run ./cmd/repocket.go dump

list:
	REPOCKET_CONSUMER_KEY=${CONSUMER_KEY} \
		${GO} run ./cmd/repocket.go list

next:
	REPOCKET_CONSUMER_KEY=${CONSUMER_KEY} \
		${GO} run ./cmd/repocket.go next

.PHONY: build
build:
	${GO} build ./...