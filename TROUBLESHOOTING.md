# Troubleshooting - Lỗi "project id is required"

## Vấn đề

Lỗi: `Failed to initialize Firebase: project id is required to access Firestore`

## Nguyên nhân có thể

1. **Code mới chưa được deploy** - Railway vẫn đang chạy code cũ
2. **Environment variable chưa được set** - `GOOGLE_APPLICATION_CREDENTIALS_JSON` không có hoặc rỗng
3. **JSON không hợp lệ** - Format JSON sai hoặc thiếu field `project_id`
4. **JSON bị escape sai** - Quotes bị escape không đúng

## Giải pháp từng bước

### Bước 1: Kiểm tra code đã được deploy chưa

1. **Commit và push code mới:**
   ```bash
   git add .
   git commit -m "Fix Firebase initialization with project_id"
   git push
   ```

2. **Trong Railway:**
   - Vào **Deployments** tab
   - Kiểm tra commit mới nhất có phải là code mới không
   - Nếu không, click **"Redeploy"** hoặc **"Deploy"**

### Bước 2: Kiểm tra Environment Variable

1. **Trong Railway Dashboard:**
   - Vào **Settings** > **Variables**
   - Kiểm tra có biến `GOOGLE_APPLICATION_CREDENTIALS_JSON` không
   - Kiểm tra giá trị có rỗng không

2. **Kiểm tra logs:**
   - Vào **Deployments** > **View Logs**
   - Tìm dòng: `GOOGLE_APPLICATION_CREDENTIALS_JSON: found (X chars)`
   - Nếu thấy `not set` → Environment variable chưa được set

### Bước 3: Kiểm tra JSON hợp lệ

1. **Local test:**
   ```bash
   # Kiểm tra file JSON có project_id không
   cat firebase-service-account.json | grep project_id
   ```

2. **Validate JSON:**
   ```bash
   # Kiểm tra JSON hợp lệ
   cat firebase-service-account.json | python3 -m json.tool
   ```

3. **Kiểm tra trong Railway:**
   - Copy toàn bộ nội dung file JSON
   - Paste vào Railway environment variable
   - **Lưu ý:** Không cần escape quotes, paste trực tiếp

### Bước 4: Xem logs chi tiết

Sau khi deploy code mới, logs sẽ hiển thị:

```
Checking Firebase credentials...
FIREBASE_CREDENTIALS: 
GOOGLE_APPLICATION_CREDENTIALS_JSON: found (1234 chars)
Using JSON credentials from environment variable
JSON parsed successfully, found 10 keys
Found project_id: your-project-id
Initializing Firebase with project_id: your-project-id
Firebase initialized successfully
```

**Nếu thấy lỗi:**
- `GOOGLE_APPLICATION_CREDENTIALS_JSON: not set` → Set environment variable
- `Error parsing JSON credentials` → JSON không hợp lệ
- `project_id not found` → JSON thiếu field project_id
- `Available keys: [...]` → Xem các keys có trong JSON

## Cách set Environment Variable đúng

### Trong Railway:

1. Vào **Settings** > **Variables**
2. Click **"New Variable"**
3. **Name:** `GOOGLE_APPLICATION_CREDENTIALS_JSON`
4. **Value:** Paste toàn bộ nội dung file JSON (từ `{` đến `}`)
5. Click **"Add"**

### Format JSON đúng:

```json
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "...",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-xxxxx@your-project-id.iam.gserviceaccount.com",
  "client_id": "...",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "..."
}
```

**Lưu ý:**
- ✅ Có thể là nhiều dòng hoặc một dòng
- ✅ Không cần escape quotes trong Railway
- ✅ Phải có field `project_id`
- ❌ Không được thiếu dấu `{` hoặc `}`

## Test local trước khi deploy

1. **Set environment variable:**
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS_JSON='$(cat firebase-service-account.json)'
   ```

2. **Chạy app:**
   ```bash
   go run main.go
   ```

3. **Kiểm tra logs:**
   - Phải thấy: `Found project_id: your-project-id`
   - Phải thấy: `Firebase initialized successfully`

## Nếu vẫn lỗi

1. **Kiểm tra file JSON gốc:**
   ```bash
   # Xem project_id
   cat firebase-service-account.json | jq .project_id
   
   # Validate JSON
   cat firebase-service-account.json | jq .
   ```

2. **Tạo lại Service Account:**
   - Vào Firebase Console > Project Settings > Service Accounts
   - Generate new private key
   - Tải file mới và thử lại

3. **Kiểm tra Firestore đã enable:**
   - Vào Firebase Console > Firestore Database
   - Đảm bảo database đã được tạo

4. **Xem logs đầy đủ trong Railway:**
   - Copy toàn bộ logs
   - Tìm các dòng có "Error" hoặc "Failed"
   - Gửi logs để được hỗ trợ

## Checklist

- [ ] Code mới đã được commit và push
- [ ] Railway đã deploy code mới
- [ ] Environment variable `GOOGLE_APPLICATION_CREDENTIALS_JSON` đã được set
- [ ] JSON có field `project_id`
- [ ] JSON hợp lệ (có thể parse được)
- [ ] Firestore đã được enable trong Firebase Console
- [ ] Service Account có quyền truy cập Firestore

## Liên hệ hỗ trợ

Nếu vẫn không giải quyết được, cung cấp:
1. Logs từ Railway (toàn bộ)
2. Screenshot của Environment Variables (ẩn giá trị)
3. Kết quả của `cat firebase-service-account.json | jq .project_id`

