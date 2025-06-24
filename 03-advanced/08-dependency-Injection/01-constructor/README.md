# 🔨 Constructor Injection (Manual DI)

## 🎯 วัตถุประสงค์
เรียนรู้ **Constructor Injection** แบบ Manual ที่ง่ายที่สุดและนิยมใช้มากที่สุด

## 💡 หลักการ Constructor Injection

```go
// ✅ GOOD: Dependency ส่งผ่าน Constructor
type UserService struct {
    repo UserRepository  // รับจากภายนอก
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}  // Inject ตอนสร้าง
}

// ❌ BAD: สร้าง Dependency ข้างใน
func NewUserService() *UserService {
    repo := NewUserRepository()  // Hard-coded dependency
    return &UserService{repo: repo}
}
```

## 🏗️ โครงสร้าง

```
01-constructor/
├── models/
│   └── user.go           # Entity
├── repository/
│   └── user_repository.go # Data layer
├── service/
│   └── user_service.go    # Business logic
├── controller/
│   └── user_controller.go # HTTP handlers
├── main.go               # 🔌 Manual DI wiring
└── go.mod
```

## 🔄 Data Flow

```
main.go → สร้าง repo → สร้าง service → สร้าง controller → รัน server
         ↓
    Constructor Injection Chain:
    Repository ← Service ← Controller
```

## 🔌 DI Wiring ใน main.go

```go
func main() {
    // === MANUAL DEPENDENCY INJECTION ===
    
    // 1. สร้าง Repository (ชั้นล่างสุด)
    userRepo := repository.NewUserRepository()
    
    // 2. สร้าง Service (ฉีด repo เข้าไป)
    userService := service.NewUserService(userRepo)
    
    // 3. สร้าง Controller (ฉีด service เข้าไป)  
    userController := controller.NewUserController(userService)
    
    // 4. Setup routes และรัน server
    app := fiber.New()
    setupRoutes(app, userController)
    app.Listen(":3000")
}
```

## ✅ ข้อดี Constructor Injection

- **🎯 Simple**: เข้าใจง่ายที่สุด ไม่ต้องเรียนรู้ library เพิ่ม
- **⚡ Fast**: ไม่มี overhead จาก reflection หรือ code generation
- **🛡️ Safe**: Compile-time safety ตรวจสอบได้ตอน compile
- **🧪 Testable**: Mock dependencies ได้ง่าย
- **📦 Zero Dependency**: ไม่ต้องติดตั้ง external packages

## ❌ ข้อเสีย

- **📝 Verbose**: เขียนโค้ดเยอะ เมื่อมี dependencies เยอะ
- **🔄 Manual Order**: ต้องจำลำดับการสร้าง dependencies
- **🔧 Maintenance**: เปลี่ยนแปลงยาก เมื่อโปรเจคใหญ่

## 🧪 การทดสอบ

```bash
# รันแอป
go run main.go

# ทดสอบ API
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'

curl http://localhost:3000/users
```

## 📊 API Endpoints

| Method | Endpoint     | Description |
|--------|-------------|-------------|
| GET    | /users      | ดึงผู้ใช้ทั้งหมด |
| POST   | /users      | สร้างผู้ใช้ใหม่ |
| GET    | /users/:id  | ดึงผู้ใช้ตาม ID |
| PUT    | /users/:id  | แก้ไขผู้ใช้ |
| DELETE | /users/:id  | ลบผู้ใช้ |

## 🎯 เมื่อไหร่ควรใช้?

✅ **ใช้เมื่อ:**
- โปรเจคเล็ก-กลาง (< 10 services)
- ทีมใหม่ๆ ที่เพิ่งเรียน Go
- ต้องการความเรียบง่าย
- ไม่มี complex dependencies

❌ **ไม่ควรใช้เมื่อ:**
- โปรเจคใหญ่มาก (> 50 services)
- มี circular dependencies
- ต้องการ advanced features (singleton, scoped lifecycle) 