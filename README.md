# Plantheon Backend API

<img width="200" height="200" alt="logo" src="https://github.com/user-attachments/assets/92e4733b-0978-4c3e-9bf0-d7bb6e7068a7" />

> A comprehensive RESTful API for plant disease management and agricultural activity tracking, built with Go and the Gin framework.

[![Go Version](https://img.shields.io/badge/Go-1.23.0-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Database Setup](#database-setup)
  - [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
  - [Authentication](#authentication)
  - [User Management](#user-management)
  - [Disease Management](#disease-management)
  - [Activity Tracking](#activity-tracking)
  - [Scan History](#scan-history)
  - [Social Features](#social-features)
  - [Plant Guides](#plant-guides)
  - [Complaints & Analytics](#complaints--analytics)
- [Data Models](#data-models)
- [Security](#security)
- [Deployment](#deployment)
- [Contributing](#contributing)

## Overview

Plantheon Backend API is a robust server-side application designed to support agricultural management and plant disease detection systems. It provides comprehensive endpoints for managing plant diseases, tracking farming activities, handling user-generated content, and analyzing disease prediction accuracy through a complaint and verification system.

The API serves as the backend for the Plantheon mobile application, supporting features such as:
- AI-powered plant disease detection with user feedback mechanisms
- Agricultural activity and financial tracking
- Social community features for farmers
- Plant cultivation guides with stage-based instructions
- News and educational content management

## Features

### Core Functionality
- ✅ **JWT-based Authentication** with role-based access control (User/Admin)
- ✅ **Disease Management** with comprehensive CRUD operations
- ✅ **Activity Tracking** with financial analytics (monthly, annual, multi-year)
- ✅ **Scan History** for disease detection results
- ✅ **Social Platform** with posts, comments, likes, and user profiles
- ✅ **Plant Cultivation Guides** with multi-stage instructions
- ✅ **News & Blog Management** with tagging system
- ✅ **Complaint System** for AI prediction feedback
- ✅ **Analytics Dashboard** for disease trends and user contributions

### Technical Features
- ✅ Pagination, search, and filtering on all list endpoints
- ✅ Input validation and sanitization
- ✅ Password hashing with bcrypt
- ✅ CORS support with configurable origins
- ✅ PostgreSQL database with GORM ORM
- ✅ Email notifications (OTP for password reset)
- ✅ Excel import/export functionality
- ✅ Soft delete support for content moderation

## Technology Stack

- **Language**: Go 1.23.0
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **Excel Processing**: excelize
- **Environment Management**: godotenv

## Project Structure

```
plantheon-backend/
├── main.go                          # Application entry point
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── .env                             # Environment configuration
├── common/
│   ├── database.go                  # Database connection manager
│   ├── utils.go                     # JWT & password utilities
│   └── email_service.go             # Email service for OTP
├── models/
│   ├── users/                       # User authentication & management
│   │   ├── models.go
│   │   ├── serializers.go
│   │   ├── routers.go
│   │   ├── middlewares.go
│   │   └── validators.go
│   ├── diseases/                    # Plant disease information
│   ├── activities/                  # Agricultural activity tracking
│   ├── activity_keywords/           # Activity categorization
│   ├── disease_activity_keywords/   # Disease-activity relationships
│   ├── scan_history/                # Disease scan results
│   ├── posts/                       # Social posts
│   ├── comments/                    # Post comments
│   ├── post_likes/                  # Post engagement
│   ├── blogs/                       # News articles
│   ├── blog_tags/                   # Blog categorization
│   ├── plants/                      # Plant information
│   ├── guide_stages/                # Cultivation guide stages
│   ├── sub_guide_stages/            # Detailed stage instructions
│   ├── complaints/                  # AI prediction feedback
│   └── noti/                        # User notifications
└── uploads/                         # File upload directory
```

Each model directory follows a consistent structure:
- `models.go` - Database models and ORM operations
- `serializers.go` - Request/response DTOs
- `routers.go` - HTTP handlers
- `validators.go` - Input validation
- `service.go` - Business logic (where applicable)

## Getting Started

### Prerequisites

- Go 1.23.0 or higher
- PostgreSQL 12 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/KKuyen/SE505_Plantheon_BackEnd.git
cd SE505_Plantheon_BackEnd
```

2. Install dependencies:
```bash
go mod download
```

### Configuration

Create a `.env` file in the project root:

```env
# Database Configuration
DATABASE_URL=postgresql://username:password@host:5432/database_name?sslmode=require

# Server Configuration
PORT=3000

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# CORS Configuration (optional)
ALLOWED_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com

# Email Configuration (for password reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@plantheon.com
```

**Example with Supabase:**
```env
DATABASE_URL=postgresql://postgres:your_password@db.xxxxx.supabase.co:5432/postgres
PORT=3000
JWT_SECRET=plantheon-jwt-secret-2024-change-this-in-production
```

### Database Setup

The application uses GORM's auto-migration feature. On first run, all necessary tables will be created automatically. If you need to run migrations manually or have an existing database:

```bash
# The application will auto-migrate on startup
go run main.go
```

For production deployments, consider using a migration tool like [golang-migrate](https://github.com/golang-migrate/migrate) for better version control.

### Running the Application

**Development:**
```bash
go run main.go
```

**Production Build:**
```bash
go build -o plantheon-backend
./plantheon-backend
```

The server will start on the configured port (default: 3000). You can verify it's running:
```bash
curl http://localhost:3000/api/v1/health
```

Expected response:
```json
{
  "status": "OK",
  "message": "Plantheon Backend API is running"
}
```

## API Documentation

All API endpoints are prefixed with `/api/v1`. Below is a comprehensive overview of available endpoints.

### Authentication

#### Register a New User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### Forgot Password
```http
POST /api/v1/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

#### Verify OTP
```http
POST /api/v1/auth/verify-otp
Content-Type: application/json

{
  "email": "user@example.com",
  "otp": "123456"
}
```

#### Reset Password
```http
POST /api/v1/auth/reset-password
Content-Type: application/json

{
  "email": "user@example.com",
  "otp": "123456",
  "new_password": "NewSecurePass123!"
}
```

### User Management

#### Get Current User Profile
```http
GET /api/v1/users/profile
Authorization: Bearer <token>
```

#### Update Profile
```http
PUT /api/v1/users/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "newusername",
  "full_name": "New Name",
  "avatar": "https://example.com/avatar.jpg"
}
```

#### Get Public User Profile (with posts)
```http
GET /api/v1/public/users/:userId
```

#### Admin: Get All Users
```http
GET /api/v1/admin/users?page=1&limit=20
Authorization: Bearer <admin_token>
```

#### Admin: Update User
```http
PUT /api/v1/admin/users/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "role": "admin",
  "full_name": "Updated Name"
}
```

#### Admin: Disable/Enable User
```http
PATCH /api/v1/admin/users/:id/disable
PATCH /api/v1/admin/users/:id/enable
Authorization: Bearer <admin_token>
```

### Disease Management

#### Get Diseases (with pagination, search, filter)
```http
GET /api/v1/diseases?page=1&limit=10&type=fungal&search=rust
```

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `type` - Filter by disease type
- `search` - Search in name and description

#### Get All Diseases (no pagination)
```http
GET /api/v1/diseases/all
```

#### Get Disease Count
```http
GET /api/v1/diseases/count
```

#### Get Disease by Class Name
```http
GET /api/v1/diseases/:className
```

#### Admin: Create Disease
```http
POST /api/v1/diseases
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Tomato Late Blight",
  "class_name": "Tomato___Late_blight",
  "type": "fungal",
  "description": "A devastating disease caused by Phytophthora infestans...",
  "solution": [
    "Remove and destroy infected plants",
    "Apply copper-based fungicides",
    "Improve air circulation"
  ],
  "image_link": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ]
}
```

#### Admin: Update Disease
```http
PUT /api/v1/diseases/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Updated Disease Name",
  "description": "Updated description"
}
```

#### Admin: Import Diseases from Excel
```http
POST /api/v1/diseases/import-excel
Authorization: Bearer <admin_token>
Content-Type: multipart/form-data

file: <excel_file>
```

#### Admin: Delete Disease
```http
DELETE /api/v1/diseases/:className
Authorization: Bearer <admin_token>
```

### Activity Tracking

#### Get Activities
```http
GET /api/v1/activities?page=1&limit=10&user_id=uuid
```

#### Get All Activities
```http
GET /api/v1/activities/all?user_id=uuid
```

#### Get Activities by Day
```http
GET /api/v1/activities/by-day?date=2024-01-15&user_id=uuid
```

#### Get Activities Calendar by Month
```http
GET /api/v1/activities/get-activites-by-month?year=2024&month=1&user_id=uuid
```

#### Get Financial Summary (Monthly)
```http
GET /api/v1/activities/financial/monthly?year=2024&month=1&user_id=uuid
```

#### Get Financial Summary (Annual)
```http
GET /api/v1/activities/financial/annual?year=2024&user_id=uuid
```

#### Get Financial Summary (Multi-Year)
```http
GET /api/v1/activities/financial/multi-year?start_year=2022&end_year=2024&user_id=uuid
```

#### Create Activity
```http
POST /api/v1/activities
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Fertilizer Application",
  "description": "Applied NPK fertilizer to tomato field",
  "activity_date": "2024-01-15T10:00:00Z",
  "cost": 150.50,
  "revenue": 0,
  "location": "Field A",
  "keywords": ["fertilizer", "tomato"]
}
```

#### Update Activity
```http
PUT /api/v1/activities/:id
Authorization: Bearer <token>
Content-Type: application/json
```

#### Delete Activity
```http
DELETE /api/v1/activities/:id
Authorization: Bearer <token>
```

#### Delete Multiple Activities
```http
DELETE /api/v1/activities
Authorization: Bearer <token>
Content-Type: application/json

{
  "ids": ["uuid1", "uuid2", "uuid3"]
}
```

### Activity Keywords

#### Get All Activity Keywords
```http
GET /api/v1/activity-keywords
```

#### Search Activity Keywords
```http
GET /api/v1/activity-keywords/search?q=fertilizer
```

#### Get Activity Keywords for Disease
```http
GET /api/v1/activity-keywords/disease/:disease_id
```

#### Admin: Create Activity Keyword
```http
POST /api/v1/activity-keywords
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Organic Fertilizer",
  "description": "Keywords related to organic fertilization"
}
```

#### Admin: Import Activity Keywords from CSV
```http
POST /api/v1/disease-activity-keywords/import-csv
Authorization: Bearer <admin_token>
Content-Type: multipart/form-data

file: <csv_file>
```

### Scan History

#### Get User's Scan History
```http
GET /api/v1/scan-history?page=1&limit=20
Authorization: Bearer <token>
```

#### Get Scan History by ID
```http
GET /api/v1/scan-history/:id
Authorization: Bearer <token>
```

#### Create Scan History
```http
POST /api/v1/scan-history
Authorization: Bearer <token>
Content-Type: application/json

{
  "disease_id": "uuid",
  "image_url": "https://example.com/scan-image.jpg",
  "confidence_score": 0.95,
  "scan_date": "2024-01-15T10:00:00Z"
}
```

#### Delete Scan History
```http
DELETE /api/v1/scan-history/:id
Authorization: Bearer <token>
```

#### Delete All Scan History
```http
DELETE /api/v1/scan-history
Authorization: Bearer <token>
```

### Social Features

#### Get Posts
```http
GET /api/v1/posts?page=1&limit=10
Authorization: Bearer <token>
```

#### Search Posts
```http
GET /api/v1/posts/search?q=tomato&page=1&limit=10
Authorization: Bearer <token>
```

#### Get My Posts
```http
GET /api/v1/posts/my?page=1&limit=10
Authorization: Bearer <token>
```

#### Get Post by ID
```http
GET /api/v1/posts/:id
Authorization: Bearer <token>
```

#### Create Post
```http
POST /api/v1/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Just harvested my first tomatoes!",
  "images": ["https://example.com/image1.jpg"],
  "visibility": "public"
}
```

#### Update Post
```http
PUT /api/v1/posts/:id
Authorization: Bearer <token>
Content-Type: application/json
```

#### Delete Post
```http
DELETE /api/v1/posts/:id
Authorization: Bearer <token>
```

#### Like/Unlike Post
```http
PUT /api/v1/posts/:id/like
PUT /api/v1/posts/:id/unlike
Authorization: Bearer <token>
```

#### Share Post
```http
PUT /api/v1/posts/:id/share
Authorization: Bearer <token>
```

#### Get Comments
```http
GET /api/v1/posts/:id/comments
Authorization: Bearer <token>
```

#### Add Comment
```http
POST /api/v1/posts/:id/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Great harvest!"
}
```

#### Like/Unlike Comment
```http
PUT /api/v1/posts/comments/:commentId/like
PUT /api/v1/posts/comments/:commentId/unlike
Authorization: Bearer <token>
```

### Plant Guides

#### Get All Plants
```http
GET /api/v1/plants
```

#### Get Plant by ID
```http
GET /api/v1/plants/:id
```

#### Get Guide Stages for Plant
```http
GET /api/v1/guide-stages/plant/:plant_id
```

#### Get Sub-Guide Stages
```http
GET /api/v1/sub-guide-stages/guide-stage/:guide_stage_id
```

#### Admin: Create/Update/Delete Plants and Guides
```http
POST /api/v1/plants
PUT /api/v1/plants/:id
DELETE /api/v1/plants/:id
Authorization: Bearer <admin_token>
```

### News & Blogs

#### Get News
```http
GET /api/v1/news?page=1&limit=10&tag=farming
Authorization: Bearer <token>
```

#### Get News by ID
```http
GET /api/v1/news/:id
Authorization: Bearer <token>
```

#### Get All News Tags
```http
GET /api/v1/news-tags
```

#### Admin: Manage News
```http
POST /api/v1/news
PUT /api/v1/news/:id
DELETE /api/v1/news/:id
Authorization: Bearer <admin_token>
```

### Complaints & Analytics

#### Create Scan Complaint
```http
POST /api/v1/complaints/scan
Authorization: Bearer <token>
Content-Type: application/json

{
  "scan_history_id": "uuid",
  "predicted_disease_id": "uuid",
  "confidence_score": 0.85,
  "category": "wrong_disease",
  "content": "The prediction was incorrect, my plant has rust not blight",
  "image_url": "https://example.com/evidence.jpg",
  "user_suggested_disease_id": "uuid"
}
```

#### Get My Complaints
```http
GET /api/v1/complaints/my?page=1&limit=10
Authorization: Bearer <token>
```

#### Admin: Get All Complaints
```http
GET /api/v1/complaints?page=1&limit=20&status=pending
Authorization: Bearer <admin_token>
```

#### Admin: Get Unverified Scan Complaints
```http
GET /api/v1/complaints/unverified?page=1&limit=20
Authorization: Bearer <admin_token>
```

#### Admin: Verify Complaint
```http
POST /api/v1/complaints/:id/verify
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "is_verified": true,
  "admin_notes": "Verified as correct feedback"
}
```

#### Admin: Analytics - Problematic Diseases
```http
GET /api/v1/analytics/problematic-diseases?limit=10
Authorization: Bearer <admin_token>
```

#### Admin: Analytics - Complaint Trends
```http
GET /api/v1/analytics/trends?period=monthly&start_date=2024-01-01&end_date=2024-12-31
Authorization: Bearer <admin_token>
```

#### Admin: Analytics - Overall Stats
```http
GET /api/v1/analytics/overall-stats
Authorization: Bearer <admin_token>
```

#### Admin: Analytics - Top Contributors
```http
GET /api/v1/analytics/top-contributors?limit=10
Authorization: Bearer <admin_token>
```

#### Admin: Export Training Data
```http
GET /api/v1/ml/export-training-data?format=csv
Authorization: Bearer <admin_token>
```

## Data Models

### User
```go
type User struct {
    ID        string    `json:"id"`         // UUID
    Email     string    `json:"email"`      // Unique
    Username  string    `json:"username"`   // Unique
    Password  string    `json:"-"`          // Hashed, not exposed in JSON
    FullName  string    `json:"full_name"`
    Avatar    string    `json:"avatar"`
    Role      string    `json:"role"`       // "user" or "admin"
    IsActive  bool      `json:"is_active"`  // Account status
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Disease
```go
type Disease struct {
    ID          string    `json:"id"`           // UUID
    Name        string    `json:"name"`
    ClassName   string    `json:"class_name"`   // Unique identifier for ML model
    Type        string    `json:"type"`         // e.g., "fungal", "bacterial", "viral"
    Description string    `json:"description"`
    Solution    []string  `json:"solution"`     // Array of treatment steps
    ImageLink   []string  `json:"image_link"`   // Array of reference images
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Activity
```go
type Activity struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    ActivityDate time.Time `json:"activity_date"`
    Cost         float64   `json:"cost"`
    Revenue      float64   `json:"revenue"`
    Location     string    `json:"location"`
    Keywords     []string  `json:"keywords"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### ScanHistory
```go
type ScanHistory struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`
    DiseaseID       string    `json:"disease_id"`
    ImageURL        string    `json:"image_url"`
    ConfidenceScore float64   `json:"confidence_score"`
    ScanDate        time.Time `json:"scan_date"`
    CreatedAt       time.Time `json:"created_at"`
}
```

### Complaint
```go
type Complaint struct {
    ID                      string    `json:"id"`
    UserID                  string    `json:"user_id"`
    TargetType              string    `json:"target_type"`      // "SCAN", "POST", "COMMENT"
    TargetID                string    `json:"target_id"`
    Category                string    `json:"category"`
    Content                 string    `json:"content"`
    ImageURL                string    `json:"image_url"`
    IsVerified              bool      `json:"is_verified"`
    AdminNotes              string    `json:"admin_notes"`
    // Scan-specific fields
    PredictedDiseaseID      string    `json:"predicted_disease_id"`
    UserSuggestedDiseaseID  string    `json:"user_suggested_disease_id"`
    ConfidenceScore         float64   `json:"confidence_score"`
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
}
```

### Post
```go
type Post struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    Content      string    `json:"content"`
    Images       []string  `json:"images"`
    Visibility   string    `json:"visibility"`  // "public", "private"
    LikeCount    int       `json:"like_count"`
    CommentCount int       `json:"comment_count"`
    ShareCount   int       `json:"share_count"`
    IsDeleted    bool      `json:"is_deleted"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

## Security

### Authentication & Authorization

The API uses JWT (JSON Web Tokens) for authentication. Tokens include:
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "user|admin",
  "exp": 1234567890
}
```

**Middleware:**
- `AuthMiddleware()` - Validates JWT token for protected routes
- `RequireAdmin()` - Ensures user has admin role
- `RequireRole(role)` - Checks for specific role

### Password Security

- Passwords are hashed using bcrypt with a cost factor of 10
- Minimum password requirements should be enforced client-side
- Password reset uses time-limited OTP sent via email

### CORS Configuration

The API supports configurable CORS origins:
- Development: `localhost:3000`, `localhost:5173`, etc.
- Production: Set via `ALLOWED_ORIGINS` environment variable
- Supports Vercel preview deployments automatically

### Best Practices

1. **Change JWT Secret**: Always use a strong, unique JWT secret in production
2. **Use HTTPS**: Deploy behind HTTPS in production
3. **Environment Variables**: Never commit `.env` files
4. **Rate Limiting**: Consider implementing rate limiting for public endpoints
5. **Input Validation**: All inputs are validated before processing
6. **SQL Injection**: GORM provides protection against SQL injection
7. **XSS Protection**: Sanitize user-generated content on the frontend

## Deployment

### Environment Variables for Production

```env
DATABASE_URL=postgresql://user:pass@prod-host:5432/plantheon?sslmode=require
PORT=3000
JWT_SECRET=<strong-random-secret-min-32-chars>
ALLOWED_ORIGINS=https://plantheon.com,https://admin.plantheon.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=noreply@plantheon.com
SMTP_PASSWORD=<app-specific-password>
```

### Build for Production

```bash
# Build binary
go build -o plantheon-backend

# Run
./plantheon-backend
```

### Docker Deployment

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o plantheon-backend

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/plantheon-backend .
EXPOSE 3000
CMD ["./plantheon-backend"]
```

### Database Migrations

For production, consider using a migration tool:
```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path ./migrations -database "postgresql://user:pass@host:5432/db?sslmode=require" up
```

## Contributing

We welcome contributions! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add comments for exported functions
- Write tests for new features

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please:
- Open an issue on GitHub
- Contact the development team
- Check the documentation

---

**Built with ❤️ for the agricultural community**
