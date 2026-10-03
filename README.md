# Mộc - Auth Backend (Golang)

Hệ thống Authentication Backend hiệu năng cao viết bằng **Golang**, tích hợp cơ chế **Cookie + JWT** với tính năng **Transparent In-Flight Auto-Refresh (Zero-Roundtrip Retry)** và dịch vụ Mail local qua **Mailpit** (Docker).

---

## 1. Điểm Nổi Bật Về Thiết Kế Kiến Trúc

### 🚀 Transparent In-Flight Auto-Refresh (Không cần request thứ 2)
Trong hệ thống truyền thống:
- Khi Access Token hết hạn, API trả về `401 Unauthorized`.
- Client Interceptor phải hứng lỗi, tạm dừng luồng, gửi request `POST /refresh`, rồi retry lại request ban đầu (tổng cộng 2 round-trip mạng, giao diện bị gián đoạn).

**Trong backend Mộc:**
1. Khi client gửi request (ví dụ: `GET /api/auth/me`), nếu Cookie `access_token` đã hết hạn nhưng `refresh_token` còn hợp lệ:
2. **Auth Middleware** lập tức tự động xoay vòng refresh token (**Refresh Token Rotation** - RTR).
3. Đính kèm các header `Set-Cookie` chứa `access_token` và `refresh_token` mới vào ngay Response hiện tại.
4. Nạp thông tin người dùng vào Context và **tiếp tục chuyển tiếp request sang Handler nghiệp vụ ban đầu**.
5. Client nhận được kết quả **HTTP 200 OK ngay trong lần gọi đầu tiên** kèm Cookie mới mà trình duyệt tự động cập nhật!

### 🛡️ Rotation Grace Period (30s) Chống Race Condition
- Khi trang web tải nhiều API đồng thời (`GET /me`, `GET /notifications`), các request sẽ mang cùng refresh token cũ.
- Server áp dụng khoảng đệm **Grace Period 30 giây**: nếu token cũ vừa được xoay vòng tới trong vòng 30s, server vẫn chấp nhận và cấp access token mà không báo lỗi reuse.
- Sau 30s, nếu token cũ tiếp tục xuất hiện, server nhận diện đây là hành vi chiếm đoạt token (Replay Attack) và hủy toàn bộ phiên đăng nhập của người dùng.

### 📬 Local Mail Testing với Mailpit (Docker)
- Không cần cấu hình tài khoản Gmail hay SMTP ngoài tốn kém / dễ bị spam block.
- Mailpit chạy qua Docker với cổng SMTP `1025` và Web UI trực quan tại `http://localhost:8025`.
- Kiểm tra toàn bộ nội dung email đặt lại mật khẩu với giao diện HTML chuẩn phong cách "Mộc".

---

## 2. Cấu Trúc Thư Mục

```
ui-design-backend/
├── cmd/
│   └── server/
│       └── main.go              # Khởi chạy server & Graceful shutdown
├── internal/
│   ├── auth/                    # Bcrypt hash, JWT service & Crypto tokens
│   ├── config/                  # Quản lý cấu hình & biến môi trường
│   ├── domain/                  # Entities (User, Session, PasswordReset) & Interfaces
│   ├── mailer/                  # SMTP client gửi email HTML (Mailpit)
│   ├── middleware/              # AuthMiddleware (Transparent Auto-Refresh), CORS
│   ├── repository/              # Repository cho PostgreSQL & SQLite
│   ├── service/                 # AuthService (Register, Login, Rotate, Password Reset)
│   └── transport/http/          # HTTP Handlers (Register, Login, Logout, Me, Reset)
├── tests/
│   └── auth_flow_test.go        # Kiểm thử tích hợp toàn diện luồng auth & auto-refresh
├── docker-compose.yml           # Mailpit & PostgreSQL
├── .env.example                 # Biến môi trường mẫu
└── .env                         # Cấu hình hiện tại
```

---

## 3. Hướng Dẫn Cài Đặt & Chạy

### 3.1 Khởi động Docker Services (Mailpit & PostgreSQL)
```bash
docker compose up -d
```
- **Mailpit Web UI**: [http://localhost:8025](http://localhost:8025)
- **Mailpit SMTP**: `localhost:1025`
- **PostgreSQL**: `localhost:5432` (database: `uidesign`, user: `postgres`, password: `postgrespassword`)

### 3.2 Khởi chạy Backend Go
Mặc định backend có thể chạy ngay với SQLite (không cần cài thêm DB ngoài):
```bash
go run ./cmd/server
```
Server sẽ lắng nghe tại: [http://localhost:8080](http://localhost:8080).

Nếu muốn dùng PostgreSQL trong Docker, chỉnh sửa file `.env`:
```ini
DB_DRIVER=postgres
DB_DSN=postgres://postgres:postgrespassword@localhost:5432/uidesign?sslmode=disable
```

### 3.3 Chạy Test Suite Tích Hợp
```bash
go test -v -count=1 ./tests/...
```
Bộ test tự động kiểm tra:
1. Đăng ký & Đăng nhập thành công, cấp Cookie bảo mật `HttpOnly`.
2. Token hết hạn -> Tự động refresh in-flight -> Trả về `200 OK` và `Set-Cookie` mới mà không cần retry!
3. Cơ chế Grace Period cho các request đồng thời.
4. Gửi mail đặt lại mật khẩu thành công đến Mailpit container.

---

## 4. API Endpoints

| Method | Endpoint | Xác thực | Mô tả |
| :--- | :--- | :---: | :--- |
| `POST` | `/api/auth/register` | Public | Đăng ký tài khoản mới & thiết lập cookies |
| `POST` | `/api/auth/login` | Public | Đăng nhập tài khoản & thiết lập cookies |
| `POST` | `/api/auth/logout` | Cookie | Đăng xuất, hủy session & xóa cookies |
| `GET` | `/api/auth/me` | Cookie (Auto-refresh) | Lấy thông tin user (tự refresh nếu hết hạn) |
| `POST` | `/api/auth/password-reset` | Public | Gửi email liên kết đặt lại mật khẩu |
| `POST` | `/api/auth/password-reset/confirm`| Public | Cập nhật mật khẩu mới bằng token |
| `GET` | `/health` | Public | Health check |

---

## 5. Kết Nối Frontend (`ui-design`)

Trong dự án `ui-design`, khởi tạo adapter kết nối vào `auth.service.ts`:

```ts
import { createApiClient } from '../../../services/api-client'
import type { AuthAdapter } from '../types'

export function createHttpAuthAdapter(baseUrl = 'http://localhost:8080'): AuthAdapter {
  const client = createApiClient(baseUrl)
  return {
    login: (values, signal) => client.request('/api/auth/login', { method: 'POST', body: JSON.stringify(values), signal }),
    register: (values, signal) => client.request('/api/auth/register', { method: 'POST', body: JSON.stringify(values), signal }),
    requestPasswordReset: (email, signal) => client.request('/api/auth/password-reset', { method: 'POST', body: JSON.stringify({ email }), signal }),
  }
}
```

Nhờ `credentials: 'include'` sẵn có trong `api-client.ts`, toàn bộ cookie sẽ được trình duyệt tự động gửi và nhận, kích hoạt trơn tru trải nghiệm Transparent Auto-Refresh!
