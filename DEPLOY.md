# Hướng dẫn Deploy Ứng dụng Miễn phí

Hướng dẫn deploy ứng dụng Todo List Backend lên các platform miễn phí.

## Các Platform Miễn phí

1. **Railway** ⭐ (Khuyến nghị - Dễ nhất)
2. **Render**
3. **Fly.io**
4. **Google Cloud Run**
5. **Vercel** (Serverless)

---

## 1. Railway (Khuyến nghị)

### Ưu điểm:
- ✅ Free tier: $5 credit/tháng
- ✅ Dễ setup, tự động deploy từ GitHub
- ✅ Hỗ trợ environment variables
- ✅ Auto HTTPS

### Các bước:

1. **Tạo tài khoản:**
   - Truy cập [railway.app](https://railway.app)
   - Đăng nhập bằng GitHub

2. **Tạo project mới:**
   - Click "New Project"
   - Chọn "Deploy from GitHub repo"
   - Chọn repository của bạn

3. **Cấu hình:**
   - Railway tự động detect Go project
   - Thêm environment variables:
     ```
     FIREBASE_CREDENTIALS_JSON=<paste toàn bộ nội dung file JSON>
     PORT=8080
     ```
   - Hoặc upload file credentials (không khuyến nghị)

4. **Deploy:**
   - Railway tự động build và deploy
   - Lấy URL từ Settings > Domains

### Lưu ý:
- Cần set `FIREBASE_CREDENTIALS_JSON` thay vì `FIREBASE_CREDENTIALS` (vì không có file system)
- Cập nhật code để hỗ trợ JSON từ env var

---

## 2. Render

### Ưu điểm:
- ✅ Free tier: 750 giờ/tháng
- ✅ Auto deploy từ GitHub
- ✅ Free SSL

### Các bước:

1. **Tạo tài khoản:**
   - Truy cập [render.com](https://render.com)
   - Đăng nhập bằng GitHub

2. **Tạo Web Service:**
   - Click "New +" > "Web Service"
   - Connect GitHub repository
   - Cấu hình:
     - **Name:** `todo-list-api`
     - **Environment:** `Go`
     - **Build Command:** `go mod download && go build -o todoapp main.go`
     - **Start Command:** `./todoapp`

3. **Environment Variables:**
   - Thêm trong Settings > Environment:
     ```
     GOOGLE_APPLICATION_CREDENTIALS_JSON=<paste JSON content>
     PORT=8080
     ```

4. **Deploy:**
   - Click "Create Web Service"
   - Render tự động deploy

### Lưu ý:
- Free tier sẽ sleep sau 15 phút không có traffic
- Lần request đầu tiên sau khi sleep sẽ chậm (~30s)

---

## 3. Fly.io

### Ưu điểm:
- ✅ Free tier: 3 VMs nhỏ
- ✅ Global edge network
- ✅ Không sleep

### Các bước:

1. **Cài đặt Fly CLI:**
   ```bash
   curl -L https://fly.io/install.sh | sh
   ```

2. **Đăng nhập:**
   ```bash
   fly auth login
   ```

3. **Tạo file `fly.toml`:**
   ```toml
   app = "your-app-name"
   primary_region = "sin"  # Singapore

   [build]

   [env]
     PORT = "8080"

   [[services]]
     internal_port = 8080
     protocol = "tcp"

     [[services.ports]]
       handlers = ["http"]
       port = 80

     [[services.ports]]
       handlers = ["tls", "http"]
       port = 443

     [[services.http_checks]]
       interval = "10s"
       timeout = "2s"
       grace_period = "5s"
       method = "GET"
       path = "/health"
   ```

4. **Deploy:**
   ```bash
   fly launch
   fly secrets set GOOGLE_APPLICATION_CREDENTIALS_JSON="<JSON content>"
   fly deploy
   ```

---

## 4. Google Cloud Run

### Ưu điểm:
- ✅ Free tier: 2 triệu requests/tháng
- ✅ Chỉ trả tiền khi có request
- ✅ Tích hợp tốt với Firebase

### Các bước:

1. **Cài đặt Google Cloud SDK:**
   ```bash
   # macOS
   brew install google-cloud-sdk
   
   # Hoặc download từ: https://cloud.google.com/sdk/docs/install
   ```

2. **Đăng nhập và setup:**
   ```bash
   gcloud init
   gcloud auth login
   ```

3. **Tạo Dockerfile:**
   ```dockerfile
   FROM golang:1.21-alpine AS builder
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN go build -o todoapp main.go

   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/todoapp .
   EXPOSE 8080
   CMD ["./todoapp"]
   ```

4. **Build và deploy:**
   ```bash
   # Set project
   gcloud config set project YOUR_PROJECT_ID
   
   # Build container
   gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/todoapp
   
   # Deploy
   gcloud run deploy todoapp \
     --image gcr.io/YOUR_PROJECT_ID/todoapp \
     --platform managed \
     --region asia-southeast1 \
     --allow-unauthenticated \
     --set-env-vars="GOOGLE_APPLICATION_CREDENTIALS_JSON=<JSON>"
   ```

---

## 5. Vercel (Serverless)

### Ưu điểm:
- ✅ Free tier: 100GB bandwidth/tháng
- ✅ Serverless, scale tự động
- ✅ Edge network

### Cần điều chỉnh code:

Vercel yêu cầu serverless function. Cần tạo `api/index.go` hoặc sử dụng adapter.

**Tạo `vercel.json`:**
```json
{
  "builds": [
    {
      "src": "main.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "main.go"
    }
  ]
}
```

**Deploy:**
```bash
npm i -g vercel
vercel
```

---

## Cập nhật Code để Hỗ trợ Deploy

Cần cập nhật `firebase/firebase.go` để hỗ trợ JSON từ environment variable tốt hơn:

```go
// Đã có sẵn trong code hiện tại
if jsonCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON"); jsonCreds != "" {
    opt := option.WithCredentialsJSON([]byte(jsonCreds))
    app, err = firebase.NewApp(ctx, nil, opt)
}
```

---

## Chuẩn bị trước khi Deploy

### 1. Tạo file `.env.production` (local testing):
```env
GOOGLE_APPLICATION_CREDENTIALS_JSON={"type":"service_account",...}
PORT=8080
```

### 2. Test build local:
```bash
go build -o todoapp main.go
./todoapp
```

### 3. Chuẩn bị Firebase credentials:
- Copy toàn bộ nội dung file JSON
- Sẵn sàng paste vào environment variables

### 4. Đảm bảo `.gitignore` đã có:
```
.env
firebase-service-account.json
*.json
```

---

## So sánh các Platform

| Platform | Free Tier | Sleep | Setup | Tốc độ |
|----------|-----------|-------|-------|--------|
| **Railway** | $5 credit/tháng | ❌ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Render** | 750h/tháng | ✅ (15 phút) | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Fly.io** | 3 VMs | ❌ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Cloud Run** | 2M requests | ❌ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Vercel** | 100GB bandwidth | ❌ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

---

## Khuyến nghị

### Cho Development/Testing:
- **Railway** - Dễ nhất, không sleep

### Cho Production nhỏ:
- **Fly.io** - Tốc độ tốt, không sleep
- **Cloud Run** - Tích hợp tốt với Firebase

### Cho Production lớn:
- **Cloud Run** - Scale tốt, giá rẻ

---

## Troubleshooting

### Lỗi: "Failed to initialize Firebase"

**Nguyên nhân:** Environment variable không đúng format

**Giải pháp:**
- Đảm bảo JSON là một dòng (escape quotes)
- Hoặc sử dụng base64 encoding:
  ```bash
  cat firebase-service-account.json | base64
  ```
  Và decode trong code

### Lỗi: "Port already in use"

**Nguyên nhân:** Platform tự động set PORT

**Giải pháp:**
- Đảm bảo code đọc `PORT` từ environment:
  ```go
  port := os.Getenv("PORT")
  if port == "" {
      port = "8080"
  }
  ```

### Lỗi: "Build failed"

**Nguyên nhân:** Thiếu dependencies

**Giải pháp:**
- Đảm bảo `go.mod` và `go.sum` đã commit
- Test build local trước:
  ```bash
  go mod tidy
  go build
  ```

---

## Next Steps sau khi Deploy

1. ✅ Test API endpoints
2. ✅ Setup custom domain (nếu cần)
3. ✅ Setup monitoring (nếu platform hỗ trợ)
4. ✅ Setup CI/CD tự động deploy từ GitHub
5. ✅ Review security rules trong Firestore

---

## Lưu ý Bảo mật

⚠️ **QUAN TRỌNG:**
- ❌ **KHÔNG** commit file credentials lên Git
- ✅ Sử dụng environment variables trong platform
- ✅ Rotate credentials định kỳ
- ✅ Sử dụng Firestore security rules
- ✅ Giới hạn IP access nếu cần (cho production)

