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
| `POST` | `/api/auth/admin/login` | Public | Đăng nhập riêng cho admin hoặc moderator |
| `POST` | `/api/auth/password-reset` | Public | Gửi email liên kết đặt lại mật khẩu |
| `POST` | `/api/auth/password-reset/validate` | Public | Kiểm tra token còn hiệu lực trước khi hiển thị form, không tiêu thụ token |
| `POST` | `/api/auth/password-reset/confirm`| Public | Cập nhật mật khẩu mới bằng token |
| `GET` | `/api/profile/me` | Cookie (Auto-refresh) | Lấy hồ sơ cá nhân của user đăng nhập (tự tạo mặc định nếu chưa có) |
| `PUT` | `/api/profile/me` | Cookie (Auto-refresh) | Cập nhật hồ sơ cá nhân (avatar, banner, bio, themeColor,...) |
| `GET` | `/api/profile/{userId}` | Public | Xem hồ sơ công khai của người dùng khác |
| `POST` | `/api/upload` | Cookie (Auto-refresh) | Tải file (ảnh avatar, banner) lên thư mục cục bộ (`uploads/`) |
| `GET` | `/uploads/{filename}` | Public | Truy cập file tĩnh đã tải lên |
| `GET` | `/health` | Public | Health check |
| `GET` | `/api/stories` | Public | Truyện đã đăng/hoàn thành với ít nhất 2 chương có nội dung |
| `GET` | `/api/stories?reading=true` | Cookie (Auto-refresh) | Truyện trong lịch sử đọc của tài khoản, mới đọc trước |
| `GET` | `/api/stories/mine` | Cookie (Auto-refresh) | Toàn bộ truyện của tác giả, bao gồm bản nháp và chương rỗng |
| `GET` | `/api/stories/{id}` | Public | Thông tin truyện và nội dung các chương theo thứ tự |
| `GET` | `/api/stories/{id}/reading-progress` | Cookie (Auto-refresh) | Chương, trang và điểm trong nội dung đọc gần nhất; chưa đọc có `chapterId` rỗng |
| `PUT` | `/api/stories/{id}/reading-progress` | Cookie (Auto-refresh) | Lưu `chapterId`, `mode`, `pageIndex`, `blockIndex`, `textOffset`, `scrollProgress`, trả `204`; chương phải thuộc truyện và có nội dung |
| `GET` | `/api/admin/stories` | Staff | Danh sách toàn bộ truyện, gồm cả bản nháp và truyện ít chương |
| `GET` | `/api/admin/users` | Admin | Danh sách người dùng |
| `PUT` | `/api/admin/users/{id}/role` | Admin | Đặt role `admin` hoặc `moderator` |
| `DELETE` | `/api/admin/users/{id}` | Admin | Xóa người dùng cùng dữ liệu truyện, chương, hồ sơ và phiên |

## Quản trị

Trang quản trị frontend nằm tại `/admin/login`. Backend chỉ chấp nhận tài khoản có `role` là `admin` hoặc `moderator` ở endpoint này. Moderator quản lý được toàn bộ truyện: sửa thông tin, chương, ảnh, trạng thái và xóa truyện. Admin có thêm danh sách người dùng và được đặt role hoặc xóa người dùng.

Tài khoản đăng ký mới không có role quản trị. Để bootstrap admin khi chạy server, đặt email của một tài khoản đã tồn tại trong `.env`:

```ini
ADMIN_EMAIL=admin@example.com
```

Server sẽ đặt tài khoản đó thành `admin` mỗi lần khởi động. Sau đó admin có thể cấp `moderator` cho tài khoản khác trong `/admin/users`. Role được đọc lại từ database ở mỗi request được bảo vệ nên việc thu hồi quyền có hiệu lực ngay ở request kế tiếp.
Sau khi bootstrap, có thể xóa `ADMIN_EMAIL` khỏi `.env` để role được quản lý hoàn toàn trong trang quản trị.

Khi khởi động, backend tự chạy migration `reader-v1` trong transaction: thêm `chapters.has_content`, tính lại cờ từ nội dung hiện có và tạo bảng `reading_progress` với khóa `(user_id, story_id)`. Migration không xóa truyện và có thể chạy lại. Cờ nội dung được cập nhật khi tạo/sửa chương; HTML rỗng, khoảng trắng, script/style không được tính là nội dung. Chương chỉ có ảnh hợp lệ được tính là có nội dung. Khi xóa chương hoặc truyện, lịch sử liên quan được xóa cùng trong transaction.

Sau khi cập nhật mã nguồn, khởi động lại server bằng `go run ./cmd/server` để áp dụng API và migration mới. Không cần chạy lại seed.

Các API danh sách truyện trả thêm `previewImages`: tối đa 3 URL ảnh khác nhau, ảnh bìa trước rồi đến ảnh minh họa trong chương theo thứ tự chương. Ảnh nhúng trong HTML cũng được lấy khi chương chưa có danh sách `illustrations`. Không nhân bản ảnh để đủ 3 lớp, không trả toàn bộ nội dung chương trong danh sách và không cần thêm migration cho trường này.

`GET /api/stories?reading=true` trả thêm `readingProgress` trên từng truyện: `chapterId`, `chapterNumber`, `chapterTitle` và `updatedAt` (thời điểm đọc UTC). Danh sách được sắp xếp theo `reading_progress.updated_at` giảm dần, không theo ngày tạo/sửa truyện. Truyện đã đọc còn 1 chương vẫn có trong danh sách này; bản nháp chỉ xuất hiện với tác giả. `GET /api/stories/{id}/reading-progress` cũng trả thông tin chương và thời gian tương tự. Mỗi lần mở chương có nội dung, frontend gọi PUT; backend tự lấy thời gian hiện tại và cập nhật database kể cả khi đọc lại cùng chương. Bảng hiện tại đã lưu các dữ liệu cần thiết, không cần migration hay seed lại database.

`/password-reset/validate` nhận JSON `{ "token": "..." }`, trả `204` khi token còn hiệu lực và `400` nếu token sai, đã dùng hoặc hết hạn. Response không được cache. `/password-reset/confirm` tiêu thụ token, cập nhật mật khẩu và thu hồi các phiên đăng nhập trong một transaction; mỗi token chỉ có thể hoàn tất một lần, kể cả khi có các yêu cầu đồng thời.

`/api/upload` nhận multipart/form-data với dung lượng tối đa 10MB (hỗ trợ JPG, PNG, WEBP, GIF, SVG) và trả về URL tĩnh dạng `/uploads/<uuid>.<ext>`.


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
