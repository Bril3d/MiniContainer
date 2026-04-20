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

## 🏗️ Architecture

MiniContainer consists of two primary components joined seamlessly:

### 1. Go Backend & CLI
Built using Go 1.25.0 and the Cobra CLI framework.
- **CLI (`cmd/`)**: Docker-like commands wrapped around Podman for easier local development. Manages lifecycle (`up`, `start`, `stop`, `run`), cleanup (`rm`, `down`), debugging (`exec`, `logs`, `stats`), and presets.
- **API Server**: Starts a Gin-based RESTful API on `localhost:8080`, exposing Podman operations natively to the GUI. Seamlessly handles automatic launching of the Tauri GUI dev server.
- **Interoperability**: Contains built-in volume-mounting logic for Windows to translate paths automatically for WSL Podman interop.

### 2. Frontend GUI (`gui/`)
A fast desktop application built on top of the Tauri v2 desktop runtime.
- **Tech Stack**: React 19 + Vite + TailwindCSS 4 + Framer Motion.
- **Operation Mechanics**: Executes cross-origin HTTP requests to the Go API backend to retrieve container stats, list images, and trigger lifecycle actions natively without relying on raw Tauri rust-commands.

## 📋 Requirements

- [Podman](https://podman.io/getting-started/installation) installed
- Windows (with WSL2) or Linux

## 🚧 Status

Under active development. The CLI offers essential Podman orchestration, and the Tauri GUI interacts seamlessly via the built-in Go REST API.

## 📜 License

[MIT](LICENSE)
