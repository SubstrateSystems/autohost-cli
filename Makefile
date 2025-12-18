

.PHONY: vm-run vm-update vm-delete incus-run incus-update incus-delete

# ===== MultiPass ====== #

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

