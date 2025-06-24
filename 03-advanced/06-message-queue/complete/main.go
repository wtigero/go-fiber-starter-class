package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Job structures
type Job struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Payload     string     `json:"payload"`
	Status      string     `json:"status"` // pending, processing, completed, failed
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Error       string     `json:"error,omitempty"`
}

type Worker struct {
	ID       int
	JobQueue chan Job
	Quit     chan bool
}

type Queue struct {
	jobs       []Job
	workers    []Worker
	jobChannel chan Job
	mutex      sync.RWMutex
}

// Request structures
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

// Global queue instance
var queue *Queue

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

	// Initialize queue with 3 workers
	queue = NewQueue(3)
	queue.Start()

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Message Queue API - Complete Version",
			"status":  "ok",
			"workers": len(queue.workers),
		})
	})

	// Job submission endpoints
	app.Post("/jobs/email", createEmailJobHandler)
	app.Post("/jobs/notification", createNotificationJobHandler)
	
	// Job monitoring endpoints
	app.Get("/jobs", getJobsHandler)
	app.Get("/jobs/:id", getJobHandler)
	app.Get("/dashboard", dashboardHandler)

	// Legacy synchronous endpoints for comparison
	app.Post("/send-email-sync", sendEmailSyncHandler)
	app.Post("/send-notification-sync", sendNotificationSyncHandler)

	log.Println("🚀 Message Queue API started on port 3000")
	log.Println("📬 Workers ready to process background jobs")
	log.Println("📊 Visit http://localhost:3000/dashboard for monitoring")
	log.Fatal(app.Listen(":3000"))
}

// Queue implementation
func NewQueue(numWorkers int) *Queue {
	queue := &Queue{
		jobs:       make([]Job, 0),
		workers:    make([]Worker, numWorkers),
		jobChannel: make(chan Job, 100), // Buffer for 100 jobs
	}

	// Create workers
	for i := 0; i < numWorkers; i++ {
		worker := Worker{
			ID:       i,
			JobQueue: make(chan Job),
			Quit:     make(chan bool),
		}
		queue.workers[i] = worker
	}

	return queue
}

func (q *Queue) Start() {
	// Start job dispatcher
	go q.dispatcher()

	// Start workers
	for i := range q.workers {
		go q.startWorker(&q.workers[i])
	}

	log.Printf("✅ Started %d workers", len(q.workers))
}

func (q *Queue) dispatcher() {
	for {
		select {
		case job := <-q.jobChannel:
			// Find available worker
			for i := range q.workers {
				select {
				case q.workers[i].JobQueue <- job:
					goto nextJob
				default:
					continue
				}
			}
			// If all workers busy, put back to channel
			go func() {
				time.Sleep(100 * time.Millisecond)
				q.jobChannel <- job
			}()
		nextJob:
		}
	}
}

func (q *Queue) startWorker(worker *Worker) {
	log.Printf("👷 Worker %d started", worker.ID)

	for {
		select {
		case job := <-worker.JobQueue:
			q.processJob(worker, job)
		case <-worker.Quit:
			log.Printf("👷 Worker %d stopped", worker.ID)
			return
		}
	}
}

func (q *Queue) processJob(worker *Worker, job Job) {
	log.Printf("👷 Worker %d processing job %s (%s)", worker.ID, job.ID, job.Type)

	// Update job status
	now := time.Now()
	q.updateJobStatus(job.ID, "processing", &now, nil, "")

	// Process based on job type
	var err error
	switch job.Type {
	case "email":
		err = q.processEmailJob(job)
	case "notification":
		err = q.processNotificationJob(job)
	default:
		err = fmt.Errorf("unknown job type: %s", job.Type)
	}

	// Update final status
	completedAt := time.Now()
	if err != nil {
		log.Printf("❌ Job %s failed: %v", job.ID, err)
		q.updateJobStatus(job.ID, "failed", nil, &completedAt, err.Error())
	} else {
		log.Printf("✅ Job %s completed in %v", job.ID, completedAt.Sub(now))
		q.updateJobStatus(job.ID, "completed", nil, &completedAt, "")
	}
}

func (q *Queue) processEmailJob(job Job) error {
	// Parse email data
	var email Email
	if err := json.Unmarshal([]byte(job.Payload), &email); err != nil {
		return fmt.Errorf("invalid email payload: %w", err)
	}

	// Simulate email sending (longer process)
	time.Sleep(2 * time.Second)

	log.Printf("📧 Email sent to %s: %s", email.To, email.Subject)
	return nil
}

