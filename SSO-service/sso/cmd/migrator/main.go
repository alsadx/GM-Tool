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
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	migrationsPath := "/migrations" // Путь к миграциям
	migrationsTable := "sso_migrations"

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

	// var dbURL, migrationsPath, migrationsTable string

	// flag.StringVar(&dbURL, "db-url", "", "PostgreSQL connection URL (e.g., postgres://user:password@host:port/dbname?sslmode=disable)")
	// flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	// flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	// flag.Parse()

	// if dbURL == "" {
	// 	panic("db-url is required")
	// }
	// if migrationsPath == "" {
	// 	panic("migrations-path is required")
	// }
	// log.Printf("Connecting to database with url: %s", dbURL)

	// m, err := migrate.New(
	// 	"file://"+migrationsPath,
	// 	dbURL,
	// )
	// if err != nil {
	// 	log.Println("migrate.New")
	// 	panic(err)
	// }

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
