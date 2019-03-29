GO_MOD := GO111MODULE=on

.PHONY: zsh-completion bash-completion build test install
zsh-completion:
	@echo "source $(shell pwd)/on-completion.sh" >> ~/.zshrc
bash-completion:
	@echo "source $(shell pwd)/on-completion.sh" >> ~/.bashrc
build:
	@$(GO_MOD) go build
test:
	@$(GO_MOD) go test ./... -v -cover -count 1
install:
	@$(GO_MOD) go install