func (q *Queue) processNotificationJob(job Job) error {
	// Parse notification data
	var notification Notification
	if err := json.Unmarshal([]byte(job.Payload), &notification); err != nil {
		return fmt.Errorf("invalid notification payload: %w", err)
	}

	// Simulate notification sending
	time.Sleep(500 * time.Millisecond)

	log.Printf("🔔 Notification sent to user %d: %s", notification.UserID, notification.Message)
	return nil
}

func (q *Queue) AddJob(job Job) {
	q.mutex.Lock()
	job.Status = "pending"
	job.CreatedAt = time.Now()
	q.jobs = append(q.jobs, job)
	q.mutex.Unlock()

	// Send to job channel
	select {
	case q.jobChannel <- job:
		log.Printf("📝 Job %s queued successfully", job.ID)
	default:
		log.Printf("⚠️ Job queue full, job %s delayed", job.ID)
		go func() {
			time.Sleep(1 * time.Second)
			q.jobChannel <- job
		}()
	}
}

func (q *Queue) updateJobStatus(jobID, status string, startedAt, completedAt *time.Time, errorMsg string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i := range q.jobs {
		if q.jobs[i].ID == jobID {
			q.jobs[i].Status = status
			if startedAt != nil {
				q.jobs[i].StartedAt = startedAt
			}
			if completedAt != nil {
				q.jobs[i].CompletedAt = completedAt
			}
			if errorMsg != "" {
				q.jobs[i].Error = errorMsg
			}
			break
		}
	}
}

func (q *Queue) GetJobs() []Job {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	// Return copy of jobs slice
	jobs := make([]Job, len(q.jobs))
	copy(jobs, q.jobs)
	return jobs
}

func (q *Queue) GetJobByID(id string) *Job {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	for _, job := range q.jobs {
		if job.ID == id {
			return &job
		}
	}
	return nil
}

// Handlers
func createEmailJobHandler(c *fiber.Ctx) error {
	var email Email
	if err := c.BodyParser(&email); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if email.To == "" || email.Subject == "" {
		return c.Status(400).JSON(fiber.Map{"error": "To and subject are required"})
	}

	// Create job payload
	payload, _ := json.Marshal(email)

	job := Job{
		ID:      generateJobID(),
		Type:    "email",
		Payload: string(payload),
	}

	queue.AddJob(job)

	return c.Status(202).JSON(fiber.Map{
		"success": true,
		"job_id":  job.ID,
		"message": "Email job queued successfully",
	})
}

func createNotificationJobHandler(c *fiber.Ctx) error {
	var notification Notification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if notification.UserID == 0 || notification.Message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "UserID and message are required"})
	}

	// Create job payload
	payload, _ := json.Marshal(notification)

	job := Job{
		ID:      generateJobID(),
		Type:    "notification",
		Payload: string(payload),
	}

	queue.AddJob(job)

	return c.Status(202).JSON(fiber.Map{
		"success": true,
		"job_id":  job.ID,
		"message": "Notification job queued successfully",
	})
}

func getJobsHandler(c *fiber.Ctx) error {
	jobs := queue.GetJobs()

	// Calculate statistics
	stats := map[string]int{
		"total":      len(jobs),
		"pending":    0,
		"processing": 0,
		"completed":  0,
		"failed":     0,
	}

	for _, job := range jobs {
		stats[job.Status]++
	}

	return c.JSON(fiber.Map{
		"success": true,
		"jobs":    jobs,
		"stats":   stats,
	})
}

func getJobHandler(c *fiber.Ctx) error {
	jobID := c.Params("id")

	job := queue.GetJobByID(jobID)
	if job == nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Job not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"job":     job,
	})
}

// Legacy synchronous handlers for comparison
func sendEmailSyncHandler(c *fiber.Ctx) error {
	var email Email
	if err := c.BodyParser(&email); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// This blocks the request for 2 seconds
	time.Sleep(2 * time.Second)

	log.Printf("📧 Email sent synchronously to %s: %s", email.To, email.Subject)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email sent synchronously",
	})
}

func sendNotificationSyncHandler(c *fiber.Ctx) error {
	var notification Notification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// This blocks the request for 500ms
	time.Sleep(500 * time.Millisecond)

	log.Printf("🔔 Notification sent synchronously to user %d: %s", notification.UserID, notification.Message)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification sent synchronously",
	})
}

