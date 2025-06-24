# 🔌 Dependency Injection Patterns

## 🎯 วัตถุประสงค์
เรียนรู้ **3 แบบ Dependency Injection** ที่นิยมใช้ในโลก Go พร้อมตัวอย่างจริง

## 📁 โครงสร้าง

```
08-dependency-injection/
├── 01-constructor/         # 🔨 Manual Constructor (ง่ายที่สุด)
├── 02-uber-dig/           # 🏗️ Uber Dig (Container-based)
├── 03-google-wire/        # ⚡ Google Wire (Code Generation)
└── README.md              # 📖 เปรียบเทียบทั้ง 3 แบบ
```

## 🔄 ความแตกต่างของ DI Patterns

### 1. 🔨 **Constructor Injection** (Manual)
```go
// สร้าง dependencies เอง
userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo)
userController := controller.NewUserController(userService)
```

### 2. 🏗️ **Uber Dig** (Container-based)
```go
// ใช้ Container จัดการ
container := dig.New()
container.Provide(repository.NewUserRepository)
container.Provide(service.NewUserService)
container.Invoke(func(uc *controller.UserController) {
    // ใช้งาน
})
```

### 3. ⚡ **Google Wire** (Code Generation)
```go
//go:build wireinject
// +build wireinject

//go:generate wire
func InitializeUserController() *controller.UserController {
    wire.Build(
        repository.NewUserRepository,
        service.NewUserService,
        controller.NewUserController,
    )
    return nil
}
```

## 📊 เปรียบเทียบ 3 แบบ

| แบบ | ความง่าย | Performance | Setup | เมื่อไรใช้ |
|-----|----------|-------------|--------|----------|
| **Constructor** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | โปรเจคเล็ก-กลาง |
| **Uber Dig** | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ | โปรเจคกลาง-ใหญ่ |
| **Google Wire** | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | โปรเจคใหญ่ |

## ✅ ข้อดี - ข้อเสีย

### 🔨 Constructor Injection
✅ **ข้อดี:**
- ง่ายที่สุด ไม่ต้องเรียนรู้เพิ่ม
- ไม่มี external dependencies
- Compile-time safety
- เร็วที่สุด (ไม่มี overhead)

❌ **ข้อเสีย:**
- เขียนโค้ดเยอะ เมื่อมี dependencies เยอะ
- ต้องจำลำดับการสร้าง
- เปลี่ยนแปลงยาก เมื่อโปรเจคใหญ่

### 🏗️ Uber Dig
✅ **ข้อดี:**
- ไม่ต้องจำลำดับ การสร้าง
- Lifecycle management
- ทำงานแบบ runtime

❌ **ข้อเสีย:**
- เรียนรู้ยาก (Reflection-based)
- ช้ากว่า Constructor
- Runtime errors

### ⚡ Google Wire
✅ **ข้อดี:**
- Compile-time safety 
- เร็วเท่า Constructor
- Handle ลำดับให้อัตโนมัติ

❌ **ข้อเสีย:**
- ต้อง code generation
- Learning curve สูง
- Setup ซับซ้อน

## 🎯 แนะนำการใช้งาน

### 👶 **เริ่มต้น**: Constructor Injection
- โปรเจคเล็ก (< 10 services)
- ทีมใหม่ๆ ที่เพิ่งเรียน Go
- ต้องการ simplicity

### 🚀 **ขั้นกลาง**: Uber Dig  
- โปรเจคกลาง (10-50 services)
- ต้องการ flexibility
- มี complex dependencies

### 💪 **ขั้นสูง**: Google Wire
- โปรเจคใหญ่ (50+ services)
- ต้องการ performance สูงสุด
- ทีม experienced

## 🧪 การทดสอบ

```bash
# ทดสอบแต่ละแบบ
cd 01-constructor && go run main.go
cd 02-uber-dig && go run main.go  
cd 03-google-wire && go run main.go

# ทดสอบ API
curl http://localhost:3000/users
```

ทุกตัวอย่างใช้ **API เดียวกัน** แต่ใช้ DI ต่างกัน! 🎯 