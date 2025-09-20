package common

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	// Lấy connection string từ environment variable
	// Hỗ trợ cả PostgreSQL URL format và DSN format
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Default connection string - bạn cần thay thế bằng connection string của mình
		dsn = "postgresql://postgres:password@localhost:5432/plantheon?sslmode=disable"
	}

	// Cấu hình GORM để sử dụng IPv4
	config := &gorm.Config{}
	
	// Tạo PostgreSQL config với prefer_simple_protocol=true để tránh IPv6
	pgConfig := postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}

	db, err := gorm.Open(postgres.New(pgConfig), config)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
