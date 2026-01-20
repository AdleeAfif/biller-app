# Development Guide

## Project Structure

```
biller-app/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/                           # Private application code
│   ├── admin/
│   │   └── handler.go                 # Admin endpoints
│   ├── auth/
│   │   └── handler.go                 # Authentication
│   ├── commitment/
│   │   ├── handler.go                 # Commitment management
│   │   └── utils.go                   # Helper functions
│   ├── config/
│   │   └── config.go                  # Configuration loader
│   ├── middleware/
│   │   └── auth.go                    # JWT middleware
│   ├── models/
│   │   ├── dto.go                     # Request/Response DTOs
│   │   ├── models.go                  # Database models
│   │   └── models_test.go             # Model tests
│   ├── summary/
│   │   └── handler.go                 # Summary/reports
│   └── user/
│       ├── handler.go                 # User management
│       └── utils.go                   # Helper functions
├── pkg/                                # Public reusable packages
│   ├── db/
│   │   └── mongodb.go                 # Database connection
│   └── jwt/
│       └── jwt.go                     # JWT utilities
├── .env.example                        # Environment template
├── .gitignore
├── API_EXAMPLES.md                     # API usage examples
├── Biller-API.postman_collection.json # Postman collection
├── docker-compose.yml                  # Docker setup
├── Dockerfile                          # Container image
├── go.mod                              # Go dependencies
├── go.sum                              # Dependency checksums
├── Makefile                            # Build commands
├── PROJECT_OVERVIEW.md                 # Architecture docs
├── QUICKSTART.md                       # Setup guide
└── README.md                           # Main documentation
```

## Adding New Features

### 1. Adding a New Endpoint

Example: Adding a "Get User Profile" endpoint

**Step 1**: Add the handler in `internal/user/handler.go`:

```go
func (h *Handler) GetProfile(c *gin.Context) {
    userID, err := middleware.GetUserID(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
        return
    }

    var user models.User
    err = h.db.Collection("users").FindOne(
        context.Background(),
        bson.M{"_id": userID},
    ).Decode(&user)

    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to get profile"})
        return
    }

    c.JSON(http.StatusOK, user)
}
```

**Step 2**: Register route in `cmd/server/main.go`:

```go
protected.GET("/profile", userHandler.GetProfile)
```

### 2. Adding a New Model

**Step 1**: Define in `internal/models/models.go`:

```go
type Budget struct {
    ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
    Name   string             `json:"name" bson:"name"`
    Amount float64            `json:"amount" bson:"amount"`
}
```

**Step 2**: Add DTO in `internal/models/dto.go`:

```go
type CreateBudgetRequest struct {
    Name   string  `json:"name" binding:"required"`
    Amount float64 `json:"amount" binding:"required,gt=0"`
}
```

### 3. Adding Database Indexes

Add to `pkg/db/mongodb.go` in the `createIndexes` function:

```go
budgetIndexes := []mongo.IndexModel{
    {
        Keys: bson.D{{Key: "user_id", Value: 1}},
    },
}
if _, err := db.Collection("budgets").Indexes().CreateMany(ctx, budgetIndexes); err != nil {
    return err
}
```

## Code Style Guidelines

### Naming Conventions

- **Files**: lowercase with underscores (e.g., `user_handler.go`)
- **Packages**: lowercase, single word (e.g., `auth`, `user`)
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase
- **Constants**: PascalCase or UPPER_SNAKE_CASE

### Error Handling

Always return meaningful errors:

```go
// Good
if err != nil {
    c.JSON(http.StatusInternalServerError, models.ErrorResponse{
        Error: "failed to create record",
    })
    return
}

// Bad
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()}) // Don't expose internal errors
    return
}
```

### HTTP Status Codes

- `200 OK`: Success (GET, PATCH)
- `201 Created`: Resource created (POST)
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing/invalid token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

## Testing

### Unit Tests

Create test files alongside source files with `_test.go` suffix:

```go
// internal/models/models_test.go
package models

import "testing"

func TestUserRole(t *testing.T) {
    if RoleUser != "user" {
        t.Errorf("Expected 'user', got '%s'", RoleUser)
    }
}
```

Run tests:

```bash
go test ./...
# or
make test
```

### Integration Tests

Test handlers with a test database:

```go
func TestRegister(t *testing.T) {
    // Setup test DB
    // Create handler
    // Make request
    // Assert response
}
```

## Database Operations

### Best Practices

1. Always use context with timeout
2. Index frequently queried fields
3. Use projection to limit returned fields
4. Handle `mongo.ErrNoDocuments` explicitly

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var user models.User
err := db.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
if err == mongo.ErrNoDocuments {
    // Handle not found
} else if err != nil {
    // Handle other errors
}
```

## Environment Variables

Add new variables to:

1. `.env.example` (with placeholder values)
2. `internal/config/config.go` (with defaults)
3. `docker-compose.yml` (if using Docker)
4. Documentation

## Debugging

### Enable Gin Debug Mode

```go
gin.SetMode(gin.DebugMode) // In main.go
```

### Log Database Queries

Enable MongoDB query logging by setting log level:

```go
clientOptions := options.Client().
    ApplyURI(uri).
    SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
```

### Check JWT Token

Decode at https://jwt.io to inspect claims

## Common Issues

### "Duplicate Key Error"

- Check if unique index violated
- Ensure proper upsert logic

### "Unauthorized" Error

- Verify token in Authorization header
- Check token expiration
- Ensure JWT_SECRET matches

### MongoDB Connection Failed

- Verify MongoDB is running
- Check MONGODB_URI in .env
- Ensure network connectivity

## Pull Request Process

1. Create feature branch from `main`
2. Make changes with clear commit messages
3. Add/update tests
4. Update documentation
5. Run `make fmt` and `make test`
6. Create PR with description

## Useful Commands

```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run linter (if installed)
golangci-lint run

# View dependencies
go list -m all

# Update dependencies
go get -u ./...
go mod tidy

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go
```

## Resources

- [Gin Documentation](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
