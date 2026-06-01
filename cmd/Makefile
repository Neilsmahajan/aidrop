# aidrop Makefile
# Author: Neil Mahajan <neilsmahajan@gmail.com>

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=aidrop
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_MACOS=$(BINARY_NAME)_macos

# Build the project
.PHONY: all
all: test build

.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_MACOS)

.PHONY: run
run: build
	./$(BINARY_NAME)

.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) verify

.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Cross compilation
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v

.PHONY: build-macos
build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MACOS) -v

.PHONY: build-all
build-all: build-linux build-windows build-macos

# Development helpers
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

.PHONY: vet
vet:
	$(GOCMD) vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: coverage
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

.PHONY: benchmark
benchmark:
	$(GOTEST) -bench=. -benchmem ./...

# Install/Uninstall
.PHONY: install
install:
	$(GOCMD) install

.PHONY: uninstall
uninstall:
	rm -f $(shell $(GOCMD) env GOPATH)/bin/$(BINARY_NAME)

# Shell completions
# Generates completion scripts and installs them to the standard system locations.
# After running, reload your shell or source the installed file.
.PHONY: completion-bash
completion-bash: build
	@mkdir -p /etc/bash_completion.d 2>/dev/null || mkdir -p $(HOME)/.local/share/bash-completion/completions
	./$(BINARY_NAME) completion bash > /etc/bash_completion.d/$(BINARY_NAME) 2>/dev/null || \
		./$(BINARY_NAME) completion bash > $(HOME)/.local/share/bash-completion/completions/$(BINARY_NAME)
	@echo "Bash completion installed. Restart your shell or run: source ~/.bashrc"

.PHONY: completion-zsh
completion-zsh: build
	@mkdir -p $(HOME)/.zsh/completions
	./$(BINARY_NAME) completion zsh > $(HOME)/.zsh/completions/_$(BINARY_NAME)
	@echo "Zsh completion installed to ~/.zsh/completions/_$(BINARY_NAME)"
	@echo "Make sure your .zshrc contains: fpath=(~/.zsh/completions \$$fpath) && autoload -U compinit && compinit"

.PHONY: completion-fish
completion-fish: build
	@mkdir -p $(HOME)/.config/fish/completions
	./$(BINARY_NAME) completion fish > $(HOME)/.config/fish/completions/$(BINARY_NAME).fish
	@echo "Fish completion installed to ~/.config/fish/completions/$(BINARY_NAME).fish"

.PHONY: completion
completion: completion-bash completion-zsh completion-fish

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              - Run tests and build"
	@echo "  build            - Build the binary"
	@echo "  test             - Run tests"
	@echo "  clean            - Clean build artifacts"
	@echo "  run              - Build and run the binary"
	@echo "  deps             - Download dependencies"
	@echo "  tidy             - Tidy module dependencies"
	@echo "  build-linux      - Cross-compile for Linux"
	@echo "  build-windows    - Cross-compile for Windows"
	@echo "  build-macos      - Cross-compile for macOS"
	@echo "  build-all        - Cross-compile for all platforms"
	@echo "  fmt              - Format Go code"
	@echo "  vet              - Run go vet"
	@echo "  lint             - Run golangci-lint"
	@echo "  coverage         - Generate test coverage report"
	@echo "  benchmark        - Run benchmarks"
	@echo "  install          - Install binary to GOPATH/bin"
	@echo "  uninstall        - Remove binary from GOPATH/bin"
	@echo "  completion-bash  - Install bash tab completion"
	@echo "  completion-zsh   - Install zsh tab completion"
	@echo "  completion-fish  - Install fish tab completion"
	@echo "  completion       - Install all shell completions"
	@echo "  help             - Show this help message"

