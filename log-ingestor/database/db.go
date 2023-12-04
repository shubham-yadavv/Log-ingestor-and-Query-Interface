package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Log struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Level          string `json:"level" gorm:"index:idx_level"`
	Message        string `json:"message"`
	ResourceID     string `json:"resourceId" gorm:"index:idx_resource_id"`
	Timestamp      string `json:"timestamp" gorm:"index:idx_timestamp"`
	TraceID        string `json:"traceId" gorm:"index:idx_trace_id"`
	SpanID         string `json:"spanId" gorm:"index:idx_span_id"`
	Commit         string `json:"commit" gorm:"index:idx_commit"`
	ParentResource string `json:"parentResource" gorm:"index:idx_parent_resource"`
}

var DB *gorm.DB

func InitializeDB() {
	var err error

	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	// dbConnectionString := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=" + dbSSLMode
	dbConnectionString := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=" + dbSSLMode

	DB, err = gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		os.Exit(1)
	}

	err = DB.AutoMigrate(&Log{})
	if err != nil {
		log.Fatal("Failed to AutoMigrate:", err)
	}
}
