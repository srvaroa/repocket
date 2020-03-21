BIN_DIR=./bin
BINARY=repocket
LINUX_ARCHS=386 amd64 arm arm64
DARWIN_ARCHS=386 amd64

VERSION=0.1

build_all:
	mkdir -p ${BIN_DIR}
	$(foreach GOARCH, $(LINUX_ARCHS), \
		$(shell export GOOS=linux; export GOARCH=$(GOARCH); export GO111MODULE=auto; go build -v -o $(BIN_DIR)/$(BINARY)-linux-$(GOARCH)-$(VERSION) cmd/repocket.go) \
	)
	$(foreach GOARCH, $(DARWIN_ARCHS), \
		$(shell export GOOS=darwin; export GOARCH=$(GOARCH); export GO111MODULE=auto; go build -v -o $(BIN_DIR)/$(BINARY)-darwin-$(GOARCH)-$(VERSION) cmd/repocket.go) \
	)

clean:
	rm -r ${BIN_DIR} || true
