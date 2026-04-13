

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


release:
	@CURRENT=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	echo "📌 Versión actual: $$CURRENT"; \
	printf "🔖 Nueva versión (ej. v1.2.3): "; \
	read NEW_VERSION; \
	if [ -z "$$NEW_VERSION" ]; then echo "❌ La versión no puede estar vacía"; exit 1; fi; \
	echo "🏷️  Creando tag $$NEW_VERSION..."; \
	git tag -a "$$NEW_VERSION" -m "Release $$NEW_VERSION"; \
	echo "🚀 Compilando release $$NEW_VERSION para: $(PLATFORMS)"; \
	mkdir -p dist; \
	for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; GOARCH=$${platform#*/}; \
		out="dist/$(BINARY_NAME)-$${GOOS}-$${GOARCH}"; \
		echo "  → $$out"; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags "-s -w -X autohost-cli/cmd/autohost-cli.Version=$$NEW_VERSION" -o "$$out" main.go; \
	done; \
	echo "🔐 Generating checksums..."; \
	cd dist && sha256sum $(BINARY_NAME)-* > checksums_$${NEW_VERSION}.txt; \
	echo "✅ Release artifacts in dist/"; \
	ls -lh dist/; \
	echo ""; \
	echo "💡 Para publicar el tag ejecuta: git push origin $$NEW_VERSION"



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
