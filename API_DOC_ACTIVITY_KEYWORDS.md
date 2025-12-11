# Activity Keywords API Documentation

**Base URL:** `http://localhost:8080/api/v1`

Quản lý Activity Keywords và liên kết với Disease trong một hệ thống thống nhất.

---

## Tổng quan

- **Activity Keyword**: Từ khóa hoạt động (VD: "Tưới nước", "Bón phân", "Phun thuốc")
- **Disease-Activity Keyword**: Bảng liên kết keyword với bệnh (nhiều-nhiều)

**Ý tưởng "Làm 1 được 2"**: Khi lấy danh sách keywords, response sẽ kèm theo các bệnh liên kết. Khi CRUD keyword, có thể cập nhật liên kết bệnh luôn.

---

## 1. ACTIVITY KEYWORDS CRUD

### 1.1 Get All Activity Keywords (Public)

**Endpoint:** `GET /activity-keywords`

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid-1",
      "name": "Phun thuốc trừ sâu",
      "description": "Phun thuốc định kỳ...",
      "type": "TREATMENT",
      "base_days_offset": 7,
      "is_free_time": false,
      "hour_time": 8,
      "end_hour_time": 10,
      "time_duration": 30,
      "created_at": "2025-12-10T00:00:00Z",
      "updated_at": "2025-12-10T00:00:00Z"
    }
  ]
}
```

---

### 1.2 Get Activity Keyword by ID (Public)

**Endpoint:** `GET /activity-keywords/:id`

**Response (200):**
```json
{
  "data": {
    "id": "uuid-1",
    "name": "Phun thuốc trừ sâu",
    "description": "...",
    "type": "TREATMENT",
    "base_days_offset": 7,
    "is_free_time": false,
    "hour_time": 8,
    "end_hour_time": 10,
    "time_duration": 30,
    "created_at": "...",
    "updated_at": "..."
  }
}
```

**Response Error (404):**
```json
{
  "error": "Activity keyword not found"
}
```

---

### 1.3 Search Activity Keywords (Public)

**Endpoint:** `GET /activity-keywords/search?q=<keyword>`

**Query Parameters:**
| Parameter | Type   | Required | Description |
|-----------|--------|----------|-------------|
| q         | string | Yes      | Từ khóa tìm kiếm |

**Example:** `GET /activity-keywords/search?q=phun`

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid-1",
      "name": "Phun thuốc trừ sâu",
      "type": "TREATMENT",
      ...
    }
  ]
}
```

---

### 1.4 Get Activity Keywords for Disease (Public)

Lấy tất cả keywords liên kết với một bệnh cụ thể.

**Endpoint:** `GET /activity-keywords/disease/:disease_id`

**Example:** `GET /activity-keywords/disease/abc-123`

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid-1",
      "name": "Phun thuốc trừ sâu",
      "description": "...",
      "type": "TREATMENT",
      "base_days_offset": 7,
      "is_free_time": false,
      "hour_time": 8,
      "end_hour_time": 10,
      "time_duration": 30,
      "created_at": "...",
      "updated_at": "..."
    },
    {
      "id": "uuid-2",
      "name": "Cắt tỉa lá bệnh",
      "type": "PREVENTION",
      ...
    }
  ]
}
```

---

### 1.5 Create Activity Keyword (Admin)

**Endpoint:** `POST /activity-keywords`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_token>
```

**Request Body:**
```json
{
  "name": "Phun thuốc trừ sâu",
  "description": "Phun thuốc định kỳ để phòng sâu bệnh",
  "type": "TREATMENT",
  "base_days_offset": 7,
  "is_free_time": false,
  "hour_time": 8,
  "end_hour_time": 10,
  "time_duration": 30
}
```

**Request Fields:**

