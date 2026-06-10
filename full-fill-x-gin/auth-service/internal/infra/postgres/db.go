package postgres

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	SlowThreshold   time.Duration
	LogLevel        string
}

type Database struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

func Open(cfg Config) (*Database, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(parseLogLevel(cfg.LogLevel)),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return &Database{db: db, sqlDB: sqlDB}, nil
}

func (d *Database) Gorm() *gorm.DB {
	return d.db
}

func (d *Database) Ping(ctx context.Context) error {
	return d.sqlDB.PingContext(ctx)
}

func (d *Database) Close() error {
	return d.sqlDB.Close()
}

func (d *Database) Migrate() error {
	return d.db.AutoMigrate(&domain.User{})
}

func parseLogLevel(value string) logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
