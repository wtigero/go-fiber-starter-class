# 🔰 Beginner Level - Go Fiber พื้นฐาน (30 นาที)

## 📚 จุดประสงค์
เรียนรู้การสร้าง REST API พื้นฐานด้วย Go Fiber framework สำหรับผู้เริ่มต้น

## 🎯 สิ่งที่จะเรียนรู้
- การสร้าง Fiber application
- การจัดการ HTTP Routes (GET, POST)
- การใช้ Middleware เบื้องต้น (Logger, Recover)
- การทำงานกับ JSON data
- การเก็บข้อมูลใน In-memory storage

## 📋 API Endpoints ที่จะสร้าง
- `GET /hello` - Hello World endpoint
- `GET /todos` - ดูรายการ Todo ทั้งหมด
- `GET /todos/:id` - ดูรายการ Todo เดียว
- `POST /todos` - เพิ่มรายการ Todo ใหม่

## 📊 Data Structure
```go
type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}
```

## 🏃‍♂️ วิธีรัน

### สำหรับนักเรียน (starter):
```bash
cd starter
go mod tidy
go run main.go
```

### สำหรับผู้สอน (complete):
```bash
cd complete
go mod tidy
go run main.go
```

Server จะรันที่ `http://localhost:3000`

## 🧪 ทดสอบ API

### 1. Hello endpoint
```bash
curl http://localhost:3000/hello
```
**ผลลัพธ์:**
```json
{
  "message": "Hello from Go Fiber Todo API!",
  "version": "1.0.0"
}
```

### 2. ดูรายการ todos ทั้งหมด
```bash
curl http://localhost:3000/todos
```
**ผลลัพธ์:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "เรียน Go Programming",
      "done": false
    },
    {
      "id": 2,
      "title": "ทำโปรเจกต์ Todo API",
      "done": true
    },
    {
      "id": 3,
      "title": "ทบทวน Fiber Framework",
      "done": false
    }
  ],
  "count": 3
}
```

### 3. ดู todo เดียว
```bash
curl http://localhost:3000/todos/1
```
**ผลลัพธ์:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "เรียน Go Programming",
    "done": false
  }
}
```

### 4. เพิ่ม todo ใหม่
```bash
curl -X POST http://localhost:3000/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "เรียน Go Fiber",
    "done": false
  }'
```
**ผลลัพธ์:**
```json
{
  "success": true,
  "message": "สร้าง todo สำเร็จ",
  "data": {
    "id": 4,
    "title": "เรียน Go Fiber",
    "done": false
  }
}
```

## 📝 สิ่งที่ควรสังเกต

### ใน starter/main.go:
- มีเพียง Hello endpoint และ TODO comments
- นักเรียนจะพิมพ์ตามผู้สอนเพื่อเพิ่ม:
  - Todo struct
  - Middleware
  - API endpoints

### ใน complete/main.go:
- มี comment ภาษาไทยอธิบายทุกส่วน
- มี error handling และ validation
- มีข้อมูลตัวอย่าง 3 รายการ
- ใช้ in-memory slice เก็บข้อมูล

## 🔍 Key Learning Points

1. **Fiber App Creation:**
   ```go
   app := fiber.New(fiber.Config{})
   ```

2. **Middleware Usage:**
   ```go
   app.Use(logger.New())
   app.Use(recover.New())
   ```

3. **Route Handling:**
   ```go
   app.Get("/todos", getTodos)
   app.Post("/todos", createTodo)
   ```

4. **JSON Response:**
   ```go
   return c.JSON(fiber.Map{
       "success": true,
       "data": todos,
   })
   ```

## ⏭️ ถัดไปคือ...
หลังจากเสร็จระดับ Beginner แล้ว ให้ไปต่อที่:
```bash
cd ../intermediate
cat README.md
```

---
**เวลาเรียน:** 30 นาที | **ความยาก:** ⭐⭐☆☆☆ 