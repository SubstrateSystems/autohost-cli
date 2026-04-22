![autohost-cli](./autohost-cli.png)


# 🚀 AutoHost CLI

**Take back control of your services.**  
**AutoHost CLI** is a command-line tool to install, configure, and manage applications and services **on your own server**, without depending on third parties and with a simple and automated workflow.

---

## 🌟 Features

- **One-command installation**: Deploy ready-to-use applications with `app install`.
- **Multi-app support**: Nextcloud, BookStack, Redis, MySQL, PostgreSQL, and more (constantly growing!).
- **Tailscale integration**: Securely connect to your private infrastructure.
- **Docker compatibility**: Isolation and portability for your applications.
- **Privacy and control focus**: Everything runs on **your** infrastructure.

---

## ⚙️ Prerequisites

Before installing, make sure you have:
- A **Linux**-based system (compatible with modern distributions like Ubuntu/Debian).  
- **Docker** installed and running.  
- Administrator permissions (**sudo/root**).  
- Optional: **Tailscale** account if you want to enable secure private access.  

---

## 📦 Installation

Install AutoHost CLI directly from GitHub with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/mazapanuwu13/autohost-cli/main/scripts/install.sh | bash
```

The script auto-detects your OS and architecture, downloads the latest binary from GitHub Releases, verifies the checksum, and installs it.

If the binary lands in `~/.local/bin`, add it to your PATH:

```bash
export PATH="$PATH:$HOME/.local/bin"
source ~/.bashrc
```

### Specific version

```bash
VERSION=v0.4.0 curl -fsSL https://raw.githubusercontent.com/mazapanuwu13/autohost-cli/main/scripts/install.sh | bash
```

### Verify installation

```bash
autohost --version
```

---

## 🛠 Basic Usage

### Example workflow

```bash
# Initialize environment
autohost init

# Initial setup (domain, networks, etc.)
autohost setup

# Install an application (example: Nextcloud)
autohost app install

# List installed applications
autohost app ls

# Start the application
autohost app start nextcloud

# Check app status
autohost app status nextcloud

# Stop the application
autohost app stop nextcloud

# Remove an application
autohost app remove nextcloud
```

---

## 📂 Supported Applications

| App        | Default Port | Status      |
|------------|-------------|-------------|
| Nextcloud  | 8081        | ✅ Stable   |
| BookStack  | 6875        | ✅ Stable   |
| MySQL      | 3306        | ✅ Stable   |
| PostgreSQL | 5432        | ✅ Stable   |
| Redis      | 6379        | ✅ Stable   |

*(The list grows with each version. Your feedback helps prioritize new apps!)*

---

## 🏗 Architecture

The project follows **Clean Architecture** principles with the following structure:

```
autohost-cli/
├── cmd/                    # CLI commands (Cobra)
├── internal/
│   ├── adapters/          # External integrations (Docker, Caddy, Tailscale, etc.)
│   │   └── storage/       # Database repositories (SQLite)
│   ├── app/               # Application services (business logic)
│   ├── domain/            # Domain models and interfaces
│   └── platform/          # Platform utilities (config, DI, filesystem)
├── db/                    # Database migrations and seeds
├── assets/                # Embedded templates (docker-compose files)
├── utils/                 # Utility functions
└── scripts/               # Installation and testing scripts
```

### Key Components

- **Domain Layer**: Core business logic and interfaces
- **Application Layer**: Use cases and service orchestration
- **Adapters Layer**: External integrations (Docker, Caddy, Tailscale, CloudFlare, etc.)
- **Infrastructure**: Database, configuration, and platform-specific code

---

## 🧪 Testing Environment

### Multipass VM for Testing

You need to have **Multipass** installed: https://canonical.com/multipass

| Command                              | Description                                                    |
|--------------------------------------|----------------------------------------------------------------|
| `make vm-run`                        | Creates VM (autohost-test) with autohost binary in bin folder |
| `make vm-update`                     | Updates autohost binary in VM (autohost-test) bin folder      |
| `make vm-delete`                     | Deletes the VM (autohost-test)                                |

---

## 🔒 Philosophy

In a world where most applications are in the cloud, **AutoHost CLI** gives you back the power:  
- You control **your data**.  
- You eliminate dependency on multiple SaaS providers.  
- You build your own scalable and private infrastructure.  

---

## 🤝 Contributing

Want to contribute?  
1. Fork the repository.  
2. Create a branch for your feature/fix.  
3. Submit a Pull Request.  
4. Check issues labeled **good first issue** to get started.

### Development Guidelines

- Follow Go 1.23.0 standards
- Use `gofmt` and `goimports` for formatting
- Run `go vet` before committing
- Keep functions small and focused
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Maintain idempotency in commands
- No secrets in logs or error messages

---

## 📜 License

This project is licensed under the **MIT License**.

---

> 💡 **Tip:** For updates and news, visit [autohst.dev](https://autohst.dev) or follow us on social media.
