# 🛡️ Go Password Manager

A secure, modular, and production-ready password manager service built with **Go**, using **gRPC** for communication and backed by **PostgreSQL**, **Redis**, and **MinIO** for persistence, caching, and object storage.

---

## 🚀 Features

- ✅ Secure gRPC API with TLS
- 🗄️ PostgreSQL-backed data storage
- ⚡ Redis for caching
- 📦 MinIO (S3-compatible) for file/object storage
- 🧪 Unit testing with coverage
- 📁 Database migrations with Goose & Atlas
- 🐳 Fully dockerized dev environment
- 🧰 Linting, formatting, and code generation included

---

## 📦 Requirements

- Docker + Docker Compose
- Make
- Go 1.21+
- Optional: `buf`, `atlas`, `goose`, `golangci-lint`

---

## 🛠️ Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/go-password-manager.git
cd go-password-manager
```

### 2. 🔐 How to Generate TLS Certificates

TLS certificates are required for secure gRPC communication. You can generate a self-signed certificate for `localhost` using the following command:

```bash
make generate-cert
```

### 3. How to run Server

You can run server in Docker image or as a dinary

as a docker image with all images applied 
```bash
make run-docker 
```

as a binary 
```bash
make run-server
```

to debug with Delve

```bash
make run-docker-debug
```

### 3. How to run Client

to debug Client 

```bash
make debug-client
```

to make a clean run

```bash
make run-client
```

### 4. How to build image for all platforms

Simply use command in yur terminal

```bash
make build-client-all
```