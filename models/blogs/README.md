# Blogs/News API Documentation

## Overview
API endpoints for managing news/blog articles in the Plantheon system. News items can be created, updated, deleted, and viewed with user authentication.

## Base URL
`/api/v1/news`

## Authentication
- User routes (GET): Requires `AuthMiddleware()` - User must be authenticated
- Admin routes (POST, PUT, DELETE): Requires `RequireAdmin()` - User must be authenticated and have admin role

## Endpoints

### 1. Get All News
Returns a list of all published news articles.

**Endpoint:** `GET /api/v1/news`

**Auth Required:** Yes (User authenticated)

**Request:** No body required

**Query Parameters:** None

**Response:**
```json
{
  "data": {
    "news": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "How to Grow Tomatoes Successfully",
        "content": "Complete guide on growing tomatoes in your garden...",
        "cover_image_url": "https://example.com/images/tomato.jpg",
        "status": "published",
        "published_at": "2024-01-15T10:30:00Z",
        "created_at": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-15T10:30:00Z",
        "user_id": "660e8400-e29b-41d4-a716-446655440000",
        "full_name": "John Doe",
        "avatar": "https://example.com/avatars/john.jpg"
      }
    ],
    "total": 1
  }
}
```

**Status Codes:**
- `200 OK`: Success
- `500 Internal Server Error`: Server error

---

### 2. Get News by ID
Returns detailed information about a specific news article.

**Endpoint:** `GET /api/v1/news/:id`

**Auth Required:** Yes (User authenticated)

**Path Parameters:**
- `id` (string, required): UUID of the news article

**Request:** No body required

**Response:**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "How to Grow Tomatoes Successfully",
    "content": "Complete guide on growing tomatoes in your garden. This includes soil preparation, planting, watering, and harvesting techniques.",
    "cover_image_url": "https://example.com/images/tomato.jpg",
    "status": "published",
    "published_at": "2024-01-15T10:30:00Z",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "user_id": "660e8400-e29b-41d4-a716-446655440000",
    "full_name": "John Doe",
    "avatar": "https://example.com/avatars/john.jpg"
  }
}
```

**Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: News article not found
- `500 Internal Server Error`: Server error

---

### 3. Create News (Admin Only)
Creates a new news article. Only users with admin role can create news.

**Endpoint:** `POST /api/v1/news`

**Auth Required:** Yes (Admin role required)

**Request Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "New Gardening Techniques for 2024",
  "content": "Discover the latest gardening techniques and innovations...",
  "cover_image_url": "https://example.com/images/techniques.jpg",
  "status": "published"
}
```

**Request Fields:**
- `title` (string, required): News article title (max 200 characters)
- `content` (string, required): News article content (max 10,000 characters)
- `cover_image_url` (string, optional): URL to the cover image (max 500 characters)
- `status` (string, optional): Article status - `"draft"` or `"published"` (default: `"draft"`)

**Response:**
```json
{
  "message": "News created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "New Gardening Techniques for 2024",
    "content": "Discover the latest gardening techniques and innovations...",
    "cover_image_url": "https://example.com/images/techniques.jpg",
    "status": "published",
    "published_at": "2024-01-20T08:00:00Z",
    "user_id": "660e8400-e29b-41d4-a716-446655440000",
    "created_at": "2024-01-20T08:00:00Z"
  }
}
```

**Status Codes:**
- `201 Created`: News created successfully
- `400 Bad Request`: Invalid request data or validation error
- `401 Unauthorized`: User not authenticated or not admin
- `500 Internal Server Error`: Server error

**Validation Rules:**
- Title is required and must be less than 200 characters
- Content is required and must be less than 10,000 characters
- Cover image URL must be less than 500 characters
- Status must be either "draft" or "published"

---

### 4. Update News (Admin Only)
Updates an existing news article. Only users with admin role can update news.

**Endpoint:** `PUT /api/v1/news/:id`

**Auth Required:** Yes (Admin role required)

**Path Parameters:**
- `id` (string, required): UUID of the news article to update

