# Go Secrets API

A lightweight API for temporarily storing secrets using a single token-based authentication mechanism. The service utilizes Redis for storage.

> **Disclaimer:** This project was created for learning Go. If you find practical use cases for it, let me know by [opening an issue](https://github.com/Draugelis/go-secrets/issues).

## Features
- Temporary secret storage
- Token-based authentication
- Redis as the backend
- Simple API interface with Swagger documentation

## 🚀 Getting Started

### Using Task (Recommended)

For easier management, the project includes a `Taskfile.yml`. To start everything up, simply run:

```sh
task start
```

This will spin up the API and Redis instance using Docker.

### Running Without Task

#### Using Docker
Ensure you have Docker installed, then run:

```sh
docker-compose -f docker/docker-compose.yml up  -d
```

#### Running Locally
If you prefer to run the application without Docker, ensure Redis is installed and running, then execute:

```sh
go run main.go
```

## ⚙️ Configuration

Set the following environment variables:

```sh
export REDIS_URL=redis://localhost:6379
export APP_PORT=8888
```

Alternatively, these can be defined in a `.env` file.

## 📡 API Usage

Swagger documentation is available at:

```sh
http://localhost:8888/swagger/index.html
```

### Endpoints

#### 🔑 Token Management
- `GET /token?ttl={ttl}` - Generates a short-lived token
- `DELETE /token` - Invalidates the token and associated secrets
- `GET /token/valid` - Checks if the token is still valid

#### 🔐 Secret Management
- `POST /secret/{key}` - Stores a secret
- `GET /secret/{key}` - Retrieves a secret
- `DELETE /secret/{key}` - Deletes a secret

## 🛠 Taskfile Usage

A `Taskfile.yml` is included for easier project management. To see available tasks, run:

```sh
task
```

