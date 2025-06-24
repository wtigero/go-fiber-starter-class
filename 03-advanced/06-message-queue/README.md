# 📬 Message Queue & Events - ส่งข้อความแบบ Async (75 นาที)

## 📚 จุดประสงค์
เรียนรู้การส่งข้อความระหว่าง services และทำงาน Background jobs แบบ asynchronous

## 🎯 สิ่งที่จะเรียนรู้
- เข้าใจ Event-driven architecture
- ใช้ Channel ของ Go เป็น Message Queue แบบง่าย
- ส่ง Email notifications แบบ async
- ประมวลผล Background tasks
- Worker pattern สำหรับจัดการ jobs

## 🔄 Event Flow ที่จะสร้าง

```
User Action --> API --> Queue --> Worker --> External Service
     |                              |
     v                              v
 Immediate Response            Email/Notification
```

### ตัวอย่าง Events:
1. **UserRegistered** → ส่ง Welcome Email
2. **TodoCompleted** → ส่ง Achievement Notification  
3. **TodoOverdue** → ส่ง Reminder Email

## 📋 API Endpoints

### Main APIs
- `POST /users` - สร้าง user (trigger UserRegistered event)
- `PUT /todos/:id/complete` - ทำ todo เสร็จ (trigger TodoCompleted event)
- `GET /todos/overdue` - หา todos ที่เลยกำหนด (trigger TodoOverdue events)

### Monitoring APIs
- `GET /queue/status` - ดูสถานะ queue
- `GET /jobs/history` - ดูประวัติ jobs
- `GET /workers` - ดูสถานะ workers

## 🏃‍♂️ วิธีรัน

### รัน Application พร้อม Workers:
```bash
cd starter
go mod tidy
go run main.go
```

Application จะเริ่ม:
- API Server (port 3000)
- 3 Email Workers (background)
- 2 Notification Workers (background)

## 🧪 ทดสอบ Message Queue

### 1. สร้าง User (trigger Welcome Email)
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "นายทดสอบ",
    "email": "test@example.com"
  }'
```

**Response:**
```json
{
  "success": true,
  "user": {"id": 1, "name": "นายทดสอบ"},
  "message": "User created, welcome email queued"
}
```

### 2. ทำ Todo เสร็จ (trigger Achievement)
```bash
curl -X PUT http://localhost:3000/todos/1/complete
```

### 3. ดูสถานะ Queue
```bash
curl http://localhost:3000/queue/status
```
**Response:**
```json
{
  "email_queue": {
    "pending": 2,
    "processing": 1,
    "completed": 15,
    "failed": 0
  },
  "notification_queue": {
    "pending": 0,
    "processing": 0,
    "completed": 8,
    "failed": 1
  },
  "workers": {
    "email_workers": 3,
    "notification_workers": 2,
    "active": 1
  }
}
```

### 4. สร้าง Load เพื่อดู Workers ทำงาน
```bash
# สร้าง users เยอะ ๆ
for i in {1..20}; do
  curl -X POST http://localhost:3000/users \
    -H "Content-Type: application/json" \
    -d '{"name":"User '${i}'","email":"user'${i}'@example.com"}'
  echo
done
```

## 🔍 สิ่งสำคัญที่เรียนรู้

### 1. Event System
```go
type Event struct {
    Type      string      `json:"type"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
    ID        string      `json:"id"`
}

type EventBus struct {
    channels map[string]chan Event
    mutex    sync.RWMutex
}

func (eb *EventBus) Publish(eventType string, data interface{}) {
    event := Event{
        Type:      eventType,
        Data:      data,
        Timestamp: time.Now(),
        ID:        generateID(),
    }
    
    if ch, exists := eb.channels[eventType]; exists {
        select {
        case ch <- event:
            // Event sent successfully
        default:
            // Channel full, handle overflow
            log.Printf("Queue full for event type: %s", eventType)
        }
    }
}
```

### 2. Worker Pattern
```go
type Worker struct {
    ID       int
    JobQueue chan Job
    Quit     chan bool
}

func (w *Worker) Start() {
    go func() {
        for {
            select {
            case job := <-w.JobQueue:
                w.processJob(job)
            case <-w.Quit:
                log.Printf("Worker %d stopping", w.ID)
                return
            }
        }
    }()
}

func (w *Worker) processJob(job Job) {
    log.Printf("Worker %d processing job: %s", w.ID, job.Type)
    
    switch job.Type {
    case "send_email":
        w.sendEmail(job.Data)
    case "send_notification":
        w.sendNotification(job.Data)
    }
}
```

### 3. Email Service (Mock)
```go
type EmailService struct {
    // จำลองการส่ง email (ไม่ส่งจริง)
}

func (es *EmailService) SendWelcomeEmail(user User) error {
    // จำลองเวลาส่ง email
    time.Sleep(2 * time.Second)
    
    log.Printf("📧 Welcome email sent to %s (%s)", user.Name, user.Email)
    return nil
}

func (es *EmailService) SendReminderEmail(user User, todo Todo) error {
    time.Sleep(1 * time.Second)
    
    log.Printf("⏰ Reminder email sent to %s for todo: %s", user.Email, todo.Title)
    return nil
}
```

### 4. Job Queue Management
```go
type JobQueue struct {
    queue    chan Job
    workers  []*Worker
    quit     chan bool
    stats    JobStats
}

func (jq *JobQueue) AddJob(jobType string, data interface{}) {
    job := Job{
        Type:      jobType,
        Data:      data,
        CreatedAt: time.Now(),
    }
    
    select {
    case jq.queue <- job:
        jq.stats.Queued++
    default:
        jq.stats.Dropped++
        log.Printf("Job queue full, dropping job: %s", jobType)
    }
}
```

## 📊 Real-time Monitoring
```go
// WebSocket สำหรับดู queue status แบบ real-time
func setupWebSocket(app *fiber.App, jobQueue *JobQueue) {
    app.Get("/ws/queue", websocket.New(func(c *websocket.Conn) {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                status := jobQueue.GetStatus()
                c.WriteJSON(status)
            case <-c.Context().Done():
                return
            }
        }
    }))
}
```

## 📝 ใน starter/ จะมี:
- [ ] TODO: สร้าง Event Bus system
- [ ] TODO: สร้าง Worker pool
- [ ] TODO: เพิ่ม Email service (mock)
- [ ] TODO: เพิ่ม Job monitoring
- [ ] TODO: Handle graceful shutdown

## ✅ ใน complete/ จะมี:
- ✅ Event-driven architecture สมบูรณ์
- ✅ Worker pool ที่ scale ได้
- ✅ Job retry mechanism
- ✅ Monitoring และ stats
- ✅ Graceful shutdown

## 💡 ประโยชน์ในการใช้งานจริง

### Performance:
- API response เร็วขึ้น (ไม่รอ background tasks)
- Handle traffic spike ได้ดีขึ้น
- Resource utilization ดีขึ้น

### Reliability:
- Job retry หาก fail
- ไม่สูญหาย events หาก server restart
- Graceful degradation

### Scalability:
- เพิ่ม workers ตามความต้องการ
- แยก queue ตาม priority
- Scale ข้าม servers ได้

## ⏭️ ขั้นต่อไป
หากต้องการ Production-grade:
- ใช้ Redis หรือ RabbitMQ
- เพิ่ม Dead Letter Queue
- ทำ Distributed queue
- เพิ่ม Circuit breaker

---
**เวลาเรียน:** 75 นาที | **ความยาก:** ⭐⭐⭐⭐☆  
**เหมาะสำหรับ:** App ที่มี background processing 