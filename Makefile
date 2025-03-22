# Go settings
GO := go
BIN_DIR := _bin
CMD_DIR := cmd
GOPATH_BIN := $(shell go env GOPATH)/bin

# Find all commands under cmd/
CMDS := $(shell find $(CMD_DIR) -mindepth 1 -maxdepth 1 -type d -exec basename {} \;)

# Generate binary paths in _bin/
BINARIES := $(CMDS:%=$(BIN_DIR)/%)

# Default target: Build all binaries
all: $(BINARIES)

# Rule to build each command
$(BIN_DIR)/%: $(CMD_DIR)/%/*.go
	mkdir -p $(BIN_DIR)
	$(GO) build -o $@ ./$(CMD_DIR)/$*

# Install commands to GOPATH/bin
install: $(BINARIES)
	@for bin in $(CMDS); do \
		echo "Installing $$bin to $(GOPATH_BIN)..."; \
		cp $(BIN_DIR)/$$bin $(GOPATH_BIN)/$$bin; \
	done

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
