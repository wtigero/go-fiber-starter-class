# 🧅 Onion Architecture (Onion Pattern)

## 🎯 วัตถุประสงค์
เรียนรู้การออกแบบ **Onion Architecture** ที่เน้นการแยกชั้นเป็นวงกลม โดยชั้นในไม่รู้จักชั้นนอก

## 🏗️ โครงสร้าง Onion Architecture

```
🧅 ONION LAYERS (จากในออกนอก)
┌─────────────────────────────────────┐
│  Infrastructure (DB, HTTP, Cache)   │ ← Outermost
├─────────────────────────────────────┤
│     Application (Use Cases)         │
├─────────────────────────────────────┤ 
│        Domain (Entities)            │ ← Core/Innermost
└─────────────────────────────────────┘
```

### 📁 Folder Structure
```
04-onion/
├── domain/                    # 🔥 CORE BUSINESS
│   ├── entities/             # ข้อมูลหลัก
│   └── value-objects/        # กฎการตรวจสอบ
├── application/              # 🧠 USE CASES  
│   └── user_service.go       # บิซิเนสโลจิก
├── infrastructure/           # 🔧 EXTERNAL
│   └── memory_repository.go  # การเก็บข้อมูล
├── presentation/             # 🌐 ADAPTERS
│   └── http_controller.go    # HTTP APIs
└── main.go                   # 🚀 Dependency Injection
```

## 🔄 Data Flow (Onion Pattern)

```
HTTP Request → Presentation → Application → Domain
                    ↑              ↓          ↓
Infrastructure ←────┴──────────────┴──────────┘
```

### 📜 Key Principles

1. **🎯 Domain Independence**: Core business logic ไม่พึ่งพาอะไรเลย
2. **🔄 Dependency Inversion**: ชั้นนอกพึ่งพาชั้นใน (ไม่ใช่กลับกัน)
3. **🧩 Isolation**: แต่ละชั้นมีหน้าที่ชัดเจน
4. **🛡️ Protection**: Business rules ถูกปกป้องในชั้นใน

## 🧪 การทดสอบ

```bash
# ติดตั้ง dependencies
go mod tidy

# รันแอป
go run main.go

# ทดสอบ API
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

## ✅ ข้อดี Onion Architecture

- **🛡️ Business Protection**: Core business ปลอดภัยจากการเปลี่ยนแปลง
- **🧪 Easy Testing**: Mock dependencies ง่าย
- **🔄 Flexible**: เปลี่ยน Infrastructure ได้ไม่กระทบ Core
- **📈 Scalable**: เพิ่มฟีเจอร์ใหม่ได้ง่าย

## ❌ ข้อเสีย

- **🔧 Complex Setup**: ต้องสร้างหลายชั้น
- **📝 More Code**: มี Interface เยอะ
- **🤔 Learning Curve**: ต้องเข้าใจหลักการ DI ก่อน

## 🎯 เมื่อไหร่ควรใช้?

✅ **ใช้เมื่อ:**
- โปรเจคขนาดใหญ่ที่ซับซ้อน
- ต้องการ Business Logic ที่แน่นอน
- ทีมงานใหญ่ ต้องการแยกงานชัดเจน
- ต้องการทดสอบครอบคลุม

❌ **ไม่ควรใช้เมื่อ:**
- โปรเจคเล็กๆ แบบ CRUD
- ต้องการพัฒนาเร็ว (MVP)
- ทีมงานเล็ก ไม่คุ้นชินกับ Clean Architecture

## 🔗 Pattern เปรียบเทียบ

| Pattern | Complexity | Learning | When to Use |
|---------|------------|----------|-------------|
| Layered | ⭐⭐ | Easy | Simple CRUD |
| Clean | ⭐⭐⭐ | Medium | Medium Projects |
| Hexagonal | ⭐⭐⭐⭐ | Hard | Complex APIs |
| **Onion** | ⭐⭐⭐⭐⭐ | **Expert** | **Enterprise** | 