# 🚀 Advanced Level - Go Fiber Production Ready (60+ นาที)

## 📚 จุดประสงค์
เรียนรู้การสร้าง Production-ready API ด้วย Go Fiber พร้อมระบบที่ซับซ้อนขึ้น

## 🎯 Advanced Topics (เลือกเรียนตามความสนใจ)

### 1. 🔐 JWT Authentication & Authorization (60 นาที)
**ระบบ Login/Register พร้อม Role-based access control**
- JWT Token generation & validation
- User registration และ login
- Role-based middleware (Admin, User)
- Password hashing ด้วย bcrypt
- Refresh token mechanism

### 2. 🚦 Rate Limiting & Caching (45 นาที)  
**การควบคุม API calls และ caching ด้วย Redis**
- Rate limiting ต่อ IP และ User
- Redis caching สำหรับ API responses
- Cache invalidation strategies
- Memory vs Redis performance

### 3. 🏗️ Microservices Architecture (90 นาที)
**แยก API เป็น services หลาย ๆ ตัว**
- User Service + Todo Service แยกกัน
- Service-to-service communication
- API Gateway pattern
- Service discovery

### 4. 📬 Message Queue & Events (75 นาที)
**Async processing ด้วย Message Queue**
- RabbitMQ หรือ NATS integration
- Event-driven architecture
- Background job processing
- Email notifications

### 5. 📊 Monitoring & Metrics (60 นาที)
**การติดตาม application performance**
- Prometheus metrics
- Health checks advanced
- Logging ด้วย structured logs
- Performance monitoring

### 6. 🗄️ Database Migration & Advanced Patterns (45 นาที)
**การจัดการ Database อย่างมืออาชีพ**
- Migration system
- Database seeding
- Repository pattern
- Transaction management

## 📁 โครงสร้างโฟลเดอร์

```
advanced/
├── 01-jwt-auth/           # JWT Authentication
│   ├── starter/
│   ├── complete/
│   └── README.md
├── 02-rate-limit-cache/   # Rate Limiting & Caching  
│   ├── starter/
│   ├── complete/
│   └── README.md
├── 03-microservices/      # Microservices
│   ├── starter/
│   ├── complete/
│   └── README.md
├── 04-message-queue/      # Message Queue
│   ├── starter/
│   ├── complete/
│   └── README.md
├── 05-monitoring/         # Monitoring
│   ├── starter/
│   ├── complete/
│   └── README.md
└── 06-database-advanced/  # Database Advanced
    ├── starter/
    ├── complete/
    └── README.md
```

## 🎯 แนะนำเส้นทางการเรียน

### สำหรับผู้ที่ต้องการทำ Production App:
1. **JWT Authentication** (ขาดไม่ได้)
2. **Rate Limiting & Caching** (สำคัญมาก)  
3. **Database Advanced** (ต้องมี)
4. **Monitoring** (ใช้จริง)

### สำหรับผู้ที่ต้องการเป็น Backend Architect:
1. **Microservices Architecture**
2. **Message Queue & Events**
3. **Monitoring & Metrics**
4. **Database Advanced**

## 🏃‍♂️ วิธีเริ่มต้น

เลือกหัวข้อที่สนใจ:

### ตัวอย่าง: เริ่มต้นด้วย JWT Authentication
```bash
cd 01-jwt-auth
cat README.md
```

### หรือดู Rate Limiting & Caching  
```bash
cd 02-rate-limit-cache
cat README.md
```

## 💡 ข้อแนะนำ

1. **เรียงลำดับตามความยาก**: JWT → Rate Limiting → Database → Monitoring → Microservices → Message Queue
2. **เลือกตามโปรเจกต์**: ถ้าทำ API เดียว ไม่จำเป็นต้อง Microservices
3. **ลองทีละหัวข้อ**: แต่ละหัวข้อค่อนข้างซับซ้อน ใช้เวลาเรียนแล้วลองทำตาม
4. **Production Ready**: ทุกตัวอย่างเน้นใช้ได้จริงในการทำงาน

## 🔍 เทคโนโลยีที่จะได้เรียนรู้

- **JWT** - JSON Web Tokens
- **Redis** - In-memory caching
- **RabbitMQ/NATS** - Message brokers  
- **Prometheus** - Metrics collection
- **Docker Compose** - Multi-service setup
- **PostgreSQL** - Advanced database features
- **gRPC** - Service communication
- **Grafana** - Metrics visualization

---

**คำแนะนำ**: เริ่มจากหัวข้อที่ 1 หรือเลือกตามความสนใจ แต่ละหัวข้อออกแบบให้เรียนรู้ได้อย่างอิสระ

**เวลาเรียนรวม:** 375+ นาที (6+ ชั่วโมง) | **ความยาก:** ⭐⭐⭐⭐⭐ 