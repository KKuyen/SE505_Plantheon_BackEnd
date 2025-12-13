# Complaints API Documentation

**Base URL:** `http://localhost:8080/api/v1`

Quản lý khiếu nại cho bài viết (Post) và bình luận (Comment).

---

## Tổng quan

- **Complaint**: Khiếu nại từ người dùng về nội dung vi phạm
- **Target Type**: Loại nội dung bị khiếu nại (`POST` hoặc `COMMENT`)
- **Category**: Lý do khiếu nại (SPAM, HARASSMENT, etc.)
- **Status**: Trạng thái xử lý (PENDING → REVIEWED → RESOLVED/REJECTED)

---

## Các loại khiếu nại (Category)

| Category | Mô tả |
|----------|-------|
| `SPAM` | Spam, quảng cáo |
| `HARASSMENT` | Quấy rối, bắt nạt |
| `HATE_SPEECH` | Phát ngôn thù địch |
| `VIOLENCE` | Nội dung bạo lực |
| `MISINFORMATION` | Thông tin sai lệch |
| `INAPPROPRIATE` | Nội dung không phù hợp |
| `OTHER` | Lý do khác |

## Trạng thái khiếu nại (Status)

| Status | Mô tả |
|--------|-------|
| `PENDING` | Chờ xử lý (mặc định) |
| `REVIEWED` | Đang xem xét |
| `RESOLVED` | Đã giải quyết |
| `REJECTED` | Từ chối (không vi phạm) |

---

## 1. USER ENDPOINTS (Cần Bearer Token)

### 1.1 Tạo khiếu nại mới

