package main

import (
	"fmt"
	"log"

	"messenger-auth/config"
	"messenger-auth/internal/models"
	"messenger-auth/internal/repositories"
)

func main() {
	appConfig := config.NewAppConfig()
	// Connect to database server
	dbServer, err := repositories.InitPostgresServer(appConfig)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
		return
	}
	defer dbServer.Close()

	// Try create database if not exists
	dbServer.Exec(fmt.Sprintf("CREATE DATABASE %s", appConfig.DBName))

	// Connect to database
	db, err := repositories.InitPostgres(appConfig)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
		return
	}
	defer db.Close()

	log.Println("Running migrations of tables...")
	db.AutoMigrate(
		&models.User{},
	)
	log.Println("Migrations completed")
}
