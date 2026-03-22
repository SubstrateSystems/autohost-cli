

.PHONY: build clean release vm-run vm-update vm-delete incus-run incus-update incus-delete incus-start incus-stop

BINARY_NAME = autohost
REPO        = mazapanuwu13/autohost-cli
PLATFORMS   = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  = -s -w -X autohost-cli/cmd/autohost-cli.Version=$(VERSION)

build:
	@echo "🔨 Building $(BINARY_NAME) $(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) main.go
	@echo "✅ Build complete: ./$(BINARY_NAME)"

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

release:
	@echo "🚀 Building release $(VERSION) for: $(PLATFORMS)"
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/}; \
		out="dist/$(BINARY_NAME)-$${GOOS}-$${GOARCH}"; \
		echo "  → $${out}"; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags "$(LDFLAGS)" -o "$$out" main.go; \
	done
	@echo "🔐 Generating checksums..."
	@cd dist && sha256sum $(BINARY_NAME)-* > checksums_$(VERSION).txt
	@echo "✅ Release artifacts in dist/"
	@ls -lh dist/


vm-run:
	@echo "🚀 Creating Multipass VM ($(VM_NAME))..."
	@bash scripts/autohost-multipass.sh run

vm-update:
	@echo "🔄 Updating autohost binary in VM ($(VM_NAME))..."
	@bash scripts/autohost-multipass.sh update

vm-delete:
	@echo "🧹 Deleting Multipass VM ($(VM_NAME))..."
	@bash scripts/autohost-multipass.sh delete

# ===== Incus ====== #

incus-run:
	@echo "🚀 Creating Incus instance ($(VM_NAME))..."
	@bash scripts/autohost-incus.sh run

incus-update:
	@echo "🔄 Updating autohost binary in Incus instance ($(VM_NAME))..."
	@bash scripts/autohost-incus.sh update

incus-delete:
	@echo "🧹 Deleting Incus instance ($(VM_NAME))..."
	@bash scripts/autohost-incus.sh delete

incus-start:
	@echo "▶️ Starting Incus instance ($(VM_NAME))..."
	@bash scripts/autohost-incus.sh start

incus-stop:
	@echo "⏹ Stopping Incus instance ($(VM_NAME))..."
	@bash scripts/autohost-incus.sh stop