**Endpoint:** `POST /complaints`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <user_token>
```

**Request Body:**
```json
{
  "target_id": "uuid-of-post-or-comment",
  "target_type": "POST",
  "category": "SPAM",
  "content": "Bài viết này chứa nội dung quảng cáo spam"
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| target_id | string (UUID) | Yes | ID của bài viết hoặc bình luận |
| target_type | string | Yes | `POST` hoặc `COMMENT` |
| category | string | Yes | Loại khiếu nại (xem bảng trên) |
| content | string | No | Nội dung chi tiết (tối đa 1000 ký tự) |

**Response Success (201):**
```json
{
  "message": "Complaint submitted successfully",
  "data": {
    "id": "uuid-complaint",
    "user_id": "uuid-user",
    "target_id": "uuid-post",
    "target_type": "POST",
    "category": "SPAM",
    "content": "Bài viết này chứa nội dung quảng cáo spam",
    "status": "PENDING",
    "created_at": "2025-12-12T10:00:00Z",
    "updated_at": "2025-12-12T10:00:00Z"
  }
}
```

**Response Error (409 - Duplicate):**
```json
{
  "error": "You have already submitted a complaint for this content"
}
```

---

### 1.2 Lấy danh sách khiếu nại của tôi

**Endpoint:** `GET /complaints/my`

**Headers:**
```
Authorization: Bearer <user_token>
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Trang hiện tại |
| limit | int | 10 | Số lượng mỗi trang (max 100) |

**Example:** `GET /complaints/my?page=1&limit=10`

**Response (200):**
```json
{
  "data": {
    "complaints": [
      {
        "id": "uuid-1",
        "user_id": "uuid-user",
        "target_id": "uuid-post-1",
        "target_type": "POST",
        "category": "SPAM",
        "content": "Nội dung spam",
        "status": "PENDING",
        "created_at": "2025-12-12T10:00:00Z",
        "updated_at": "2025-12-12T10:00:00Z"
      },
      {
        "id": "uuid-2",
        "user_id": "uuid-user",
        "target_id": "uuid-comment-1",
        "target_type": "COMMENT",
        "category": "HARASSMENT",
        "content": "Bình luận quấy rối",
        "status": "RESOLVED",
        "admin_notes": "Đã xóa bình luận vi phạm",
        "resolved_at": "2025-12-12T15:00:00Z",
        "resolved_by": "uuid-admin",
        "created_at": "2025-12-11T08:00:00Z",
        "updated_at": "2025-12-12T15:00:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "limit": 10
  }
}
```

---

### 1.3 Xem chi tiết khiếu nại

**Endpoint:** `GET /complaints/:id`

**Headers:**
```
Authorization: Bearer <user_token>
```

**Response (200):**
```json
{
  "data": {
    "id": "uuid-complaint",
    "user_id": "uuid-user",
    "target_id": "uuid-post",
    "target_type": "POST",
    "category": "SPAM",
    "content": "Nội dung chi tiết...",
    "status": "PENDING",
    "created_at": "2025-12-12T10:00:00Z",
    "updated_at": "2025-12-12T10:00:00Z"
  }
}
```

**Response Error (403):**
```json
{
  "error": "You can only view your own complaints"
}
```

---

### 1.4 Xóa khiếu nại (chỉ PENDING)

**Endpoint:** `DELETE /complaints/:id`

**Headers:**
```
Authorization: Bearer <user_token>
```

**Response Success (200):**
```json
{
  "message": "Complaint deleted successfully"
}
```

**Response Error (400):**
```json
{
  "error": "You can only delete pending complaints"
}
```

---

## 2. ADMIN ENDPOINTS (Cần Admin Token)

### 2.1 Lấy tất cả khiếu nại

**Endpoint:** `GET /complaints`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | int | 1 | Trang hiện tại |
| limit | int | 10 | Số lượng mỗi trang (max 100) |
| status | string | | Filter theo status: `PENDING`, `REVIEWED`, `RESOLVED`, `REJECTED` |
| target_type | string | | Filter theo loại: `POST` hoặc `COMMENT` |

**Examples:**
- `GET /complaints?page=1&limit=20`
- `GET /complaints?status=PENDING`
- `GET /complaints?target_type=POST&status=PENDING`

**Response (200):**
```json
{
  "data": {
    "complaints": [
      {
        "id": "uuid-1",
        "user_id": "uuid-user-1",
        "target_id": "uuid-post-1",
        "target_type": "POST",
        "category": "SPAM",
        "content": "...",
        "status": "PENDING",
        "created_at": "2025-12-12T10:00:00Z",
        "updated_at": "2025-12-12T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "limit": 10
  }
}
```

---

### 2.2 Đếm số khiếu nại

**Endpoint:** `GET /complaints/count`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| status | string | Filter theo status |
| target_type | string | Filter theo loại |

**Examples:**
- `GET /complaints/count` - Tổng số khiếu nại
- `GET /complaints/count?status=PENDING` - Số khiếu nại chờ xử lý
- `GET /complaints/count?target_type=POST` - Số khiếu nại bài viết

**Response (200):**
```json
{
  "data": {
    "count": 25,
    "status": "PENDING",
    "target_type": ""
  }
}
```

---

### 2.3 Lấy khiếu nại theo bài viết/bình luận

**Endpoint:** `GET /complaints/target/:target_id`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| target_type | string | Yes | `POST` hoặc `COMMENT` |

**Example:** `GET /complaints/target/uuid-post-123?target_type=POST`

**Response (200):**
```json
{
  "data": [
    {
      "id": "uuid-complaint-1",
      "user_id": "uuid-user-1",
      "target_id": "uuid-post-123",
      "target_type": "POST",
      "category": "SPAM",
      "content": "Spam quảng cáo",
      "status": "PENDING",
      "created_at": "2025-12-12T10:00:00Z",
      "updated_at": "2025-12-12T10:00:00Z"
    },
    {
      "id": "uuid-complaint-2",
      "user_id": "uuid-user-2",
      "target_id": "uuid-post-123",
      "target_type": "POST",
      "category": "MISINFORMATION",
      "content": "Thông tin sai sự thật",
      "status": "PENDING",
      "created_at": "2025-12-12T11:00:00Z",
      "updated_at": "2025-12-12T11:00:00Z"
    }
  ]
}
```

---

### 2.4 Cập nhật trạng thái khiếu nại

**Endpoint:** `PUT /complaints/:id/status`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_token>
```

**Request Body:**
```json
{
  "status": "RESOLVED",
  "admin_notes": "Đã xóa bài viết vi phạm. Cảm ơn bạn đã báo cáo."
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| status | string | Yes | `PENDING`, `REVIEWED`, `RESOLVED`, `REJECTED` |
| admin_notes | string | No | Ghi chú từ admin (tối đa 1000 ký tự) |

**Response Success (200):**
```json
{
  "message": "Complaint status updated successfully",
  "data": {
    "id": "uuid-complaint",
    "user_id": "uuid-user",
    "target_id": "uuid-post",
    "target_type": "POST",
    "category": "SPAM",
    "content": "...",
    "status": "RESOLVED",
    "admin_notes": "Đã xóa bài viết vi phạm. Cảm ơn bạn đã báo cáo.",
    "resolved_at": "2025-12-12T15:00:00Z",
    "resolved_by": "uuid-admin",
    "created_at": "2025-12-12T10:00:00Z",
    "updated_at": "2025-12-12T15:00:00Z"
  }
}
```

---

### 2.5 Xóa khiếu nại (Admin)

**Endpoint:** `DELETE /complaints/admin/:id`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response Success (200):**
```json
{
  "message": "Complaint deleted successfully"
}
```

---

### 2.6 Xem bài viết (Admin - bao gồm cả bài đã xóa)

**Endpoint:** `GET /admin/posts/:id`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response (200):**
```json
{
  "data": {
    "id": "uuid-post",
    "user_id": "uuid-user",
    "content": "Nội dung bài viết...",
    "image_link": ["url1", "url2"],
    "disease_link": "uuid-disease",
    "tags": ["tag1", "tag2"],
    "like_num": 10,
    "comment_num": 5,
    "share_num": 2,
    "status": "AVAILABLE",
    "is_deleted": false,
    "created_at": "2025-12-12T10:00:00Z",
    "updated_at": "2025-12-12T10:00:00Z"
  }
}
```

---

### 2.7 Cập nhật trạng thái xóa bài viết (Admin)

**Endpoint:** `PUT /admin/posts/:id/is-deleted`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_token>
```

**Request Body:**
```json
{
  "is_deleted": true
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| is_deleted | boolean | Yes | `true` để gỡ bài, `false` để khôi phục |

**Response Success (200):**
```json
{
  "message": "Post deleted successfully",
  "data": {
    "id": "uuid-post",
    "is_deleted": true
  }
}
```

**Response (khôi phục bài viết):**
```json
{
  "message": "Post restored successfully",
  "data": {
    "id": "uuid-post",
    "is_deleted": false
  }
}
```

---

### 2.8 Xem bình luận (Admin - bao gồm cả bình luận đã xóa)

**Endpoint:** `GET /admin/comments/:commentId`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response (200):**
```json
{
  "data": {
    "id": "uuid-comment",
    "post_id": "uuid-post",
    "user_id": "uuid-user",
    "parent_id": "uuid-parent-comment",
    "content": "Nội dung bình luận...",
    "like_num": 5,
    "status": "AVAILABLE",
    "is_deleted": false,
    "created_at": "2025-12-12T10:00:00Z",
    "updated_at": "2025-12-12T10:00:00Z"
  }
}
```

---

### 2.9 Cập nhật trạng thái xóa bình luận (Admin)

**Endpoint:** `PUT /admin/comments/:commentId/is-deleted`

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin_token>
```

**Request Body:**
```json
{
  "is_deleted": true
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| is_deleted | boolean | Yes | `true` để gỡ bình luận, `false` để khôi phục |

**Response Success (200):**
```json
{
  "message": "Comment deleted successfully",
  "data": {
    "id": "uuid-comment",
    "post_id": "uuid-post",
    "is_deleted": true
  }
}
```

**Response (khôi phục bình luận):**
```json
{
  "message": "Comment restored successfully",
  "data": {
    "id": "uuid-comment",
    "post_id": "uuid-post",
    "is_deleted": false
  }
}
```

---

## 3. JavaScript Integration

```javascript
const API_URL = 'http://localhost:8080/api/v1';
const getToken = () => localStorage.getItem('token');

// ============ USER FUNCTIONS ============

// Tạo khiếu nại mới
const createComplaint = async (targetId, targetType, category, content = '') => {
  const res = await fetch(`${API_URL}/complaints`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({
      target_id: targetId,
      target_type: targetType,  // 'POST' or 'COMMENT'
      category: category,        // 'SPAM', 'HARASSMENT', etc.
      content: content
    })
  });
  return await res.json();
};

// Lấy khiếu nại của tôi
const getMyComplaints = async (page = 1, limit = 10) => {
  const res = await fetch(
    `${API_URL}/complaints/my?page=${page}&limit=${limit}`,
    { headers: { 'Authorization': `Bearer ${getToken()}` } }
  );
  return (await res.json()).data;
};

// Xem chi tiết khiếu nại
const getComplaintById = async (id) => {
  const res = await fetch(`${API_URL}/complaints/${id}`, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return (await res.json()).data;
};

// Xóa khiếu nại (chỉ pending)
const deleteMyComplaint = async (id) => {
  const res = await fetch(`${API_URL}/complaints/${id}`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return await res.json();
};

// ============ ADMIN FUNCTIONS ============

// Lấy tất cả khiếu nại (có filter)
const getAllComplaints = async (page = 1, limit = 10, status = '', targetType = '') => {
  let url = `${API_URL}/complaints?page=${page}&limit=${limit}`;
  if (status) url += `&status=${status}`;
  if (targetType) url += `&target_type=${targetType}`;
  
  const res = await fetch(url, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return (await res.json()).data;
};

// Đếm số khiếu nại
const getComplaintsCount = async (status = '', targetType = '') => {
  let url = `${API_URL}/complaints/count`;
  const params = [];
  if (status) params.push(`status=${status}`);
  if (targetType) params.push(`target_type=${targetType}`);
  if (params.length > 0) url += '?' + params.join('&');
  
  const res = await fetch(url, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return (await res.json()).data;
};

// Lấy khiếu nại theo target (bài viết/bình luận)
const getComplaintsByTarget = async (targetId, targetType) => {
  const res = await fetch(
    `${API_URL}/complaints/target/${targetId}?target_type=${targetType}`,
    { headers: { 'Authorization': `Bearer ${getToken()}` } }
  );
  return (await res.json()).data;
};

// Cập nhật trạng thái khiếu nại
const updateComplaintStatus = async (id, status, adminNotes = '') => {
  const res = await fetch(`${API_URL}/complaints/${id}/status`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({
      status: status,  // 'PENDING', 'REVIEWED', 'RESOLVED', 'REJECTED'
      admin_notes: adminNotes
    })
  });
  return await res.json();
};

// Xóa khiếu nại (admin)
const adminDeleteComplaint = async (id) => {
  const res = await fetch(`${API_URL}/complaints/admin/${id}`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return await res.json();
};

// ============ ADMIN POST/COMMENT MODERATION ============

// Xem bài viết (admin - bao gồm cả bài đã xóa)
const adminGetPost = async (postId) => {
  const res = await fetch(`${API_URL}/admin/posts/${postId}`, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return (await res.json()).data;
};

// Gỡ/Khôi phục bài viết (admin)
const adminUpdatePostIsDeleted = async (postId, isDeleted) => {
  const res = await fetch(`${API_URL}/admin/posts/${postId}/is-deleted`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({ is_deleted: isDeleted })
  });
  return await res.json();
};

// Xem bình luận (admin - bao gồm cả bình luận đã xóa)
const adminGetComment = async (commentId) => {
  const res = await fetch(`${API_URL}/admin/comments/${commentId}`, {
    headers: { 'Authorization': `Bearer ${getToken()}` }
  });
  return (await res.json()).data;
};

// Gỡ/Khôi phục bình luận (admin)
const adminUpdateCommentIsDeleted = async (commentId, isDeleted) => {
  const res = await fetch(`${API_URL}/admin/comments/${commentId}/is-deleted`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getToken()}`
    },
    body: JSON.stringify({ is_deleted: isDeleted })
  });
  return await res.json();
};

// ============ HELPER CONSTANTS ============

const COMPLAINT_CATEGORIES = [
  { value: 'SPAM', label: 'Spam / Quảng cáo' },
  { value: 'HARASSMENT', label: 'Quấy rối / Bắt nạt' },
  { value: 'HATE_SPEECH', label: 'Phát ngôn thù địch' },
  { value: 'VIOLENCE', label: 'Nội dung bạo lực' },
  { value: 'MISINFORMATION', label: 'Thông tin sai lệch' },
  { value: 'INAPPROPRIATE', label: 'Nội dung không phù hợp' },
  { value: 'OTHER', label: 'Lý do khác' }
];

const COMPLAINT_STATUSES = [
  { value: 'PENDING', label: 'Chờ xử lý', color: 'warning' },
  { value: 'REVIEWED', label: 'Đang xem xét', color: 'info' },
  { value: 'RESOLVED', label: 'Đã giải quyết', color: 'success' },
  { value: 'REJECTED', label: 'Từ chối', color: 'error' }
];

const TARGET_TYPES = [
  { value: 'POST', label: 'Bài viết' },
  { value: 'COMMENT', label: 'Bình luận' }
];
```

---

## 4. React Component Examples

### 4.1 Nút báo cáo bài viết/bình luận

```jsx
import { useState } from 'react';

const ReportButton = ({ targetId, targetType, onSuccess }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [category, setCategory] = useState('');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    if (!category) {
      setError('Vui lòng chọn lý do báo cáo');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const result = await createComplaint(targetId, targetType, category, content);
      if (result.error) {
        setError(result.error);
      } else {
        setIsOpen(false);
        onSuccess?.();
      }
    } catch (err) {
      setError('Có lỗi xảy ra, vui lòng thử lại');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button onClick={() => setIsOpen(true)}>
        🚩 Báo cáo
      </button>

      {isOpen && (
        <div className="modal">
          <h3>Báo cáo {targetType === 'POST' ? 'bài viết' : 'bình luận'}</h3>
          
          <select value={category} onChange={(e) => setCategory(e.target.value)}>
            <option value="">-- Chọn lý do --</option>
            {COMPLAINT_CATEGORIES.map(cat => (
              <option key={cat.value} value={cat.value}>{cat.label}</option>
            ))}
          </select>

          <textarea
            placeholder="Mô tả chi tiết (không bắt buộc)"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            maxLength={1000}
          />

          {error && <p className="error">{error}</p>}

          <div className="actions">
            <button onClick={() => setIsOpen(false)}>Hủy</button>
            <button onClick={handleSubmit} disabled={loading}>
              {loading ? 'Đang gửi...' : 'Gửi báo cáo'}
            </button>
          </div>
        </div>
      )}
    </>
  );
};
```

### 4.2 Admin Dashboard - Danh sách khiếu nại

```jsx
import { useState, useEffect } from 'react';

const ComplaintsDashboard = () => {
  const [complaints, setComplaints] = useState([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState('PENDING');
  const [typeFilter, setTypeFilter] = useState('');
  const [loading, setLoading] = useState(false);
  const [counts, setCounts] = useState({});

  // Load counts
  useEffect(() => {
    const loadCounts = async () => {
      const pending = await getComplaintsCount('PENDING');
      const reviewed = await getComplaintsCount('REVIEWED');
      setCounts({
        pending: pending.count,
        reviewed: reviewed.count
      });
    };
    loadCounts();
  }, []);

  // Load complaints
  useEffect(() => {
    const loadComplaints = async () => {
      setLoading(true);
      const data = await getAllComplaints(page, 10, statusFilter, typeFilter);
      setComplaints(data.complaints);
      setTotal(data.total);
      setLoading(false);
    };
    loadComplaints();
  }, [page, statusFilter, typeFilter]);

  const handleUpdateStatus = async (id, newStatus, notes = '') => {
    const result = await updateComplaintStatus(id, newStatus, notes);
    if (!result.error) {
      // Reload list
      const data = await getAllComplaints(page, 10, statusFilter, typeFilter);
      setComplaints(data.complaints);
    }
  };

  return (
    <div className="complaints-dashboard">
      <h1>Quản lý khiếu nại</h1>

      {/* Stats */}
      <div className="stats">
        <div className="stat warning">
          <span>Chờ xử lý</span>
          <strong>{counts.pending || 0}</strong>
        </div>
        <div className="stat info">
          <span>Đang xem xét</span>
          <strong>{counts.reviewed || 0}</strong>
        </div>
      </div>

      {/* Filters */}
      <div className="filters">
        <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
          <option value="">Tất cả trạng thái</option>
          {COMPLAINT_STATUSES.map(s => (
            <option key={s.value} value={s.value}>{s.label}</option>
          ))}
        </select>

        <select value={typeFilter} onChange={(e) => setTypeFilter(e.target.value)}>
          <option value="">Tất cả loại</option>
          {TARGET_TYPES.map(t => (
            <option key={t.value} value={t.value}>{t.label}</option>
          ))}
        </select>
      </div>

      {/* Table */}
      <table>
        <thead>
          <tr>
            <th>Loại</th>
            <th>Lý do</th>
            <th>Nội dung</th>
            <th>Trạng thái</th>
            <th>Ngày tạo</th>
            <th>Hành động</th>
          </tr>
        </thead>
        <tbody>
          {complaints.map(complaint => (
            <tr key={complaint.id}>
              <td>{complaint.target_type}</td>
              <td>{complaint.category}</td>
              <td>{complaint.content || '-'}</td>
              <td>
                <span className={`badge ${complaint.status.toLowerCase()}`}>
                  {COMPLAINT_STATUSES.find(s => s.value === complaint.status)?.label}
                </span>
              </td>
              <td>{new Date(complaint.created_at).toLocaleDateString('vi-VN')}</td>
              <td>
                <button onClick={() => viewTarget(complaint)}>
                  Xem nội dung
                </button>
                {complaint.status === 'PENDING' && (
                  <>
                    <button onClick={() => handleUpdateStatus(complaint.id, 'REVIEWED')}>
                      Xem xét
                    </button>
                    <button onClick={() => handleUpdateStatus(complaint.id, 'RESOLVED', 'Đã xử lý')}>
                      Giải quyết
                    </button>
                    <button onClick={() => handleUpdateStatus(complaint.id, 'REJECTED', 'Không vi phạm')}>
                      Từ chối
                    </button>
                  </>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Pagination */}
      <div className="pagination">
        <button disabled={page === 1} onClick={() => setPage(p => p - 1)}>
          Trước
        </button>
        <span>Trang {page} / {Math.ceil(total / 10)}</span>
        <button disabled={page * 10 >= total} onClick={() => setPage(p => p + 1)}>
          Sau
        </button>
      </div>
    </div>
  );
};
```

---

## 5. Tóm tắt Endpoints

### User (Cần Bearer token):
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/complaints` | Tạo khiếu nại mới |
| GET | `/complaints/my` | Lấy khiếu nại của tôi |
| GET | `/complaints/:id` | Xem chi tiết khiếu nại |
| DELETE | `/complaints/:id` | Xóa khiếu nại (chỉ PENDING) |

### Admin - Quản lý khiếu nại (Cần Admin token):
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/complaints` | Lấy tất cả khiếu nại |
| GET | `/complaints/count` | Đếm số khiếu nại |
| GET | `/complaints/target/:target_id?target_type=` | Lấy khiếu nại theo bài viết/bình luận |
| PUT | `/complaints/:id/status` | Cập nhật trạng thái |
| DELETE | `/complaints/admin/:id` | Xóa khiếu nại |

### Admin - Quản lý bài viết (Cần Admin token):
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/posts/:id` | Xem bài viết (bao gồm cả bài đã xóa) |
| PUT | `/admin/posts/:id/is-deleted` | Gỡ/Khôi phục bài viết |

### Admin - Quản lý bình luận (Cần Admin token):
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/comments/:commentId` | Xem bình luận (bao gồm cả bình luận đã xóa) |
| PUT | `/admin/comments/:commentId/is-deleted` | Gỡ/Khôi phục bình luận |

---

## 6. Workflow xử lý khiếu nại

```
1. User báo cáo bài viết/bình luận
   ↓
2. Tạo complaint với status = PENDING
   ↓
3. Admin xem danh sách PENDING complaints
   ↓
4. Admin xem nội dung bị báo cáo (sử dụng GET /admin/posts/:id hoặc GET /admin/comments/:commentId)
   ↓
5. Admin quyết định:
   ├── RESOLVED: Gỡ nội dung vi phạm (PUT /admin/posts/:id/is-deleted hoặc PUT /admin/comments/:commentId/is-deleted)
   └── REJECTED: Không vi phạm, giữ nguyên
   ↓
6. Cập nhật status + admin_notes
```
