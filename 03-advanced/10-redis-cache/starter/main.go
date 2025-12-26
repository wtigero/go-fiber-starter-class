package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/redis/go-redis/v9"
)

// TODO: สร้าง Redis client
// TODO: สร้าง Cache middleware
// TODO: สร้าง Rate limit middleware
// TODO: สร้าง Session middleware

var rdb *redis.Client
var ctx = context.Background()

func main() {
	// TODO: เชื่อมต่อ Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// ทดสอบการเชื่อมต่อ
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("ไม่สามารถเชื่อมต่อ Redis:", err)
	}
	log.Println("เชื่อมต่อ Redis สำเร็จ")

	app := fiber.New()
	app.Use(logger.New())

	// TODO: เพิ่ม routes และ middleware

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Redis Cache Starter",
			"status":  "ok",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
