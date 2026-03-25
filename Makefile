

.PHONY: build clean release dev-up vm-run vm-update vm-delete incus-run incus-update incus-delete incus-start incus-stop incus-up

BINARY_NAME = autohost
REPO        = mazapanuwu13/autohost-cli
PLATFORMS   = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
INCUS_INSTANCE = autohost-test

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


dev-up: build
	@echo "🧪 Testing autohost up against local dev (http://localhost:3000)..."
	./$(BINARY_NAME) up --cloud http://localhost:3000

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
	@echo "🚀 Creating Incus instance ($(INCUS_INSTANCE))..."
	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh run

incus-update:
	@echo "🔄 Updating autohost binary in Incus instance ($(INCUS_INSTANCE))..."
	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh update

incus-delete:
	@echo "🧹 Deleting Incus instance ($(INCUS_INSTANCE))..."
	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh delete

incus-start:
	@echo "▶️ Starting Incus instance ($(INCUS_INSTANCE))..."
	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh start

incus-stop:
	@echo "⏹ Stopping Incus instance ($(INCUS_INSTANCE))..."
	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh stop

incus-up: incus-update
	@echo "🧪 Testing autohost up inside Incus instance ($(INCUS_INSTANCE))..."
	@GATEWAY=$$(incus exec $(INCUS_INSTANCE) -- sh -c "ip route show default" | awk '/default/{print $$3}' | head -1); \
	 echo "   Container gateway (host): $$GATEWAY"; \
	 echo "   Cloud URL (browser): http://localhost:3000"; \
	 echo "   API URL (container -> host): http://$$GATEWAY:8080"; \
	 incus exec $(INCUS_INSTANCE) -- env AUTOHOST_CLOUD_URL="http://localhost:3000" AUTOHOST_API_URL="http://$$GATEWAY:8080" /usr/local/bin/autohost up
