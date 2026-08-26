package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"notification-service/internal/config"
	"notification-service/internal/handler"
	"notification-service/internal/rabbitmq"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/internal/websocket"
)

func main() {
	cfg := config.Load()

	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Cannot connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database connected")

	// Setup WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(repo, hub)
	
	// Setup RabbitMQ Consumer
	rmqConsumer, err := rabbitmq.NewConsumer(cfg.RabbitMQURL, svc)
	if err != nil {
		log.Printf("⚠️ Could not connect to RabbitMQ: %v. Running without AMQP.", err)
	} else {
		defer rmqConsumer.Close()
		rmqConsumer.StartConsuming()
		log.Println("✅ RabbitMQ Consumer started")
	}

	h := handler.NewNotificationHandler(svc, hub)
	r := handler.SetupRouter(h, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Notification Service running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}

func connectDB(dsn string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	for i := 1; i <= 5; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			return db, nil
		}
		log.Printf("⏳ DB connect attempt %d/5 failed: %v", i, err)
		time.Sleep(time.Duration(i) * 2 * time.Second)
	}
	return nil, err
}
