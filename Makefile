

.PHONY: vm-run vm-update vm-delete 

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

