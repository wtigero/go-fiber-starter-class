package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// TODO: สร้าง Job struct
// type Job struct {
//     ID       string    `json:"id"`
//     Type     string    `json:"type"`
//     Payload  string    `json:"payload"`
//     Status   string    `json:"status"` // pending, processing, completed, failed
//     CreatedAt time.Time `json:"created_at"`
//     CompletedAt *time.Time `json:"completed_at,omitempty"`
// }

// TODO: สร้าง Worker struct
// type Worker struct {
//     ID       int
//     JobQueue chan Job
//     Quit     chan bool
// }

// TODO: สร้าง Queue struct
// type Queue struct {
//     jobs    []Job
//     workers []Worker
//     mutex   sync.RWMutex
// }

// ข้อมูลตัวอย่าง
type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Notification struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		},
	})

	app.Use(logger.New())

	// TODO: สร้าง queue และ workers
	// queue := NewQueue(3) // 3 workers
	// queue.Start()

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Message Queue API - Ready to process jobs!",
			"status":  "ok",
		})
	})

	// TODO: สร้าง job submission endpoints
	// app.Post("/jobs/email", createEmailJobHandler(queue))
	// app.Post("/jobs/notification", createNotificationJobHandler(queue))
	// app.Get("/jobs", getJobsHandler(queue))
	// app.Get("/jobs/:id", getJobHandler(queue))

	// Sample endpoints (without queue)
	app.Post("/send-email", sendEmailHandler)
	app.Post("/send-notification", sendNotificationHandler)

	log.Println("🚀 Message Queue API started on port 3000")
	log.Println("📬 Ready to process background jobs")
	log.Fatal(app.Listen(":3000"))
}

// Sample handlers (synchronous)
func sendEmailHandler(c *fiber.Ctx) error {
	var email Email
	if err := c.BodyParser(&email); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if email.To == "" || email.Subject == "" {
		return c.Status(400).JSON(fiber.Map{"error": "To and subject are required"})
	}

	// จำลองการส่ง email (ใช้เวลานาน)
	time.Sleep(2 * time.Second)

	log.Printf("📧 Email sent to %s: %s", email.To, email.Subject)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email sent successfully",
	})
}

func sendNotificationHandler(c *fiber.Ctx) error {
	var notification Notification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if notification.UserID == 0 || notification.Message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "UserID and message are required"})
	}

	// จำลองการส่ง notification (ใช้เวลาปานกลาง)
	time.Sleep(500 * time.Millisecond)

	log.Printf("🔔 Notification sent to user %d: %s", notification.UserID, notification.Message)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification sent successfully",
	})
}

// TODO: สร้าง NewQueue function
// func NewQueue(numWorkers int) *Queue {
//     queue := &Queue{
//         jobs:    make([]Job, 0),
//         workers: make([]Worker, numWorkers),
//     }
//
//     // สร้าง workers
//     for i := 0; i < numWorkers; i++ {
//         worker := Worker{
//             ID:       i,
//             JobQueue: make(chan Job),
//             Quit:     make(chan bool),
//         }
//         queue.workers[i] = worker
//     }
//
//     return queue
// }

// TODO: สร้าง Start method
// func (q *Queue) Start() {
//     for i := range q.workers {
//         go q.startWorker(&q.workers[i])
//     }
// }

// TODO: สร้าง startWorker method
// func (q *Queue) startWorker(worker *Worker) {
//     log.Printf("👷 Worker %d started", worker.ID)
//
//     for {
//         select {
//         case job := <-worker.JobQueue:
//             q.processJob(worker, job)
//         case <-worker.Quit:
//             log.Printf("👷 Worker %d stopped", worker.ID)
//             return
//         }
//     }
// }

// TODO: สร้าง processJob method
// func (q *Queue) processJob(worker *Worker, job Job) {
//     log.Printf("👷 Worker %d processing job %s", worker.ID, job.ID)
//
//     // อัปเดตสถานะ
//     q.updateJobStatus(job.ID, "processing")
//
//     // ประมวลผลตาม job type
//     var err error
//     switch job.Type {
//     case "email":
//         err = q.processEmailJob(job)
//     case "notification":
//         err = q.processNotificationJob(job)
//     default:
//         err = fmt.Errorf("unknown job type: %s", job.Type)
//     }
//
//     // อัปเดตสถานะตามผลลัพธ์
//     if err != nil {
//         log.Printf("❌ Job %s failed: %v", job.ID, err)
//         q.updateJobStatus(job.ID, "failed")
//     } else {
//         log.Printf("✅ Job %s completed", job.ID)
//         q.updateJobStatus(job.ID, "completed")
//     }
// }

// TODO: สร้าง processEmailJob method
// func (q *Queue) processEmailJob(job Job) error {
//     // จำลองการส่ง email
//     time.Sleep(2 * time.Second)
//
//     log.Printf("📧 Email job processed: %s", job.Payload)
//     return nil
// }

// TODO: สร้าง processNotificationJob method
// func (q *Queue) processNotificationJob(job Job) error {
//     // จำลองการส่ง notification
//     time.Sleep(500 * time.Millisecond)
//
//     log.Printf("🔔 Notification job processed: %s", job.Payload)
//     return nil
// }

// TODO: สร้าง AddJob method
// func (q *Queue) AddJob(job Job) {
//     q.mutex.Lock()
//     defer q.mutex.Unlock()
//
//     job.Status = "pending"
//     job.CreatedAt = time.Now()
//     q.jobs = append(q.jobs, job)
//
//     // ส่ง job ไปให้ worker ที่ว่าง
//     for i := range q.workers {
//         select {
//         case q.workers[i].JobQueue <- job:
//             return
//         default:
//             continue
//         }
//     }
//
//     log.Printf("⏳ Job %s queued (all workers busy)", job.ID)
// }

// TODO: สร้าง handlers
// func createEmailJobHandler(queue *Queue) fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         var email Email
//         if err := c.BodyParser(&email); err != nil {
//             return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
//         }
//
//         job := Job{
//             ID:      generateJobID(),
//             Type:    "email",
//             Payload: fmt.Sprintf("To: %s, Subject: %s", email.To, email.Subject),
//         }
//
//         queue.AddJob(job)
//
//         return c.Status(202).JSON(fiber.Map{
//             "success": true,
//             "job_id":  job.ID,
//             "message": "Email job queued",
//         })
//     }
// }

// TODO: สร้าง utility functions
// func generateJobID() string {
//     return fmt.Sprintf("job_%d", time.Now().UnixNano())
// }
