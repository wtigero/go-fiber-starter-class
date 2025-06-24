# 🔐 JWT Authentication - แบบง่าย ๆ ที่เข้าใจได้ (60 นาที)

## 📚 จุดประสงค์
เรียนรู้การทำระบบ Login/Register ด้วย JWT ในแบบที่เข้าใจง่าย ไม่ซับซ้อน

## 🎯 สิ่งที่จะเรียนรู้
- สร้างระบบ Register (สมัครสมาชิก)
- สร้างระบบ Login (เข้าสู่ระบบ)
- สร้าง JWT Token
- ตรวจสอบ Token ก่อนเข้า API
- Hash Password ให้ปลอดภัย

## 📋 API Endpoints
- `POST /register` - สมัครสมาชิก
- `POST /login` - เข้าสู่ระบบ
- `GET /profile` - ดูข้อมูลตัวเอง (ต้อง login)
- `GET /todos` - ดู todos ของตัวเอง (ต้อง login)

## 📊 ข้อมูลที่ใช้
```go
type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Password string `json:"-"`        // ไม่ส่งใน response
    Name     string `json:"name"`
}

type Todo struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Done   bool   `json:"done"`
    UserID int    `json:"user_id"`  // เก็บว่าเป็น todo ของใคร
}
```

## 🏃‍♂️ วิธีรัน

### สำหรับนักเรียน:
```bash
cd starter
go mod tidy
go run main.go
```

### สำหรับดูเฉลย:
```bash
cd complete  
go mod tidy
go run main.go
```

## 🧪 ทดสอบทีละขั้นตอน

### 1. สมัครสมาชิก
```bash
curl -X POST http://localhost:3000/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456",
    "name": "นายทดสอบ"
  }'
```

**ผลลัพธ์:**
```json
{
  "success": true,
  "message": "สมัครสมาชิกสำเร็จ",
  "user": {
    "id": 1,
    "email": "test@example.com", 
    "name": "นายทดสอบ"
  }
}
```

### 2. เข้าสู่ระบบ
```bash
curl -X POST http://localhost:3000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "123456"
  }'
```

**ผลลัพธ์:**
```json
{
  "success": true,
  "message": "เข้าสู่ระบบสำเร็จ",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "test@example.com",
    "name": "นายทดสอบ"
  }
}
```

### 3. ดูโปรไฟล์ (ต้องใส่ token)
```bash
curl http://localhost:3000/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### 4. ดู todos ของตัวเอง
```bash
curl http://localhost:3000/todos \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

## 🔍 สิ่งสำคัญที่เรียนรู้

### 1. การ Hash Password
```go
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

### 2. การสร้าง JWT Token
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": user.ID,
    "email": user.Email,
    "exp": time.Now().Add(time.Hour * 24).Unix(),
})
```

### 3. การตรวจสอบ Token
```go
func RequireAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{"error": "ต้อง login ก่อน"})
        }
        // ตรวจสอบ token...
    }
}
```

## 📝 ใน starter/ จะมี:
- [ ] User struct และ Todo struct
- [ ] TODO: สร้าง register endpoint
- [ ] TODO: สร้าง login endpoint  
- [ ] TODO: สร้าง middleware ตรวจ token
- [ ] TODO: ใช้ bcrypt hash password

## ✅ ใน complete/ จะมี:
- ✅ ระบบ register/login สมบูรณ์
- ✅ JWT token generation
- ✅ Middleware ตรวจสอบ token
- ✅ Password hashing
- ✅ Error handling ครบถ้วน

## ⏭️ หลังจากนี้
- ลองทำ Rate Limiting (02-rate-limit-cache)
- หรือลองทำ Database Advanced (06-database-advanced)

---
**เวลาเรียน:** 60 นาที | **ความยาก:** ⭐⭐⭐☆☆
**เหมาะสำหรับ:** คนที่ต้องการทำระบบ login ใน app จริง 