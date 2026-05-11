package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"go-base-project/internal/model"
	"go-base-project/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger struct {
	*zap.SugaredLogger
}

func NewLogger(level string) *AppLogger {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := config.Build()
	return &AppLogger{logger.Sugar()}
}

func RequestLogger(logger *AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logEntry := fmt.Sprintf("[%d] %s %s (%s)", statusCode, method, path, duration)
		if statusCode >= 500 {
			logger.Error(logEntry)
		} else if statusCode >= 400 {
			logger.Warn(logEntry)
		} else {
			logger.Info(logEntry)
		}
	}
}

func Recovery(logger *AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("PANIC: %v\n%s", err, debug.Stack()))

				c.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorResponse(
					"Terjadi kesalahan internal server", nil,
				))
			}
		}()
		c.Next()
	}
}

// ── CORS ─────────────────────────────────────────────────────
// CORS mengizinkan request dari domain lain (misalnya frontend di port berbeda)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Cek apakah origin diizinkan
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Preflight request (browser nanya dulu sebelum kirim request asli)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse(
				"Token tidak ditemukan, silakan login", nil,
			))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse(
				"Format token tidak valid (harus: Bearer <token>)", nil,
			))
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateJWT(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse(
				"Token tidak valid atau sudah kedaluwarsa", nil,
			))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, model.ErrorResponse(
			"Akses ditolak, role tidak mencukupi", nil,
		))
	}
}

