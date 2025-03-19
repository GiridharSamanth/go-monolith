# Go Monolith

A web service built with Go and the Gin framework, following clean architecture principles.

## Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional, for running services in containers)

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd go-monolith
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env.local
# Edit .env.local with your configuration
```

4. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## Project Structure

```
go-monolith/
├── cmd/
│   └── server/
│       └── main.go    # Main application entry point
├── internal/
│   ├── app/          # Application core and configuration
│   ├── bff/          # Backend-For-Frontend layer
│   └── modules/      # Business modules and domain logic
├── pkg/
│   ├── auth/         # Authentication and authorization
│   ├── context/      # Context utilities
│   ├── errors/       # Error handling and custom errors
│   ├── logger/       # Logging utilities
│   └── metrics/      # Metrics and monitoring
├── .env.local        # Local environment variables
├── .gitignore        # Git ignore rules
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
└── README.md        # Project documentation
```

## Architecture

The project follows clean architecture principles with the following layers:

- **cmd/server**: Application entry point and server setup
- **internal/app**: Core application configuration and setup
- **internal/bff**: Backend-For-Frontend layer for API composition
- **internal/modules**: Business modules containing domain logic
- **pkg**: Shared packages and utilities

## Development

### Running Tests
```bash
go test ./...
```

### Code Style
The project follows standard Go formatting. Format your code before committing:
```bash
go fmt ./...
```

### Environment Variables
The application uses environment variables for configuration. Copy `.env.example` to `.env.local` and adjust the values as needed.


### Running the Server
```bash
go run cmd/server/main.go
```

### Example Curl Commands

Get story by ID (v2.0):
```bash
curl -H 'access-token:550e8400-e29b-41d4-a716-446655440000' 'http://localhost:8080/v2.0/stories/10' | jq
```

Get story by ID (v1.2):
```bash
curl -H 'access-token:550e8400-e29b-41d4-a716-446655440000' 'http://localhost:8080/v1.2/stories?id=10' | jq
```

## License

This project is licensed under the MIT License. 