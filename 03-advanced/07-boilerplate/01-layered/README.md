# 🏢 Layered Architecture

> **Traditional N-Tier Architecture** - การแบ่งชั้นตาม responsibility

## 🎯 หลักการ

Layered Architecture แบ่งระบบเป็นชั้น (layers) โดยแต่ละชั้นมี responsibility ชัดเจน และสื่อสารกับชั้นที่อยู่ติดกันเท่านั้น

## 📋 โครงสร้าง

```
01-layered/
├── controllers/      # Presentation Layer (HTTP handlers)
├── services/         # Business Layer (business logic)
├── repositories/     # Data Access Layer (database operations)
├── models/          # Data models/entities
├── database/        # Database connection
├── main.go          # Application entry point
└── README.md        # ไฟล์นี้
```

## 🔄 Flow การทำงาน

```
HTTP Request → Controller → Service → Repository → Database
                    ↓
HTTP Response ← Controller ← Service ← Repository ← Database
```

## 📊 Layers ละเอียด

### **1. 🎭 Presentation Layer (Controllers)**
- รับ HTTP requests
- Validate input data  
- เรียก Business Layer
- ส่ง HTTP responses

### **2. 💼 Business Layer (Services)**
- Business logic และ rules
- Data validation และ processing
- เรียก Data Access Layer
- ไม่รู้จัก HTTP หรือ Database details

### **3. 🗄️ Data Access Layer (Repositories)**
- Database operations (CRUD)
- SQL queries
- Data mapping
- ไม่รู้จัก business rules

### **4. 📦 Models**
- Data structures
- Entity definitions
- Shared across layers

## ✅ ข้อดี

- **เข้าใจง่าย**: โครงสร้างชัดเจน straightforward
- **แยกหน้าที่**: แต่ละ layer มี responsibility ชัดเจน
- **Reusable**: Service layer ใช้ได้หลาย controllers
- **เริ่มต้นเร็ว**: setup ง่าย เหมาะ small projects

## ❌ ข้อเสีย

- **Tight Coupling**: layers depend on กันแบบ chain
- **ทดสอบยาก**: ต้อง mock หลาย layers
- **Database Driven**: มักจะ design ตาม database structure
- **Business Logic กระจาย**: บางครั้งอยู่ใน controller หรือ repository

## 🚀 การรัน

```bash
cd 03-advanced/07-boilerplate/01-layered
go run main.go
```

## 🔧 API Endpoints

```bash
# สร้าง user
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# ดู users ทั้งหมด
curl http://localhost:3000/users

# ดู user ตาม ID
curl http://localhost:3000/users/1

# อัปเดต user
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated"}'

# ลบ user
curl -X DELETE http://localhost:3000/users/1
```

## 🎓 เมื่อไหร่ควรใช้

✅ **เหมาะกับ:**
- Small to medium projects
- Team ที่เริ่มต้น microservices
- Prototype หรือ MVP
- เวลาต้องการ delivery เร็ว

❌ **ไม่เหมาะกับ:**
- Large enterprise applications
- Complex business rules
- ต้องการ high testability
- Multiple data sources

---

💡 **คำแนะนำ**: Layered เป็น starting point ที่ดี แต่เมื่อ project โต ควรพิจารณา architecture อื่นที่ maintainable กว่า 