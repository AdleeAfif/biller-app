# Monthly Biller API - Project Overview

## 🎯 What Is This?

A complete REST API backend for managing personal finances, specifically monthly salary and financial commitments. Built following clean architecture principles with Go, Gin framework, and MongoDB.

## ✨ Key Features

### Authentication & Authorization

- JWT-based authentication
- Role-based access control (User/Admin)
- Secure password hashing with bcrypt

### Financial Management

- **Default Salary**: Set a base monthly salary
- **Monthly Overrides**: Override salary for specific months
- **Default Commitments**: Set recurring monthly commitments
- **Monthly Commitments**: Override or add commitments for specific months
- **Commitment Types**:
  - Decimal: Fixed amount (e.g., $1,200 rent)
  - Percentage: Based on salary (e.g., 20% savings)
- **Payment Tracking**: Mark commitments as paid/unpaid

### Reports & Analytics

- Monthly summaries with commitment breakdown
- Yearly summaries with monthly breakdown
- Automatic calculation of remaining balance

### Admin Features

- List all users
- Update user details
- Soft delete users

## 🏗️ Architecture

### Clean Architecture Layers

```
cmd/server/          → Entry point
internal/
  ├── auth/          → Authentication logic
  ├── user/          → User management
  ├── commitment/    → Commitment handling
  ├── summary/       → Report generation
  ├── admin/         → Admin operations
  ├── middleware/    → Auth middleware
  ├── models/        → Data models & DTOs
  └── config/        → Configuration
pkg/
  ├── db/            → Database connection
  └── jwt/           → JWT utilities
```

### Technology Stack

| Component          | Technology           |
| ------------------ | -------------------- |
| Language           | Go 1.21+             |
| Web Framework      | Gin                  |
| Database           | MongoDB              |
| Authentication     | JWT (golang-jwt/jwt) |
| Password Hashing   | bcrypt               |
| Environment Config | godotenv             |

## 📊 Data Model

### Collections

1. **users** - User accounts
2. **monthly_records** - Monthly financial records
3. **yearly_summaries** - Yearly aggregations
4. **default_commitments** - User's default commitments

### Indexes

- `users`: Unique on `username` and `email`
- `monthly_records`: Unique compound on `(user_id, year, month)`
- `yearly_summaries`: Unique compound on `(user_id, year)`
- `default_commitments`: Unique on `user_id`

## 🔐 Security

- Passwords are hashed using bcrypt before storage
- JWT tokens expire after 24 hours
- Authorization headers required for protected routes
- Admin-only endpoints protected by role-based middleware
- Soft delete for users (preserves data integrity)

## 🌐 API Design

### Authentication Flow

1. User registers → Password hashed → User created
2. User logs in → Credentials validated → JWT token issued
3. Subsequent requests → Token in Authorization header → User authenticated

### Resource Organization

- `/api/auth/*` - Public authentication endpoints
- `/api/users/me/*` - User's own resources
- `/api/admin/*` - Admin-only endpoints

### Response Patterns

- Success: HTTP 200/201 with data
- Client Error: HTTP 400/401/403/404 with error message
- Server Error: HTTP 500 with generic error

## 🚀 Deployment Options

### Docker

```bash
docker-compose up -d
```

### Local

```bash
make run
```

### Production Considerations

1. Change `JWT_SECRET` in environment
2. Use MongoDB Atlas or managed MongoDB
3. Add rate limiting
4. Enable HTTPS/TLS
5. Add logging and monitoring
6. Set up backup strategy

## 📝 API Endpoints Summary

| Method | Endpoint                                     | Auth  | Description              |
| ------ | -------------------------------------------- | ----- | ------------------------ |
| POST   | `/api/auth/register`                         | No    | Register user            |
| POST   | `/api/auth/login`                            | No    | Login user               |
| PUT    | `/api/users/me/salary/default`               | Yes   | Set default salary       |
| PUT    | `/api/users/me/salary/:year/:month`          | Yes   | Set monthly salary       |
| POST   | `/api/users/me/commitments/default`          | Yes   | Set default commitments  |
| POST   | `/api/users/me/commitments/:year/:month`     | Yes   | Set monthly commitments  |
| PATCH  | `/api/users/me/commitments/:year/:month/:id` | Yes   | Update commitment status |
| GET    | `/api/users/me/summary/monthly/:year/:month` | Yes   | Get monthly summary      |
| GET    | `/api/users/me/summary/yearly/:year`         | Yes   | Get yearly summary       |
| GET    | `/api/admin/users`                           | Admin | List all users           |
| PUT    | `/api/admin/users/:id`                       | Admin | Update user              |
| DELETE | `/api/admin/users/:id`                       | Admin | Delete user              |

## 🔄 Typical User Flow

1. **Registration**

   ```
   POST /api/auth/register
   → User account created with 'user' role
   ```

2. **Login**

   ```
   POST /api/auth/login
   → Receive JWT token
   ```

3. **Setup Finances**

   ```
   PUT /api/users/me/salary/default (salary: 5000)
   → Default salary set to $5,000

   POST /api/users/me/commitments/default
   → Add recurring commitments (rent, savings, etc.)
   ```

4. **Monthly Management**

   ```
   GET /api/users/me/summary/monthly/2026/1
   → View January 2026 summary

   PATCH /api/users/me/commitments/2026/1/{id}
   → Mark rent as paid
   ```

5. **Yearly Review**
   ```
   GET /api/users/me/summary/yearly/2026
   → View full year financial summary
   ```

## 🧪 Testing

Run tests:

```bash
make test
```

Manual API testing with curl examples in [API_EXAMPLES.md](API_EXAMPLES.md)

## 📚 Documentation Files

- **README.md** - Project overview and quick start
- **QUICKSTART.md** - Detailed setup instructions
- **API_EXAMPLES.md** - Complete API examples with curl
- **PROJECT_OVERVIEW.md** - This file
- **.env.example** - Environment configuration template

## 🛠️ Development Commands

```bash
make run      # Run the server
make build    # Build binary
make install  # Install dependencies
make test     # Run tests
make fmt      # Format code
make clean    # Clean build artifacts
```

## 🔮 Future Enhancements

- [ ] Budget categories and subcategories
- [ ] Recurring commitments scheduling
- [ ] Email notifications for unpaid commitments
- [ ] Export reports to PDF/Excel
- [ ] Multi-currency support
- [ ] Budget forecasting and alerts
- [ ] Integration with bank APIs
- [ ] Mobile app (React Native/Flutter)
- [ ] Real-time updates with WebSockets
- [ ] GraphQL API option

## 📄 License

This is a sample project for educational/portfolio purposes.

## 👤 Author

Built following REST API best practices and clean architecture principles.
