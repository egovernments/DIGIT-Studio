package scripts

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/magiconair/properties"
)

func RunMigrations() {
	// Load application.properties
	p, err := properties.LoadFile("application.properties", properties.UTF8)
	if err != nil {
		log.Fatalf("Could not load application.properties: %v", err)
	}

	dbUser := p.GetString("DB_USER", "postgres")
	dbPassword := p.GetString("DB_PASSWORD", "postgres")
	dbHost := p.GetString("DB_HOST", "localhost")
	dbPort := p.GetString("DB_PORT", "5432")
	dbName := p.GetString("DB_NAME", "egov_generic_db")

	// Build the database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Build the migrations folder path
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	migrationsDir := filepath.Join(wd, "scripts/migrations")

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Could not read migrations directory: %v", err)
	}

	// Run each SQL file
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(migrationsDir, file.Name())
			log.Printf("Running migration: %s", filePath)

			sqlBytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Could not read SQL file %s: %v", filePath, err)
			}

			sqlQuery := string(sqlBytes)

			_, err = db.Exec(sqlQuery)
			if err != nil {
				log.Fatalf("Migration failed for %s: %v", filePath, err)
			}

			log.Printf("Migration applied: %s", filePath)
		}
	}
}
