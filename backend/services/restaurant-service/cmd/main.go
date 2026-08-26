package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"restaurant-service/internal/config"
	"restaurant-service/internal/handler"
	"restaurant-service/internal/repository"
	"restaurant-service/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Cannot connect to database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Println("✅ Database connected")

	// Init layers
	restRepo := repository.NewRestaurantRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	
	svc := service.NewRestaurantService(restRepo, menuRepo)
	reviewSvc := service.NewReviewService(reviewRepo, restRepo)
	
	h := handler.NewRestaurantHandler(svc)
	reviewHandler := handler.NewReviewHandler(reviewSvc)
	r := handler.SetupRouter(h, reviewHandler, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Restaurant Service running on %s", addr)
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
