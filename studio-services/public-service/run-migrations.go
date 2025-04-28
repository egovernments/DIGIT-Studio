
package main

import (
	"log"
	"public-service/scripts" // 👈 import your `scripts` package
)

func main() {
	log.Println("Starting migrations...")
	scripts.RunMigrations()
	log.Println("Migrations completed.")
}

