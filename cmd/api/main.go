package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"go-base-project/config"
	"go-base-project/internal/database"
	"go-base-project/internal/handler"
	"go-base-project/internal/middleware"
	"go-base-project/internal/repository"
	"go-base-project/internal/service"
)

func main() {
	// ── Logger ────────────────────────────────────────────────────────────────
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Set log level from config
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
	if cfg.Log.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	log.Infof("Starting %s v%s [%s]", cfg.App.Name, cfg.App.Version, cfg.App.Env)

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := database.NewMySQL(&cfg.Database, log)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	// ── Repositories ──────────────────────────────────────────────────────────
	custRepo := repository.NewCustomerRepository(db)
	txRepo   := repository.NewTransactionRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	custSvc  := service.NewCustomerService(custRepo)
	txSvc   := service.NewTransactionService(txRepo, custRepo)
	pointSvc := service.NewPointService(custRepo, txRepo)

	// ── Handlers ──────────────────────────────────────────────────────────────
	authH    := handler.NewAuthHandler(cfg.JWT.Secret, cfg.JWT.ExpiryHours, log)
	custH    := handler.NewCustomerHandler(custSvc, log)
	txH      := handler.NewTransactionHandler(txSvc, log)
	pointH   := handler.NewPointHandler(pointSvc, log)

	// ── Router ────────────────────────────────────────────────────────────────
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "version": cfg.App.Version})
	})

	// API v1
	v1 := r.Group("/api/v1")

	// Public routes
	v1.POST("/auth/login", authH.Login)

	// Protected routes
	auth := v1.Group("/")
	auth.Use(middleware.JWTAuth(cfg.JWT.Secret))
	{
		// Customers
		auth.POST("/customers",     custH.Create)
		auth.GET("/customers",      custH.List)

		// Transactions
		auth.POST("/transactions",              txH.Create)
		auth.GET("/transactions",               txH.List)
		auth.GET("/transactions/passbook",      txH.Passbook)

		// Points
		auth.GET("/points", pointH.List)
	}

	// ── Start server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Infof("Server listening on %s", addr)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Info("Shutting down server...")
}
