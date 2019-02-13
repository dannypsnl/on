GO_MOD := GO111MODULE=on

.PHONY: build test install
build:
	@$(GO_MOD) go build
test:
	@$(GO_MOD) go test ./... -v -cover -count 1
install:
	@$(GO_MOD) go install
