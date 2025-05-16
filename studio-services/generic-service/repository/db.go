package repository

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB
var config map[string]string

// loadProperties reads application.properties and loads into a map
func loadProperties(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	properties := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // skip empty lines and comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		properties[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return properties, nil
}

// InitDB initializes the database connection
func InitDB() {
	var err error
	config, err = loadProperties("application.properties")
	if err != nil {
		log.Fatal("Error reading application.properties: ", err)
	}

	user := config["db.user"]
	password := config["db.password"]
	host := config["db.host"]
	port := config["db.port"]
	dbname := config["db.name"]

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("Database configuration is not set properly in application.properties")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error while opening DB connection: ", err)
	}

	// Verify if the database is alive
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
