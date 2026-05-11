package database

import (
	"fmt"
	"time"

	"go-base-project/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type DB = sqlx.DB

func NewMySQLConnection(cfg config.DatabaseConfig) (*DB, error) {

	// format DSN MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi DB: %w", err)
	}

	// connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	// test koneksi
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("tidak bisa ping database: %w", err)
	}

	return db, nil
}