// Dashboard handler
func dashboardHandler(c *fiber.Ctx) error {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Message Queue Dashboard</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 20px;
            text-align: center;
        }
        .stats { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px; 
            margin-bottom: 20px; 
        }
        .stat-card { 
            background: white;
            border: 1px solid #ddd; 
            padding: 20px; 
            border-radius: 10px; 
            text-align: center;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .stat-value { 
            font-size: 2em; 
            font-weight: bold; 
            margin: 10px 0;
        }
        .jobs-list {
            background: white;
            border-radius: 10px;
            padding: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .job {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px;
            margin: 5px 0;
            border: 1px solid #eee;
            border-radius: 5px;
        }
        .pending { background-color: #fff3cd; }
        .processing { background-color: #d1ecf1; }
        .completed { background-color: #d4edda; }
        .failed { background-color: #f8d7da; }
        .controls {
            margin: 20px 0;
            text-align: center;
        }
        button {
            background: #007bff;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 5px;
            margin: 0 10px;
            cursor: pointer;
        }
        button:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="header">
        <h1>📬 Message Queue Dashboard</h1>
        <p>Real-time job monitoring and queue management</p>
    </div>
    
    <div class="stats">
        <div class="stat-card">
            <div>Total Jobs</div>
            <div class="stat-value" id="total-jobs">0</div>
        </div>
        <div class="stat-card">
            <div>Pending</div>
            <div class="stat-value" style="color: #856404;" id="pending-jobs">0</div>
        </div>
        <div class="stat-card">
            <div>Processing</div>
            <div class="stat-value" style="color: #0c5460;" id="processing-jobs">0</div>
        </div>
        <div class="stat-card">
            <div>Completed</div>
            <div class="stat-value" style="color: #155724;" id="completed-jobs">0</div>
        </div>
        <div class="stat-card">
            <div>Failed</div>
            <div class="stat-value" style="color: #721c24;" id="failed-jobs">0</div>
        </div>
    </div>
    
    <div class="controls">
        <button onclick="addTestEmail()">📧 Add Test Email</button>
        <button onclick="addTestNotification()">🔔 Add Test Notification</button>
        <button onclick="addMultipleJobs()">⚡ Add 5 Jobs</button>
    </div>
    
    <div class="jobs-list">
        <h3>Recent Jobs</h3>
        <div id="jobs-container">Loading...</div>
    </div>
    
    <script>
        function updateDashboard() {
            fetch('/jobs')
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        const stats = data.stats;
                        document.getElementById('total-jobs').textContent = stats.total;
                        document.getElementById('pending-jobs').textContent = stats.pending;
                        document.getElementById('processing-jobs').textContent = stats.processing;
                        document.getElementById('completed-jobs').textContent = stats.completed;
                        document.getElementById('failed-jobs').textContent = stats.failed;
                        
                        const container = document.getElementById('jobs-container');
                        container.innerHTML = '';
                        
                        const recentJobs = data.jobs.slice(-10).reverse();
                        recentJobs.forEach(job => {
                            const div = document.createElement('div');
                            div.className = 'job ' + job.status;
                            div.innerHTML = \`
                                <div>
                                    <strong>\${job.id}</strong> (\${job.type})
                                    <br><small>\${new Date(job.created_at).toLocaleString()}</small>
                                </div>
                                <div>\${job.status.toUpperCase()}</div>
                            \`;
                            container.appendChild(div);
                        });
                        
                        if (recentJobs.length === 0) {
                            container.innerHTML = '<p>No jobs yet. Add some test jobs!</p>';
                        }
                    }
                })
                .catch(error => console.error('Error:', error));
        }
        
        function addTestEmail() {
            fetch('/jobs/email', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    to: 'test@example.com',
                    subject: 'Test Email',
                    body: 'This is a test email from the queue system.'
                })
            });
        }
        
        function addTestNotification() {
            fetch('/jobs/notification', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    user_id: Math.floor(Math.random() * 1000),
                    message: 'Test notification message',
                    type: 'info'
                })
            });
        }
        
        function addMultipleJobs() {
            for (let i = 0; i < 5; i++) {
                setTimeout(() => {
                    if (i % 2 === 0) {
                        addTestEmail();
                    } else {
                        addTestNotification();
                    }
                }, i * 200);
            }
        }
        
        // Update every 2 seconds
        updateDashboard();
        setInterval(updateDashboard, 2000);
    </script>
</body>
</html>
    `

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}

// Utility functions
func generateJobID() string {
	return fmt.Sprintf("job_%d", time.Now().UnixNano())
} 