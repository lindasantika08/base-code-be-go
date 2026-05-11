package main

// ======================File utama inisiasi configuration====================

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-base-project/config"
	"go-base-project/internal/database"
	"go-base-project/internal/handler"
	"go-base-project/internal/middleware"
	"go-base-project/internal/repository"
	"go-base-project/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(" Gagal load config: %v", err)
	}

	logger := middleware.NewLogger(cfg.App.LogLevel)
	logger.Info("Memulai aplikasi " + cfg.App.Name)

	db, err := database.NewMySQLConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Gagal konek ke database: " + err.Error())
	}
	defer db.Close()
	logger.Info(" Database terhubung")

	// database.RunMigrations(db, "./migrations")

	userRepo := repository.NewUserRepository(db)       
	userService := service.NewUserService(userRepo)    
	userHandler := handler.NewUserHandler(userService) 

	router := setupRouter(cfg, userHandler, logger)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server berjalan di http://localhost%s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(" Server error: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Menghentikan server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server paksa berhenti: " + err.Error())
	}
	logger.Info("Server berhenti dengan bersih")
}

func setupRouter(cfg *config.Config, userHandler *handler.UserHandler, logger *middleware.AppLogger) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.RequestLogger(logger))   // log setiap request masuk
	router.Use(middleware.Recovery(logger))       
	router.Use(middleware.CORS(cfg.App.AllowedOrigins)) 

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": cfg.App.Version,
		})
	})

	// ── API V1 ───────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(cfg.JWT.Secret)) 
		{
			users := protected.Group("/users")
			{
				users.GET("", userHandler.GetAll)       
				users.GET("/:id", userHandler.GetByID)  
				users.PUT("/:id", userHandler.Update)   
				users.DELETE("/:id", userHandler.Delete)
			}

			// ── TAMBAHKAN ROUTE BARU DI SINI ────────────────
		}
	}

	return router
}
