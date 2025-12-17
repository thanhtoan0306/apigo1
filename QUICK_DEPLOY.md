# Quick Deploy - Railway (5 phút)

Hướng dẫn nhanh deploy ứng dụng lên Railway (miễn phí).

## Bước 1: Chuẩn bị

1. **Có tài khoản GitHub** và code đã push lên repository
2. **Có Firebase credentials** (file JSON)

## Bước 2: Deploy trên Railway

1. **Truy cập:** [railway.app](https://railway.app)
2. **Đăng nhập** bằng GitHub
3. **Click "New Project"** > **"Deploy from GitHub repo"**
4. **Chọn repository** của bạn
5. Railway tự động detect và build Go project

## Bước 3: Cấu hình Environment Variables

1. Vào **Settings** > **Variables**
2. Thêm biến môi trường:

   **Tên:** `GOOGLE_APPLICATION_CREDENTIALS_JSON`
   
   **Giá trị:** Copy toàn bộ nội dung file `firebase-service-account.json` và paste vào
   
   ⚠️ **Lưu ý quan trọng:** 
   - Phải là JSON hợp lệ (một dòng hoặc nhiều dòng đều được)
   - JSON phải có field `project_id` (code sẽ tự động parse)
   - Không cần escape quotes khi paste vào Railway
   - Ví dụ format:
     ```json
     {"type":"service_account","project_id":"your-project-id","private_key_id":"...","private_key":"-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n","client_email":"...","client_id":"...","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url":"..."}
     ```

3. (Tùy chọn) Thêm `PORT=8080` nếu cần

## Bước 4: Lấy URL

1. Vào tab **Settings** > **Domains**
2. Railway tự động tạo domain: `your-app-name.up.railway.app`
3. Copy URL này

## Bước 5: Test

```bash
# Test health check
curl https://your-app-name.up.railway.app/health

# Test tạo todo
curl -X POST https://your-app-name.up.railway.app/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","description":"Test deploy"}'
```

## ✅ Xong!

Ứng dụng của bạn đã được deploy và có thể truy cập từ bất kỳ đâu.

---

## Troubleshooting

### Lỗi: "Failed to initialize Firebase" hoặc "project id is required"

**Nguyên nhân:**
- JSON không hợp lệ
- Thiếu field `project_id` trong JSON
- JSON bị format sai

**Giải pháp:**
1. Kiểm tra JSON có field `project_id`:
   ```bash
   # Local test
   cat firebase-service-account.json | grep project_id
   ```
2. Đảm bảo copy toàn bộ nội dung file JSON (từ `{` đến `}`)
3. Trong Railway, paste JSON vào textarea (không cần escape)
4. Kiểm tra logs trong Railway để xem lỗi chi tiết

### Lỗi: "Build failed"
- Kiểm tra `go.mod` và `go.sum` đã được commit
- Xem logs trong Railway dashboard

### App không chạy
- Kiểm tra logs trong Railway dashboard
- Đảm bảo PORT được set đúng (Railway tự động set)

---

## Xem hướng dẫn chi tiết

📖 [DEPLOY.md](./DEPLOY.md) - Hướng dẫn đầy đủ cho tất cả platforms

