# Plantheon Backend API

> ### Backend API cho ứng dụng Plantheon - Hệ thống quản lý thông tin bệnh cây trồng

API được xây dựng với **Golang/Gin** bao gồm authentication, CRUD operations, pagination và nhiều tính năng khác.

## Cấu trúc thư mục

```
.
├── main.go                 // Entry point của ứng dụng
├── go.mod                  // Go modules
├── migration.sql           // Database migration script
├── common/
│   ├── database.go         // Database connection manager
│   └── utils.go            // Utility functions (JWT, password hashing)
├── users/
│   ├── models.go           // User model & database operations
│   ├── serializers.go      // Request/Response structures
│   ├── routers.go          // API handlers
│   ├── middlewares.go      // Authentication & authorization middleware
│   └── validators.go       // Input validation
└── diseases/
    ├── models.go           // Disease model & database operations
    ├── serializers.go      // Request/Response structures
    ├── routers.go          // API handlers
    └── validators.go       // Input validation
```

## Hệ thống phân quyền

API sử dụng hệ thống role-based authentication với 2 roles:

### User Roles

- **`user`** (default): Người dùng thông thường
  - Có thể xem danh sách và chi tiết bệnh
  - Có thể cập nhật profile của mình
- **`admin`**: Quản trị viên
  - Có tất cả quyền của user
  - Có thể tạo, sửa, xóa bệnh
  - Có thể quản lý người dùng (tương lai)

### JWT Token

JWT token bây giờ bao gồm thông tin role:

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "user|admin",
  "exp": "expiration_time"
}
```

### Middleware Authorization

- `AuthMiddleware()`: Xác thực token (cho user routes)
- `RequireAdmin()`: Yêu cầu admin role
- `RequireRole(role)`: Yêu cầu role cụ thể

## Cài đặt và chạy

### 1. Cài đặt dependencies

```bash
go mod tidy
```

### 2. Cấu hình database

Tạo file `.env` trong thư mục gốc project với nội dung:

```env
# PostgreSQL URL format (Supabase, Railway, etc.)
DATABASE_URL=postgresql://postgres:your_password@host:5432/database_name?sslmode=require

# Hoặc traditional DSN format
# DATABASE_URL=host=localhost user=postgres password=your_password dbname=plantheon port=5432 sslmode=disable

# Server port
PORT=8080

# JWT Secret
JWT_SECRET=your-super-secret-jwt-key-here
```

**Ví dụ với Supabase:**

```env
DATABASE_URL=postgresql://postgres:Quyen@213@db.qlxcxrrhrrlfaqplqkmz.supabase.co:5432/postgres
PORT=8080
JWT_SECRET=plantheon-jwt-secret-2024-very-secure-key
```

### 3. Database Migration (nếu có database cũ)

Nếu bạn đã có database từ version trước, chạy migration script:

```bash
# Chạy migration SQL để thêm role column
psql -d your_database -f migration.sql
```

### 4. Chạy ứng dụng

```bash
go run main.go
```

Server sẽ chạy trên port 8080 (có thể thay đổi bằng environment variable PORT).

## API Endpoints

### Authentication

#### Đăng ký

```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "full_name": "Tên đầy đủ"
}
```

#### Đăng nhập

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### User Profile (Cần Authentication)

#### Xem profile

```http
GET /api/users/profile
Authorization: Bearer <jwt_token>
```

#### Cập nhật profile

```http
PUT /api/users/profile
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "username": "new_username",
  "full_name": "Tên mới",
  "avatar": "http://example.com/avatar.jpg"
}
```

### Diseases

#### Lấy danh sách bệnh (có pagination, search, filter)

```http
GET /api/diseases?page=1&limit=10&type=fungal&search=keyword
```

#### Lấy thông tin bệnh theo ID

```http
GET /api/diseases/:id
```

#### Lấy thông tin bệnh theo Class Name

```http
GET /api/diseases/class/:className
```

#### Tạo bệnh mới (Cần Authentication)

```http
POST /api/diseases
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "Tên bệnh",
  "class_name": "class_name",
  "type": "fungal",
  "description": "Mô tả bệnh",
  "solution": ["Giải pháp 1", "Giải pháp 2"],
  "image_link": ["http://example.com/image1.jpg", "http://example.com/image2.jpg"]
}
```

#### Cập nhật bệnh (Cần Authentication)

```http
PUT /api/diseases/:id
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "Tên bệnh mới",
  "description": "Mô tả mới"
}
```

#### Xóa bệnh (Cần Authentication)

```http
DELETE /api/diseases/:id
Authorization: Bearer <jwt_token>
```

## Models

### User Model

```go
type User struct {
    ID        string    `json:"id"`        // UUID
    Email     string    `json:"email"`     // Unique
    Username  string    `json:"username"`  // Unique
    Password  string    `json:"-"`         // Hashed
    FullName  string    `json:"full_name"`
    Avatar    string    `json:"avatar"`
    Role      UserRole  `json:"role"`      // "user" or "admin"
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Disease Model

```go
type Disease struct {
    ID          string         `json:"id"`          // UUID
    Name        string         `json:"name"`        // Tên bệnh
    ClassName  string         `json:"class_name"` // Tên class
    Type        string         `json:"type"`        // Loại bệnh
    Description string         `json:"description"` // Mô tả
    Solution    []string       `json:"solution"`    // Array giải pháp
    ImageLink   []string       `json:"image_link"`  // Array link ảnh
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
}
```

## Features

- ✅ JWT Authentication với role-based authorization
- ✅ Password hashing với bcrypt
- ✅ Role-based access control (User/Admin)
- ✅ CRUD operations cho User và Disease
- ✅ Admin-only disease management
- ✅ Pagination cho danh sách
- ✅ Search và filter
- ✅ Input validation
- ✅ CORS support
- ✅ PostgreSQL với GORM
- ✅ Array fields support (solution, image_link)

## Health Check

```http
GET /api/health
```

Response:

```json
{
  "status": "OK",
  "message": "Plantheon Backend API is running"
}
```

## Lưu ý

1. Thay đổi JWT secret trong production (file `common/utils.go`)
2. Cấu hình database connection string phù hợp
3. Sử dụng HTTPS trong production
4. Thêm rate limiting nếu cần
5. Cấu hình logging phù hợp