| Field           | Type    | Required | Description                           |
|-----------------|---------|----------|---------------------------------------|
| name            | string  | Yes      | Tên keyword                           |
| description     | string  | No       | Mô tả                                 |
| type            | string  | Yes      | Loại: GENERAL, TREATMENT, PREVENTION  |
| base_days_offset| int     | No       | Số ngày offset mặc định (default: 0)  |
| is_free_time    | bool    | No       | Thời gian tự do (default: false)      |
| hour_time       | int     | No       | Giờ bắt đầu (0-23)                    |
| end_hour_time   | int     | No       | Giờ kết thúc (0-23)                   |
| time_duration   | int     | No       | Thời lượng (phút)                     |

**Response Success (201):**
```json
{
  "data": {
    "id": "uuid-new",
    "name": "Phun thuốc trừ sâu",
    ...
  }
}
```

---

### 1.6 Update Activity Keyword (Admin)

**Endpoint:** `PUT /activity-keywords/:id`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_token>
```

**Request Body:** (Chỉ gửi field cần update)
```json
{
  "name": "Tên mới",
  "description": "Mô tả mới",
  "type": "PREVENTION"
}
```

**Response Success (200):**
```json
{
  "data": {
    "id": "uuid-1",
    "name": "Tên mới",
    ...
  }
}
```

---

### 1.7 Delete Activity Keyword (Admin)

**Endpoint:** `DELETE /activity-keywords/:id`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response Success (200):**
```json
{
  "message": "Activity keyword deleted successfully"
}
```

---

## 2. DISEASE-KEYWORD RELATIONSHIPS

### 2.1 Add Keyword to Disease (Admin)

Liên kết một keyword với một bệnh.

**Endpoint:** `POST /disease-activity-keywords/disease/:disease_id`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_token>
```

**Request Body:**
```json
{
  "activity_keyword_id": "uuid-keyword"
}
```

**Response Success (200):**
```json
{
  "message": "Activity keyword added successfully"
}
```

---

### 2.2 Remove Keyword from Disease (Admin)

Xóa liên kết keyword khỏi bệnh.

**Endpoint:** `DELETE /disease-activity-keywords/disease/:disease_id/keyword/:keyword_id`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response Success (200):**
```json
{
  "message": "Activity keyword removed successfully"
}
```

---

### 2.3 Set All Keywords for Disease (Admin)

Thay thế toàn bộ keywords của một bệnh (xóa cũ, thêm mới).

**Endpoint:** `PUT /disease-activity-keywords/disease/:disease_id`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_token>
```

**Request Body:**
```json
{
  "activity_keyword_ids": ["uuid-1", "uuid-2", "uuid-3"]
}
```

**Response Success (200):**
```json
{
  "message": "Activity keywords set successfully"
}
```

---

### 2.4 Import Keywords + Relationships from CSV (Admin)

Import keywords và tự động liên kết với bệnh theo ClassName.

**Endpoint:** `POST /disease-activity-keywords/import-csv`

**Headers:**
```
Content-Type: multipart/form-data
Authorization: Bearer <admin_token>
```

**Form Field:** `file` (CSV file)

**CSV Format (8 columns):**

| Column | Name               | Required | Description                    |
|--------|--------------------|----------|--------------------------------|
| 1      | ClassName          | Yes      | Class name của bệnh            |
| 2      | keywordName        | Yes      | Tên keyword                    |
| 3      | keywordDescription | No       | Mô tả keyword                  |
| 4      | keywordType        | Yes      | Loại: GENERAL, TREATMENT, etc. |
| 5      | keywordDayOffset   | No       | Số ngày offset                 |
| 6      | keywordIsFreeTime  | No       | true/false hoặc 1/0            |
| 7      | keywordHourTime    | No       | Giờ bắt đầu (0-23)             |
| 8      | timeDuration       | No       | Thời lượng (phút)              |

**Example CSV:**
```csv
ClassName,keywordName,keywordDescription,keywordType,keywordDayOffset,keywordIsFreeTime,keywordHourTime,timeDuration
leaf_spot,Phun thuốc,Phun thuốc định kỳ,TREATMENT,7,false,8,30
leaf_spot,Cắt tỉa lá bệnh,Loại bỏ lá nhiễm bệnh,PREVENTION,0,true,,
root_rot,Cải thiện thoát nước,Xử lý đất,PREVENTION,0,false,9,60
```

**Cách hoạt động:**
1. Tìm disease theo ClassName
2. Nếu keyword chưa tồn tại → Tạo mới
3. Nếu keyword đã tồn tại → Dùng lại
4. Tạo liên kết disease-keyword

**Response Success (200):**
```json
{
  "message": "Imported 8 out of 10 rows with 2 errors",
  "success": 8,
  "total": 10,
  "errors": [
    {
      "row": 5,
      "error": "Disease not found",
      "details": "No disease with ClassName 'unknown_disease'"
    }
  ]
}
```

---

## 3. JavaScript Integration

```javascript
const API_URL = 'http://localhost:8080/api/v1';
const getToken = () => localStorage.getItem('token');

