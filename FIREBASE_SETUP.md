# Hướng dẫn kết nối Firebase

Hướng dẫn chi tiết từng bước để kết nối Firebase Firestore với ứng dụng Todo List.

## Bước 1: Tạo Firebase Project

1. Truy cập [Firebase Console](https://console.firebase.google.com/)
2. Click **"Add project"** hoặc **"Create a project"**
3. Nhập tên project (ví dụ: `todo-list-app`)
4. Chọn có bật Google Analytics hay không (tùy chọn)
5. Click **"Create project"** và đợi Firebase tạo project

## Bước 2: Enable Firestore Database

1. Trong Firebase Console, vào menu bên trái
2. Click **"Firestore Database"** (hoặc **"Build" > "Firestore Database"**)
3. Click **"Create database"**
4. Chọn chế độ:
   - **Production mode** (cho production - cần security rules)
   - **Test mode** (cho development - cho phép đọc/ghi trong 30 ngày)
5. Chọn **location** (region) cho database (ví dụ: `asia-southeast1` cho Việt Nam)
6. Click **"Enable"**

## Bước 3: Tạo Service Account

1. Trong Firebase Console, click vào **⚙️ Settings** (bánh răng) ở góc trên bên trái
2. Chọn **"Project settings"**
3. Vào tab **"Service accounts"**
4. Click **"Generate new private key"**
5. Một popup cảnh báo sẽ hiện ra, click **"Generate key"**
6. File JSON sẽ được tải xuống (ví dụ: `todo-list-app-firebase-adminsdk-xxxxx-xxxxxxxxxx.json`)

## Bước 4: Cấu hình ứng dụng

### Cách 1: Sử dụng file .env (Khuyến nghị)

1. **Di chuyển file JSON vào project:**
   ```bash
   # Di chuyển file JSON đã tải về vào thư mục project
   mv ~/Downloads/todo-list-app-firebase-adminsdk-xxxxx.json ./firebase-service-account.json
   ```

2. **Tạo file .env:**
   ```bash
   cp .env.example .env
   ```

3. **Chỉnh sửa file .env:**
   ```env
   FIREBASE_CREDENTIALS=./firebase-service-account.json
   PORT=8080
   ```

### Cách 2: Sử dụng environment variable

1. **Set biến môi trường:**
   ```bash
   export FIREBASE_CREDENTIALS=./firebase-service-account.json
   ```

   Hoặc nếu muốn dùng JSON trực tiếp:
   ```bash
   export GOOGLE_APPLICATION_CREDENTIALS_JSON='{"type":"service_account","project_id":"your-project-id",...}'
   ```

### Cách 3: Application Default Credentials (cho GCP)

```bash
export GOOGLE_APPLICATION_CREDENTIALS=./firebase-service-account.json
```

## Bước 5: Kiểm tra kết nối

1. **Chạy ứng dụng:**
   ```bash
   go run main.go
   ```

2. **Kiểm tra log:**
   - Nếu thấy `Firebase initialized successfully` → Kết nối thành công ✅
   - Nếu có lỗi → Xem phần Troubleshooting bên dưới

3. **Test API:**
   ```bash
   # Tạo một todo mới
   curl -X POST http://localhost:8080/api/todos \
     -H "Content-Type: application/json" \
     -d '{
       "title": "Test Firebase",
       "description": "Kiểm tra kết nối Firebase"
     }'
   ```

4. **Kiểm tra trong Firebase Console:**
   - Vào Firestore Database
   - Bạn sẽ thấy collection `todos` với document mới được tạo

## Cấu trúc file credentials

File JSON service account có dạng:
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

## Troubleshooting

### Lỗi: "Failed to initialize Firebase"

**Nguyên nhân:**
- File credentials không tồn tại hoặc đường dẫn sai
- File credentials không đúng format
- Thiếu quyền truy cập

**Giải pháp:**
1. Kiểm tra đường dẫn file:
   ```bash
   ls -la firebase-service-account.json
   ```
2. Kiểm tra biến môi trường:
   ```bash
   echo $FIREBASE_CREDENTIALS
   ```
3. Kiểm tra file .env:
   ```bash
   cat .env
   ```

### Lỗi: "Permission denied" hoặc "Access denied"

**Nguyên nhân:**
- Service account chưa có quyền truy cập Firestore
- Firestore chưa được enable

**Giải pháp:**
1. Đảm bảo Firestore đã được enable trong Firebase Console
2. Kiểm tra IAM permissions trong Google Cloud Console

### Lỗi: "Collection not found" hoặc "Document not found"

**Nguyên nhân:**
- Collection sẽ được tạo tự động khi có document đầu tiên
- Đây không phải lỗi, chỉ cần tạo document đầu tiên

### Lỗi: "Invalid credentials"

**Nguyên nhân:**
- File JSON bị hỏng hoặc không đúng
- Service account đã bị xóa

**Giải pháp:**
1. Tải lại file credentials từ Firebase Console
2. Đảm bảo file JSON còn nguyên vẹn

## Security Rules (Quan trọng cho Production)

Mặc định Firestore ở chế độ test mode cho phép đọc/ghi trong 30 ngày. Cho production, cần setup security rules:

1. Vào Firestore Database > **Rules**
2. Cập nhật rules:
   ```javascript
   rules_version = '2';
   service cloud.firestore {
     match /databases/{database}/documents {
       // Chỉ cho phép đọc/ghi từ server (qua service account)
       // Client không thể truy cập trực tiếp
       match /{document=**} {
         allow read, write: if false;
       }
     }
   }
   ```

   Hoặc nếu muốn cho phép client truy cập (cần authentication):
   ```javascript
   rules_version = '2';
   service cloud.firestore {
     match /databases/{database}/documents {
       match /todos/{todoId} {
         allow read, write: if request.auth != null;
       }
     }
   }
   ```

## Kiểm tra kết nối thành công

Sau khi setup xong, bạn có thể kiểm tra:

1. **Log khi khởi động:**
   ```
   Firebase initialized successfully
   Server starting on port :8080
   ```

2. **Test tạo todo:**
   ```bash
   curl -X POST http://localhost:8080/api/todos \
     -H "Content-Type: application/json" \
     -d '{"title":"Test","description":"Test"}'
   ```

3. **Kiểm tra trong Firebase Console:**
   - Vào Firestore Database
   - Thấy collection `todos` với document mới

## Lưu ý bảo mật

⚠️ **QUAN TRỌNG:**
- ❌ **KHÔNG** commit file `firebase-service-account.json` lên Git
- ❌ **KHÔNG** commit file `.env` lên Git
- ✅ File `.gitignore` đã được cấu hình để bảo vệ các file này
- ✅ Chỉ share credentials với team members cần thiết
- ✅ Sử dụng environment variables trong production (không dùng file)

## Next Steps

Sau khi kết nối thành công:
1. Test các API endpoints
2. Setup Firestore security rules cho production
3. Cân nhắc thêm Firebase Authentication nếu cần
4. Setup monitoring và logging

