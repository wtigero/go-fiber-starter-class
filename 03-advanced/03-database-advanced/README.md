# 🗄️ Database Advanced - การจัดการฐานข้อมูลขั้นสูง (45 นาที)

## 📚 จุดประสงค์
เรียนรู้การใช้ฐานข้อมูลอย่างมืออาชีพ พร้อม pattern ที่ใช้ในการทำงานจริง

## 🎯 สิ่งที่จะเรียนรู้
- Repository Pattern แยก business logic กับ database
- Transaction Management จัดการข้อมูลที่เกี่ยวข้องกัน
- Database Migration ระบบอัปเดตโครงสร้าง
- Connection Pool จัดการ connection อย่างมีประสิทธิภาพ

## 📊 ตัวอย่างที่ใช้: ระบบ E-commerce เล็ก ๆ
```go
type User struct {
    ID    int    `json:"id" db:"id"`
    Email string `json:"email" db:"email"`
    Name  string `json:"name" db:"name"`
}

type Product struct {
    ID    int     `json:"id" db:"id"`
    Name  string  `json:"name" db:"name"`
    Price float64 `json:"price" db:"price"`
    Stock int     `json:"stock" db:"stock"`
}

type Order struct {
    ID       int     `json:"id" db:"id"`
    UserID   int     `json:"user_id" db:"user_id"`
    Total    float64 `json:"total" db:"total"`
    Status   string  `json:"status" db:"status"`
}
```

## 📋 API Endpoints
- `GET /users` - ดูรายการ users
- `POST /orders` - สร้างออเดอร์ (ใช้ transaction)
- `GET /orders/:id` - ดูรายละเอียดออเดอร์
- `POST /migrate` - รัน database migration

## 🏗️ Architecture Pattern

### Repository Pattern
```
Controller -> Service -> Repository -> Database
```

### ไฟล์ที่จะมี:
```
/models/     - structs และ interfaces
/repository/ - database operations
/service/    - business logic
/migration/  - database schema changes
```

## 🏃‍♂️ วิธีรัน

### เตรียม Database (PostgreSQL)
```bash
# ใช้ Docker
docker run --name postgres-advanced \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 -d postgres:15

# สร้าง database
docker exec postgres-advanced createdb -U postgres ecommerce
```

### รัน Application
```bash
cd starter
go mod tidy
go run main.go
```

## 🧪 ทดสอบ Transaction

### 1. สร้างออเดอร์ (สำเร็จ)
```bash
curl -X POST http://localhost:3000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"product_id": 1, "quantity": 2}
    ]
  }'
```

### 2. สร้างออเดอร์ (สินค้าไม่พอ)
```bash
curl -X POST http://localhost:3000/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,  
    "items": [
      {"product_id": 1, "quantity": 999}
    ]
  }'
```
**ผลลัพธ์:** Transaction rollback, stock ไม่เปลี่ยน

## 🔍 สิ่งสำคัญที่เรียนรู้

### 1. Repository Interface
```go
type UserRepository interface {
    GetAll() ([]User, error)
    GetByID(id int) (*User, error)
    Create(user *User) error
}

type userRepository struct {
    db *sql.DB
}
```

### 2. Transaction ใน Service
```go
func (s *OrderService) CreateOrder(order *CreateOrderRequest) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback() // rollback หากมี error
    
    // 1. ตรวจสอบ stock
    // 2. ลด stock  
    // 3. สร้าง order
    // 4. สร้าง order items
    
    return tx.Commit() // commit หากทุกอย่างสำเร็จ
}
```

### 3. Migration System
```go
type Migration struct {
    Version string
    Up      string
    Down    string
}

var migrations = []Migration{
    {
        Version: "001_create_users",
        Up: `CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE,
            name VARCHAR(255)
        )`,
        Down: `DROP TABLE users`,
    },
}
```

## 📝 ใน starter/ จะมี:
- [ ] TODO: สร้าง Repository interfaces
- [ ] TODO: สร้าง Service layer
- [ ] TODO: ใช้ Transaction ใน CreateOrder
- [ ] TODO: สร้าง Migration system

## ✅ ใน complete/ จะมี:
- ✅ Repository pattern สมบูรณ์
- ✅ Service layer พร้อม business logic
- ✅ Transaction management
- ✅ Migration system
- ✅ Connection pooling
- ✅ Error handling ครบถ้วน

## 💡 ประโยชน์ของ Pattern นี้

### Repository Pattern:
- แยก database logic ออกจาก business logic
- ง่ายต่อการเปลี่ยน database
- ทำ unit testing ได้ง่ายขึ้น

### Transaction:
- รับประกันความสมบูรณ์ของข้อมูล
- ป้องกัน race condition
- ยกเลิกการเปลี่ยนแปลงได้หาก error

### Migration:
- อัปเดต database schema อย่างปลอดภัย
- ย้อนกลับได้หาก error
- ทำงานร่วมกับทีมได้ง่าย

## ⏭️ ขั้นต่อไป
- Database Indexing สำหรับ performance
- Query optimization
- Database monitoring

---
**เวลาเรียน:** 45 นาที | **ความยาก:** ⭐⭐⭐⭐☆  
**เหมาะสำหรับ:** คนที่จะทำ production app 