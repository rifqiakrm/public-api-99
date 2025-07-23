# Public API

This service acts as the Public API Layer that integrates multiple internal microservices like `user-service` and `listing-service`. It exposes a unified interface for frontend or third-party consumers, handling request validation, error mapping, and aggregation logic.

## Features

* RESTful endpoints
* Integration with internal services via HTTP clients
* Request/response validation
* Graceful error handling
* Unit tested with mocks

## Folder Structure

```
public-api/
├── client/             # HTTP clients to other services
├── config/             # Project Config
├── handler/            # HTTP handlers (Gin)
├── model/              # Request/response & shared models
├── mocks/              # Auto-generated mocks (GoMock)
├── service/            # Business logic
├── router/             # Route registration
├── main.go             # App entry point
└── go.mod              # Dependencies
```

## Requirements

* Go 1.21+
* Gin Web Framework
* GoMock (for unit testing)

## Setup

```bash
go mod tidy
go run main.go
```

## Run Tests

```bash
go test ./...
```

## Generate Mocks (if needed)

```bash
go generate ./...
```

## Example Endpoints

### Create User

```
POST /api/v1/users
{
  "name": "John Doe"
}
```

### Create Listings

```
POST /api/v1/listings
{
  "user_id": 8,
  "listing_type": "sale",
  "price": 10000000
}
```

### Get Listings

```
GET /api/v1/listings?page=1&limit=10
```

## 🔖 Author
Rifqi Fauzan Akram  
Email: rifqiakram57@gmail.com  
GitHub: [@rifqiakrm](https://github.com/rifqiakrm)

---
