package main

import (
	"errors"
	// "flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	migrationsPath := "./migrations" // Путь к миграциям
	migrationsTable := "sso_migrations"

	fmt.Println(dbUser, dbPassword, dbHost, dbPort, dbName, migrationsPath, migrationsTable)

	// Формируем строку подключения
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, migrationsTable)
	log.Printf("Connecting to database with URL: %s", dbURL)

	// Выполняем миграции
	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		log.Println("migrate.New")
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
