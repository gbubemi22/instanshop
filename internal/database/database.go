package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetGORM() *gorm.DB
}

type service struct {
	db     *sql.DB
	gormDB *gorm.DB
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

// New initializes the database service.
func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	// Correct DSN for PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)

	// Open a low-level database connection
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Set a reasonable lifetime
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(50)

	// Initialize GORM with the existing connection
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to initialize GORM: %v", err)
	}

	dbInstance = &service{
		db:     sqlDB,
		gormDB: gormDB,
	}
	return dbInstance
}

// Health checks the health of the database connection.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get stats
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	return stats
}

// Close closes the database connection.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	return s.db.Close()
}

// GetGORM returns the GORM instance.
func (s *service) GetGORM() *gorm.DB {
	return s.gormDB
}
