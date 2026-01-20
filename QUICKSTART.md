# Quick Start Guide

## Prerequisites

1. **Go** (version 1.21 or higher)
2. **MongoDB** (running locally or remote)

## Installation Steps

### 1. Install Dependencies

```bash
go mod download
```

Or use the Makefile:

```bash
make install
```

### 2. Configure Environment

Copy the example environment file and update with your settings:

```bash
cp .env.example .env
```

Edit `.env` with your MongoDB connection details:

```env
MONGODB_URI=mongodb://localhost:27017
DATABASE_NAME=biller_app
JWT_SECRET=your-super-secret-key-change-this
PORT=8080
```

### 3. Start MongoDB

Make sure MongoDB is running. If using Docker:

```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### 4. Run the Server

```bash
go run cmd/server/main.go
```

Or using Makefile:

```bash
make run
```

The server will start on `http://localhost:8080`

## Testing the API

### 1. Register a User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

Save the `access_token` from the response.

### 3. Set Default Salary

```bash
curl -X PUT http://localhost:8080/api/users/me/salary/default \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "salary": 5000
  }'
```

### 4. Set Default Commitments

```bash
curl -X POST http://localhost:8080/api/users/me/commitments/default \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "commitments": [
      {
        "name": "Rent",
        "type": "decimal",
        "value": 1200
      },
      {
        "name": "Savings",
        "type": "percentage",
        "value": 20
      }
    ]
  }'
```

### 5. Get Monthly Summary

```bash
curl -X GET http://localhost:8080/api/users/me/summary/monthly/2026/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Project Structure

```
biller-app/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── admin/
│   │   └── handler.go           # Admin handlers
│   ├── auth/
│   │   └── handler.go           # Authentication handlers
│   ├── commitment/
│   │   ├── handler.go           # Commitment handlers
│   │   └── utils.go
│   ├── config/
│   │   └── config.go            # Configuration
│   ├── middleware/
│   │   └── auth.go              # JWT middleware
│   ├── models/
│   │   ├── dto.go               # Request/Response DTOs
│   │   └── models.go            # Database models
│   ├── summary/
│   │   └── handler.go           # Summary handlers
│   └── user/
│       ├── handler.go           # User handlers
│       └── utils.go
└── pkg/
    ├── db/
    │   └── mongodb.go           # MongoDB connection
    └── jwt/
        └── jwt.go               # JWT utilities
```

## Creating an Admin User

To create an admin user, you'll need to manually update the user role in MongoDB:

```javascript
// Connect to MongoDB
use biller_app

// Update user role to admin
db.users.updateOne(
  { username: "testuser" },
  { $set: { role: "admin" } }
)
```

Or using mongosh:

```bash
mongosh
use biller_app
db.users.updateOne({username: "testuser"}, {$set: {role: "admin"}})
```

## Common Commands

- **Run server**: `make run` or `go run cmd/server/main.go`
- **Build**: `make build`
- **Install dependencies**: `make install`
- **Format code**: `make fmt`
- **Clean build**: `make clean`

## Troubleshooting

### MongoDB Connection Error

- Ensure MongoDB is running: `mongosh` or check Docker container
- Verify `MONGODB_URI` in `.env` file

### Port Already in Use

- Change `PORT` in `.env` file
- Or kill the process using port 8080

### JWT Token Issues

- Ensure you're including the token in the Authorization header
- Format: `Authorization: Bearer YOUR_TOKEN`
- Check that `JWT_SECRET` is set in `.env`

## Next Steps

1. Explore all API endpoints in [API_EXAMPLES.md](API_EXAMPLES.md)
2. Implement frontend client
3. Add more features like budget analytics
4. Deploy to production

For more details, see the complete [README.md](README.md)
