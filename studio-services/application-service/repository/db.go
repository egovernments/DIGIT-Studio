package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net"
	"os"
)

var db *sql.DB
var config map[string]string

func InitDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		user, password, host, port, dbname,
	)

	// Force IPv4 (avoid [::1] issue)
	net.DefaultResolver.PreferGo = true
	net.DefaultResolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		return net.Dial("tcp4", address)
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error while opening DB connection: ", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database ping failed: ", err)
	}

	log.Println("Database connected successfully!")
}

// GetDB returns the database connection instance
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB() first.")
	}
	return db
}
