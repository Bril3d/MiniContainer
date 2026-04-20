# 🚀 MiniContainer

> Run dev environments instantly — no Docker complexity, no heavy setup.

MiniContainer is a **lightweight, fast container management tool** built on top of Podman. It provides a simple CLI and an intuitive GUI to run containers and development environments without the resource overhead of Docker Desktop.

## ✨ Features

- 🏃 **One-command environments** — `mini run python`, `mini run postgres`
- ⚡ **Daemonless** — Zero idle CPU/RAM usage. Relies actively on Podman.
- 📦 **Presets** — Pre-configured environments for Python, Node, Postgres, MongoDB, Jupyter
- 🔍 **Smart detection** — Auto-detects project type from your files
- 📊 **Resource monitoring** — Real-time CPU/RAM usage per container
- 📄 **MiniFile Parsing** — A lightweight `docker-compose` alternative for single-file environment orchestration (`mini up`).
- 🖥️ **Desktop GUI** — Fast, modern React & Tauri-based interface for managing containers visually.

## 📋 Requirements

- [Podman](https://podman.io/getting-started/installation) installed
- Windows (with WSL2) or Linux
- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) & [npm](https://www.npmjs.com/) (for GUI development)

## 🚀 Getting Started

### 1. Start the Podman Machine
MiniContainer runs on top of Podman. Ensure your Podman machine is running:
```bash
podman machine start
```

### 2. Run the Full Stack (API + GUI)
To launch the backend API and the desktop interface simultaneously:
```bash
go run main.go serve
```
*The API will start on `http://localhost:8080` and the Tauri GUI will launch automatically.*

### 3. Run API Only
If you prefer using only the CLI or your own frontend:
```bash
go run main.go serve --no-gui
```

### 4. Build and Use the CLI
You can compile MiniContainer into a standalone executable:
```bash
# Build the binary
go build -o mini.exe main.go

# Use it like Docker
.\mini.exe ps
.\mini.exe run alpine echo "Hello from Mini"
```

## 🏗️ Architecture

MiniContainer consists of two primary components joined seamlessly:

### 1. Go Backend & CLI
Built using Go 1.25.0 and the Cobra CLI framework.
- **CLI (`cmd/`)**: Docker-like commands wrapped around Podman for easier local development. Manages lifecycle (`up`, `start`, `stop`, `run`), cleanup (`rm`, `down`), debugging (`exec`, `logs`, `stats`), and presets.
- **API Server**: Starts a Gin-based RESTful API on `localhost:8080`, exposing Podman operations natively to the GUI.
- **Interoperability**: Contains built-in volume-mounting logic for Windows to translate paths automatically for WSL Podman interop.

### 2. Frontend GUI (`gui/`)
A fast desktop application built on top of the Tauri v2 desktop runtime.
- **Tech Stack**: React 19 + Vite + TailwindCSS 4 + Framer Motion.
- **Operation Mechanics**: Executes cross-origin HTTP requests to the Go API backend to retrieve container stats, list images, and trigger lifecycle actions.

## 🚧 Status

Under active development. CLI and GUI MVP are functional.

## 📜 License

[MIT](LICENSE)
