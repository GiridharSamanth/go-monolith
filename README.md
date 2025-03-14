# Go Monolith

A web service built with Go and the Gin framework.

## Prerequisites

- Go 1.21 or higher
- Git

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

3. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

- `GET /`: Welcome message

## Project Structure

```
go-monolith/
├── cmd/
│   └── server/
│       └── main.go    # Main application entry point
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
└── README.md         # Project documentation
```

## License

This project is licensed under the MIT License. 