// ============ ACTIVITY KEYWORDS ============

// Get all keywords
const getAllKeywords = async () => {
  const res = await fetch(`${API_URL}/activity-keywords`);
  return (await res.json()).data;
};

// Get keywords by disease
const getKeywordsByDisease = async (diseaseId) => {
  const res = await fetch(`${API_URL}/activity-keywords/disease/${diseaseId}`);
  return (await res.json()).data;
};

// Search keywords
const searchKeywords = async (query) => {
  const res = await fetch(`${API_URL}/activity-keywords/search?q=${encodeURIComponent(query)}`);
  return (await res.json()).data;
};

// Create keyword (Admin)
const createKeyword = async (data) => {
  const res = await fetch(`${API_URL}/activity-keywords`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify(data)
  });
  return (await res.json()).data;
};

// Update keyword (Admin)
const updateKeyword = async (id, data) => {
  const res = await fetch(`${API_URL}/activity-keywords/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify(data)
  });
  return (await res.json()).data;
};

// Delete keyword (Admin)
const deleteKeyword = async (id) => {
  const res = await fetch(`${API_URL}/activity-keywords/${id}`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return await res.json();
};

// ============ DISEASE-KEYWORD RELATIONSHIPS ============

// Add keyword to disease (Admin)
const addKeywordToDisease = async (diseaseId, keywordId) => {
  const res = await fetch(`${API_URL}/disease-activity-keywords/disease/${diseaseId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({ activity_keyword_id: keywordId })
  });
  return await res.json();
};

// Remove keyword from disease (Admin)
const removeKeywordFromDisease = async (diseaseId, keywordId) => {
  const res = await fetch(`${API_URL}/disease-activity-keywords/disease/${diseaseId}/keyword/${keywordId}`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return await res.json();
};

// Set all keywords for disease (Admin)
const setKeywordsForDisease = async (diseaseId, keywordIds) => {
  const res = await fetch(`${API_URL}/disease-activity-keywords/disease/${diseaseId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({ activity_keyword_ids: keywordIds })
  });
  return await res.json();
};

// Import CSV (Admin)
const importKeywordsCSV = async (file) => {
  const formData = new FormData();
  formData.append('file', file);
  
  const res = await fetch(`${API_URL}/disease-activity-keywords/import-csv`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${getToken()}` },
    body: formData
  });
  return await res.json();
};
```

---

## 4. Tóm tắt Endpoints

### Public (Không cần token):
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/activity-keywords` | Lấy tất cả keywords |
| GET | `/activity-keywords/:id` | Lấy keyword theo ID |
| GET | `/activity-keywords/search?q=` | Tìm kiếm keywords |
| GET | `/activity-keywords/disease/:disease_id` | Lấy keywords của bệnh |

### Admin (Cần Bearer token):
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/activity-keywords` | Tạo keyword |
| PUT | `/activity-keywords/:id` | Cập nhật keyword |
| DELETE | `/activity-keywords/:id` | Xóa keyword |
| POST | `/disease-activity-keywords/disease/:disease_id` | Thêm keyword vào bệnh |
| DELETE | `/disease-activity-keywords/disease/:disease_id/keyword/:keyword_id` | Xóa liên kết |
| PUT | `/disease-activity-keywords/disease/:disease_id` | Set keywords cho bệnh |
| POST | `/disease-activity-keywords/import-csv` | Import CSV |