**Request Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Updated Gardening Techniques for 2024",
  "content": "Updated content with new information...",
  "cover_image_url": "https://example.com/images/new-techniques.jpg",
  "status": "published"
}
```

**Request Fields:**
- `title` (string, required): News article title (max 200 characters)
- `content` (string, required): News article content (max 10,000 characters)
- `cover_image_url` (string, optional): URL to the cover image (max 500 characters)
- `status` (string, optional): Article status - `"draft"` or `"published"`

**Response:**
```json
{
  "message": "News updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Updated Gardening Techniques for 2024",
    "content": "Updated content with new information...",
    "cover_image_url": "https://example.com/images/new-techniques.jpg",
    "status": "published",
    "updated_at": "2024-01-25T14:30:00Z"
  }
}
```

**Status Codes:**
- `200 OK`: News updated successfully
- `400 Bad Request`: Invalid UUID format or validation error
- `401 Unauthorized`: User not authenticated or not admin
- `404 Not Found`: News article not found
- `500 Internal Server Error`: Server error

**Special Behavior:**
- If status is changed from "draft" to "published" and `published_at` is `null`, it will be automatically set to the current timestamp

---

### 5. Delete News (Admin Only)
Deletes a news article. Only users with admin role can delete news.

**Endpoint:** `DELETE /api/v1/news/:id`

**Auth Required:** Yes (Admin role required)

**Path Parameters:**
- `id` (string, required): UUID of the news article to delete

**Request:** No body required

**Response:**
```json
{
  "message": "News deleted successfully"
}
```

**Status Codes:**
- `200 OK`: News deleted successfully
- `400 Bad Request`: Invalid UUID format
- `401 Unauthorized`: User not authenticated or not admin
- `500 Internal Server Error`: Server error

---

## Error Response Format

All error responses follow this format:

```json
{
  "error": "Error message describing what went wrong"
}
```

## Common Error Messages

- `"User not found in context"`: Authentication token is missing or invalid
- `"Invalid user format"`: Internal server error with user data
- `"Invalid request"`: Request body is malformed or missing required fields
- `"News title is required"`: Title field is empty
- `"News content is required"`: Content field is empty
- `"News title must be less than 200 characters"`: Title exceeds character limit
- `"News content must be less than 10,000 characters"`: Content exceeds character limit
- `"Status must be either 'draft' or 'published'"`: Invalid status value
- `"Cover image URL must be less than 500 characters"`: Cover image URL exceeds character limit
- `"id parameter is required"`: ID parameter is missing
- `"id parameter must be a valid UUID"`: ID is not in valid UUID format
- `"Failed to create news"`: Server error during creation
- `"Failed to get news"`: Server error during retrieval
- `"Failed to update news"`: Server error during update
- `"Failed to delete news"`: Server error during deletion
- `"News article not found"`: News with specified ID does not exist

## Data Model

### Blog Model
```go
type Blog struct {
    ID               string     `json:"id"`                 // UUID
    Title            string     `json:"title"`              // Required, max 200 chars
    Content          string     `json:"content"`            // Required, max 10,000 chars
    CoverImageURL    *string    `json:"cover_image_url"`    // Optional, max 500 chars
    SubGuideStagesID *string    `json:"sub_guide_stages_id"` // Optional, UUID
    UserID           string     `json:"user_id"`            // Auto-filled from token
    Status           string     `json:"status"`              // "draft" or "published"
    PublishedAt      *time.Time `json:"published_at"`       // Auto-set when published
    CreatedAt        time.Time  `json:"created_at"`          // Auto-generated
    UpdatedAt        time.Time  `json:"updated_at"`          // Auto-generated
}
```

## Usage Examples

### Example 1: Create a Draft News Article
```bash
curl -X POST http://localhost:3000/api/v1/news \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Hydroponics",
    "content": "Learn the basics of hydroponic gardening...",
    "cover_image_url": "https://example.com/hydroponics.jpg",
    "status": "draft"
  }'
```

### Example 2: Publish a News Article
```bash
curl -X PUT http://localhost:3000/api/v1/news/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Hydroponics",
    "content": "Learn the basics of hydroponic gardening...",
    "cover_image_url": "https://example.com/hydroponics.jpg",
    "status": "published"
  }'
```

### Example 3: Get All Published News
```bash
curl -X GET http://localhost:3000/api/v1/news \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Example 4: Get Specific News Article
```bash
curl -X GET http://localhost:3000/api/v1/news/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Example 5: Delete News Article
```bash
curl -X DELETE http://localhost:3000/api/v1/news/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

