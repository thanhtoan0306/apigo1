# Todo List Backend API

Backend API cho ứng dụng Todo List được xây dựng bằng Golang với Firebase Firestore.

## Yêu cầu

- Go 1.21 hoặc cao hơn
- Firebase project với Firestore enabled
- Firebase service account credentials

## Cài đặt Firebase

📖 **Xem hướng dẫn chi tiết:** [FIREBASE_SETUP.md](./FIREBASE_SETUP.md)

Tóm tắt nhanh:
1. Tạo Firebase project tại [Firebase Console](https://console.firebase.google.com/)
2. Enable Firestore Database
3. Tạo Service Account:
   - Vào Project Settings > Service Accounts
   - Click "Generate new private key"
   - Lưu file JSON

## Cài đặt

1. Tải dependencies:
```bash
go mod download
```

2. Cấu hình environment variables:

   **Sử dụng file .env (khuyến nghị):**
   ```bash
   cp .env.example .env
   # Chỉnh sửa file .env với thông tin Firebase của bạn
   ```

   Hoặc **set environment variables trực tiếp** (chọn một trong các cách):

   **Cách 1: Sử dụng file credentials**
   ```bash
   export FIREBASE_CREDENTIALS=./firebase-service-account.json
   ```

   **Cách 2: Sử dụng environment variable**
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS_JSON='{"type":"service_account",...}'
   ```

   **Cách 3: Application Default Credentials (cho GCP)**
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS=./firebase-service-account.json
   ```

3. Chạy server:
```bash
go run main.go
```

Server sẽ chạy tại `http://localhost:8080`

## API Endpoints

### Health Check
- **GET** `/health` - Kiểm tra trạng thái server

### Todos

- **GET** `/api/todos` - Lấy tất cả todos
- **GET** `/api/todos/{id}` - Lấy todo theo ID
- **POST** `/api/todos` - Tạo todo mới
- **PUT** `/api/todos/{id}` - Cập nhật todo
- **DELETE** `/api/todos/{id}` - Xóa todo

### Blogs (Markdown)

- **GET** `/api/blogs` - Lấy tất cả blogs
- **GET** `/api/blogs/{id}` - Lấy blog theo ID
- **GET** `/api/blogs/slug/{slug}` - Lấy blog theo slug
- **POST** `/api/blogs` - Tạo blog mới
- **PUT** `/api/blogs/{id}` - Cập nhật blog
- **DELETE** `/api/blogs/{id}` - Xóa blog

## Ví dụ sử dụng

### Tạo todo mới
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hoàn thành bài tập",
    "description": "Làm bài tập về nhà"
  }'
```

### Lấy tất cả todos
```bash
curl http://localhost:8080/api/todos
```

### Lấy todo theo ID
```bash
curl http://localhost:8080/api/todos/1
```

### Cập nhật todo
```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Đã hoàn thành bài tập",
    "completed": true
  }'
```

### Xóa todo
```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## Ví dụ sử dụng Blog API

### Tạo blog mới
```bash
curl -X POST http://localhost:8080/api/blogs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hướng dẫn Golang",
    "content": "# Golang\n\nGolang là ngôn ngữ lập trình...",
    "slug": "huong-dan-golang",
    "author": "John Doe",
    "published": true,
    "tags": ["golang", "programming", "tutorial"]
  }'
```

### Lấy tất cả blogs
```bash
curl http://localhost:8080/api/blogs
```

### Lấy blog theo ID
```bash
curl http://localhost:8080/api/blogs/1
```

### Lấy blog theo slug
```bash
curl http://localhost:8080/api/blogs/slug/huong-dan-golang
```

### Cập nhật blog
```bash
curl -X PUT http://localhost:8080/api/blogs/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Hướng dẫn Golang - Updated",
    "published": false
  }'
```

### Xóa blog
```bash
curl -X DELETE http://localhost:8080/api/blogs/1
```

## Cấu trúc dự án

```
apigo1/
├── main.go                    # Entry point của ứng dụng
├── go.mod                     # Go module file
├── .env.example               # Ví dụ cấu hình environment variables
├── .gitignore                 # Git ignore file
├── models/
│   ├── todo.go                # Todo model và request structs
│   └── blog.go                # Blog model và request structs
├── store/
│   ├── store_interface.go     # Interface cho todo store
│   ├── store.go               # In-memory store (backup)
│   ├── firestore_store.go     # Firestore store implementation cho todos
│   └── blog_store.go          # Firestore store implementation cho blogs
├── firebase/
│   └── firebase.go            # Firebase initialization
└── handlers/
    ├── todo_handler.go        # HTTP handlers cho todo endpoints
    └── blog_handler.go        # HTTP handlers cho blog endpoints
```

## Lưu trữ dữ liệu

- Dữ liệu được lưu trữ trong **Firebase Firestore**
- Dữ liệu được lưu vĩnh viễn và có thể truy cập từ bất kỳ đâu
- Collection names: 
  - `todos` - Lưu trữ todos
  - `blogs` - Lưu trữ blogs (Markdown content)

## Deploy

🚀 **Hướng dẫn deploy miễn phí:** [DEPLOY.md](./DEPLOY.md)

Các platform được hỗ trợ:
- Railway (Khuyến nghị - Dễ nhất)
- Render
- Fly.io
- Google Cloud Run
- Vercel

## Lưu ý

- Đảm bảo file credentials không được commit lên Git (đã có trong .gitignore)
- Firestore cần được enable trong Firebase Console
- Cần set up Firestore security rules phù hợp cho production

# apigo1
