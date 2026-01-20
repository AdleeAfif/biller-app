# API Endpoints Reference

## Base URL

```
http://localhost:8080
```

## Authentication

### Register

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "username": "john",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "password123"
  }'
```

Response:

```json
{
  "access_token": "your-jwt-token"
}
```

## User Profile & Salary

### Set Default Salary

```bash
curl -X PUT http://localhost:8080/api/users/me/salary/default \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "salary": 5000
  }'
```

### Update Monthly Salary

```bash
curl -X PUT http://localhost:8080/api/users/me/salary/2026/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "salary": 6000
  }'
```

## Commitments

### Set Default Commitments

```bash
curl -X POST http://localhost:8080/api/users/me/commitments/default \
  -H "Authorization: Bearer YOUR_TOKEN" \
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

### Set Monthly Commitments

```bash
curl -X POST http://localhost:8080/api/users/me/commitments/2026/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "commitments": [
      {
        "name": "Car Loan",
        "type": "decimal",
        "value": 800
      }
    ]
  }'
```

### Update Commitment Paid Status

```bash
curl -X PATCH http://localhost:8080/api/users/me/commitments/2026/1/COMMITMENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_paid": true
  }'
```

## Summaries

### Get Monthly Summary

```bash
curl -X GET http://localhost:8080/api/users/me/summary/monthly/2026/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:

```json
{
  "salary": 5000,
  "total_commitment": 2500,
  "remaining_balance": 2500,
  "commitments": [
    {
      "id": "...",
      "name": "Rent",
      "amount": 1200,
      "is_paid": true
    }
  ]
}
```

### Get Yearly Summary

```bash
curl -X GET http://localhost:8080/api/users/me/summary/yearly/2026 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:

```json
{
  "year": 2026,
  "total_salary": 60000,
  "total_commitment": 30000,
  "total_remaining": 30000,
  "monthly_breakdown": [
    {
      "month": 1,
      "remaining": 2500
    }
  ]
}
```

## Admin Endpoints

### List All Users

```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/admin/users/USER_ID \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com",
    "default_salary": 7000
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/admin/users/USER_ID \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## Health Check

```bash
curl -X GET http://localhost:8080/health
```
