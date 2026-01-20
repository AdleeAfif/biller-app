# Monthly Biller API

A REST API for managing monthly salary, commitments, and financial summaries built with Go, Gin, and MongoDB.

## Features

- 🔐 User authentication with JWT
- 💰 Salary management (default and monthly overrides)
- 📊 Commitment tracking (decimal and percentage-based)
- 📈 Monthly and yearly summaries
- 👨‍💼 Admin user management
- 🗄️ MongoDB for data persistence
- 🔒 Password hashing with bcrypt

## Quick Start

### Option 1: Using Docker (Recommended)

1. Run with Docker Compose:

```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Option 2: Local Development

1. Ensure MongoDB is running locally

2. Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

3. Install dependencies:

```bash
go mod download
```

4. Run the server:

```bash
go run cmd/server/main.go
```

Or using Make:

```bash
make run
```

See [QUICKSTART.md](QUICKSTART.md) for detailed setup instructions.

## API Documentation

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login and get JWT token

### User Profile & Salary

- `PUT /api/users/me/salary/default` - Set default monthly salary
- `PUT /api/users/me/salary/:year/:month` - Update salary for specific month

### Commitments

- `POST /api/users/me/commitments/default` - Set default monthly commitments
- `POST /api/users/me/commitments/:year/:month` - Set commitments for specific month
- `PATCH /api/users/me/commitments/:year/:month/:commitment_id` - Update commitment paid status

### Summaries

- `GET /api/users/me/summary/monthly/:year/:month` - Get monthly summary
- `GET /api/users/me/summary/yearly/:year` - Get yearly summary

### Admin (requires admin role)

- `GET /api/admin/users` - List all users
- `PUT /api/admin/users/:user_id` - Update user
- `DELETE /api/admin/users/:user_id` - Delete user

## Tech Stack

- Go
- Gin Web Framework
- MongoDB
- JWT Authentication
- bcrypt for password hashing
