# ----------------------------
# AutoHost CLI - DevOps tasks
# ----------------------------

.PHONY: vm-run vm-update vm-delete

vm-run:
	@echo "🚀 Creating Multipass VM (autohost-test)..."
	@bash scripts/autohost-multipass.sh run

vm-update:
	@echo "🔄 Updating autohost binary in VM..."
	@bash scripts/autohost-multipass.sh update

vm-delete:
	@echo "🧹 Deleting Multipass VM (autohost-test)..."
	@bash scripts/autohost-multipass.sh delete
