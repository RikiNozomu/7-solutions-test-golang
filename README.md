# 7-solutions-test-golang

Project for applying to [7SOLUTIONS's Backend Golang Coding Test](https://github.com/7-solutions/backend-challenge)

## 🛠 Project Setup & Run Instructions

### Prerequisites
- Go 1.25.3
- MongoDB
- Docker

### Installation & Build
```bash
git clone https://github.com/RikiNozomu/7-solutions-test-golang.git
cd 7-solutions-test-golang
go mod download
go mod tidy
```

### Run the application
```bash
go run main.go
```
By default, the server listens on `localhost:8080`.

### Testing
```bash
go test ./adapters/handlers ./core/services
```

### Run the application with MongoDB
```bash
docker compose up -d
```

## 🔐 JWT Token Usage Guide

This project uses JWT for securing `PUT /user/:id` and `DELETE /user/:id`.

### Token Generation
- First, You must create a user via `POST /user`.
- Next, Get JWT `token` by logging in via `POST /auth/login`
- After successful login, the server returns a JWT `token`.
- The `token` must be included in the `Authorization` header for protected endpoints.

### Header Format
```
Authorization: Bearer <your_jwt_token>
```

### Middleware
- JWT validation is handled in `middlewares/authentication.go`.
- Invalid or expired tokens will result in a `401 Unauthorized` response.

## 📬 List API Requests & Responses

### 1. Create a User
**POST** `/user`

**Body**
```json
{
    "email" : "email@test.com", // email format only
    "name" : "name123456", // min=8, max=24 
    "password" : "pass1234WORD" //  min=8, max=24, no space, lower>=1, upper>=1, number>=1
}
```
**Response** `201 - Success`
```json
{
    "data": {
        "id": "68fdaa1a07efb9324b5de04a",
        "name": "test-another",
        "email": "rikitphorn@ggg.com",
        "created_at": "2025-10-26T04:56:58.934368922Z"
    }
}
```
**Response** `400 - Bad Request`
```json
{
    "errors": [
        "Name must be at least 8 characters long."
    ]
}
```
**Response** `409 - Conflict`
```json
{
    "errors": [
        "Cannot add user with exist email."
    ]
}
```
### 2. Get all Users
**GET** `/user`

**Response** `200 - Success`
```json
{
    "data": [
        {
            "id": "68fd0e2407efb9324b5de047",
            "name": "4564564456546",
            "email": "rwerewrikit@ggg.com",
            "created_at": "2025-10-25T17:51:32.507Z"
        },
        {
            "id": "68fdaa1a07efb9324b5de04a",
            "name": "test-another",
            "email": "rikitphorn@ggg.com",
            "created_at": "2025-10-26T04:56:58.934Z"
        }
    ]
}
```
### 3. Get a User
**GET** `/user/:id`

**Response** `200 - Success`
```json
{
    "data": {
        "id": "68fd0e2407efb9324b5de047",
        "name": "4564564456546",
        "email": "rwerewrikit@ggg.com",
        "created_at": "2025-10-25T17:51:32.507Z"
    }
}
```
**Response** `404 - Not Found`
```json
{
    "errors": [
        "User not found."
    ]
}
```
### 4. Login a User
**POST** `/auth/login`

**Body**
```json
{
    "email" : "rikitphorn@ggg.com",
    "password" : "pAass1234"
}
```
**Response** `200 - Success`
```json
{
    "data": {
        "expire_in_unix": 1761458524, // expire in 1 hour
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...." // JWT string token
    }
}
```
**Response** `401 - Unauthenticated`
```json
{
    "errors": [
        "Unauthenticated"
    ]
}
```
**Response** `404 - (User) Not found`
```json
{
    "errors": [
        "User not found."
    ]
}
```
### 5. Update a User
**PUT** `/user/:id`

**Headers**
```json
{
    "Authorization" : "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....",
}
```
**Body**
```json
{
    "email" : "email@test.com", // email format only
    "name" : "name123456", // min=8, max=24 
}
```
**Response** `200 - Success`
```json
{
    "data": {
        "id": "68fdaa1a07efb9324b5de04a",
        "name": "7777777777777",
        "email": "rikitphorn@ggg.com",
        "created_at": "2025-10-26T04:56:58.934Z"
    }
}
```
**Response** `400 - Bad Request`
```json
{
    "errors": [
        "Name must be at least 8 characters long."
    ]
}
```
**Response** `401 - Unauthenticated`
```json
{
    "errors": [
        "Missing barrier token"
    ]
}
```
**Response** `404 - (User) not found`
```json
{
    "errors": [
        "User not found."
    ]
}
```

### 6. Delete a User
**DELETE** `/user/:id`

**Headers**
```json
{
    "Authorization" : "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....",
}
```
**Response** `200 - Success`
```json
{
    "message": "User has been removed."
}
```
**Response** `401 - Unauthenticated`
```json
{
    "errors": [
        "Unauthenticated"
    ]
}
```
**Response** `404 - (User) not found`
```json
{
    "errors": [
        "User not found."
    ]
}
```

## 📌 Assumptions & Decisions Made

- **Architecture**: Choose Hexagonal architecture, that is easy to implement services, routers, and unit testing.
- **Security**:
  - Decided to use JWT as token for protected routes.
  - Protected routes have only 2 endpoint; `PUT /user/:id` and `DELETE /user/:id`. Because I didn't want to make any superuser before.
  - Only the owner account that can update and delete themselves.
- **Error Handling**: Built `middlewares/error-response.go` as middleware to manage showing error instead of default.
- **Validation**: Built `utils/validate.go` as utility to manage validations instead of default.
- **Logging**: Built `middlewares/log-request.go` for request logging, and go function in `main.go` for showing amount of user.
