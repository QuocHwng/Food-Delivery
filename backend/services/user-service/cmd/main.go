package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"user-service/internal/config"
	"user-service/internal/handler"
	"user-service/internal/service"
)

func main() {
	// ── Config ───────────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── Database ─────────────────────────────────────────────────────────────
	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Cannot connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Println("✅ Database connected")

	// ── Services ─────────────────────────────────────────────────────────────
	jwtCfg := service.JWTConfig{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessExpire:  cfg.JWT.AccessExpire,
		RefreshExpire: cfg.JWT.RefreshExpire,
	}
	authSvc := service.NewAuthService(db, jwtCfg)
	userSvc := service.NewUserService(db)

	// ── Handlers & Router ────────────────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	r := handler.SetupRouter(authHandler, userHandler, cfg.JWT.AccessSecret)

	// ── Start ────────────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 User Service running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}

// connectDB thử kết nối tối đa 5 lần (Docker startup race condition)
func connectDB(dsn string) (*sqlx.DB, error) {
	var (
		db  *sqlx.DB
		err error
	)
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
