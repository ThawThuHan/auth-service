# Auth Service

This is the **Auth Service** for the API Gateway Lab project. It is a microservice written in Go, responsible for authentication and authorization, designed to work in a microservices architecture.

## Features
- JWT-based authentication
- User login and jwks endpoints
- gRPC client integration
- Configurable via environment variables
- Docker and Kubernetes (Helm) deployment support

## Prerequisites
- [Go](https://golang.org/dl/) 1.25 or newer
- [Docker](https://www.docker.com/)
- [buf](https://docs.buf.build/installation) (for protobuf generation)
- [Helm](https://helm.sh/) (for Kubernetes deployment)

## Getting Started

### 1. Clone the Repository
```bash
git clone <your-repo-url>
cd auth-service
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Set Up Environment Variables
Copy the example environment file and edit as needed:
```bash
cp .env.example .env
# Edit .env to match your configuration
```

### 4. Run the Service Locally
```bash
go run ./cmd/main.go
```

### 5. Run with Docker
Build and run the Docker container:
```bash
docker build -t auth-service .
docker run --env-file .env -p 8080:8080 auth-service
```

### 6. Run with Air (Hot Reload for Development)
If you have [Air](https://github.com/cosmtrek/air) installed:
```bash
air
```

### 7. Generate Protobuf Files
If you modify proto files, regenerate Go code with:
```bash
buf generate
```

### 8. Deploy with Helm (Kubernetes)
Update values in `helm-chart/values.yaml` as needed, then:
```bash
helm install auth-service ./helm-chart
```

## Project Structure
```
auth-service/
├── cmd/                # Main application entrypoint
├── internal/           # Application core logic
│   ├── config/         # Configuration and logger
│   ├── helper/         # Helper utilities
│   ├── infrastructures/# Database and gRPC client
│   ├── modules/        # Auth module (handler, service, repo, dto)
│   ├── routes/         # Route definitions
│   └── util/           # Utilities (JWT, HTTP response)
├── proto/              # Protobuf definitions and generated code
├── helm-chart/         # Helm chart for Kubernetes deployment
├── Dockerfile          # Docker build file
├── .env.example        # Example environment variables
└── README.md           # This file
```

## Useful Commands
- **Run tests:**
	```bash
	go test ./...
	```
- **Lint code:**
	```bash
	golangci-lint run
	```

## Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the MIT License.
