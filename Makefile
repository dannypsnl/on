GO_MOD := GO111MODULE=on

.PHONY: build
build:
	@$(GO_MOD) go build
test:
	@$(GO_MOD) go test ./... -v -cover -count 1
