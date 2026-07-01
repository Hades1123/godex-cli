BINARY   := godex
MODULE   := github.com/hades/godex
PREFIX   ?= $(HOME)/.local
BINDIR   := $(PREFIX)/bin
MANDIR   := $(PREFIX)/share/man/man1

GO       := go
BUILDDIR := build

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------

.PHONY: build
build: $(BINARY)

$(BINARY): go.mod go.sum main.go cmd/*.go internal/**/*.go
	$(GO) build -o $@ .

.PHONY: build/windows
build/windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY).exe .

# ---------------------------------------------------------------------------
# Install / Uninstall
# ---------------------------------------------------------------------------

.PHONY: install
install: $(BINARY) install-completions
	@mkdir -p $(BINDIR)
	install -m 0755 $(BINARY) $(BINDIR)/$(BINARY)
	@echo "✅ Installed to $(BINDIR)/$(BINARY)"

.PHONY: uninstall
uninstall: uninstall-completions
	@rm -f $(BINDIR)/$(BINARY)
	@echo "🗑️  Removed $(BINDIR)/$(BINARY)"

# ---------------------------------------------------------------------------
# Quality
# ---------------------------------------------------------------------------

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test:
	$(GO) test -v -count=1 ./...

.PHONY: lint
lint:
	@command -v staticcheck >/dev/null 2>&1 || \
		{ echo "Install staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest"; exit 1; }
	staticcheck ./...

# ---------------------------------------------------------------------------
# Completions (bash, fish)
# ---------------------------------------------------------------------------

COMPS_DIR := completions

.PHONY: completions
completions: $(BINARY)
	@mkdir -p $(COMPS_DIR)
	./$(BINARY) completion bash > $(COMPS_DIR)/godex.bash
	./$(BINARY) completion fish > $(COMPS_DIR)/godex.fish
	@echo "✅ Completions generated (bash, fish)"

# Fish loads completions from ~/.config/fish/completions/<name>.fish automatically.
# Bash needs bash-completion; ~/.local/share/bash-completion/completions/ is on the
# default BASH_COMPLETION_USER_DIR when bash-completion is installed.
.PHONY: install-completions
install-completions: completions
	@mkdir -p $(HOME)/.config/fish/completions
	install -m 0644 $(COMPS_DIR)/godex.fish $(HOME)/.config/fish/completions/godex.fish
	@mkdir -p $(HOME)/.local/share/bash-completion/completions
	install -m 0644 $(COMPS_DIR)/godex.bash $(HOME)/.local/share/bash-completion/completions/godex
	@echo "✅ Completions installed (fish: ~/.config/fish/completions/, bash: ~/.local/share/bash-completion/completions/)"

.PHONY: uninstall-completions
uninstall-completions:
	@rm -f $(HOME)/.config/fish/completions/godex.fish
	@rm -f $(HOME)/.local/share/bash-completion/completions/godex
	@echo "🗑️  Completions removed"

# ---------------------------------------------------------------------------
# Clean
# ---------------------------------------------------------------------------

.PHONY: clean
clean:
	rm -rf $(BINARY) $(BINARY).exe $(COMPS_DIR)
	@echo "✅ Cleaned build artifacts"

# ---------------------------------------------------------------------------
# Help
# ---------------------------------------------------------------------------

.PHONY: help
help:
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Build:'
	@echo '  build            Build for current OS (default)'
	@echo '  build/windows    Cross-compile for Windows'
	@echo ''
	@echo 'Install:'
	@echo '  install          Build + install binary and fish/bash completions'
	@echo '                   (to ~/.local by default)'
	@echo '  uninstall        Remove everything installed by "make install"'
	@echo '  install-completions  Install fish + bash completions only'
	@echo ''
	@echo 'Quality:'
	@echo '  test             Run tests'
	@echo '  vet              go vet'
	@echo '  lint             staticcheck (install separately)'
	@echo ''
	@echo 'Other:'
	@echo '  completions      Generate completion scripts (bash, fish)'
	@echo '  clean            Remove build artifacts'
	@echo ''
	@echo 'Variables:'
	@echo '  PREFIX=~/.local  Installation prefix (default)'
