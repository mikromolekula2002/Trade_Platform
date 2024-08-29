package test

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/mikromolekula2002/Trade_Platform/internal/config"
	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createDatabase() error {
	//Загрузка конфига
	cfg := config.LoadConfig("./config.yaml")
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=%s", cfg.TestDatabase.Host, cfg.TestDatabase.User, cfg.TestDatabase.Password, cfg.TestDatabase.Port, cfg.TestDatabase.Sslmode)
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PostgreSQL: %v", err)
	}
	defer conn.Close(context.Background())

	// Проверка существования базы данных
	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'testdb')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("не удалось проверить существование базы данных: %v", err)
	}

	if exists {
		// База данных существует
		return nil
	}

	// Создание базы данных
	_, err = conn.Exec(context.Background(), "CREATE DATABASE testdb")
	if err != nil {
		return fmt.Errorf("не удалось создать базу данных: %v", err)
	}

	return nil
}

func setupTestDB() (*gorm.DB, error) {
	// Создание базы данных
	err := createDatabase()
	if err != nil {
		return nil, err
	}

	//Загрузка конфига
	cfg := config.LoadConfig("./config.yaml")
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=%s", cfg.TestDatabase.Host, cfg.TestDatabase.User, cfg.TestDatabase.Password, cfg.TestDatabase.Port, cfg.TestDatabase.Sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}

	// Создание таблиц для тестов
	err = db.AutoMigrate(
		&models.UserData{},
		&models.UserAds{},
		&models.User{},
		&models.Likes{},
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать таблицы в базе данных: %v", err)
	}

	return db, nil
}

// Очистка данных после тестов
func truncateTables(db *gorm.DB) error {
	tables := []string{"user_data", "user_ads", "users", "likes"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
			return fmt.Errorf("не удалось очистить таблицу %s: %v", table, err)
		}
	}
	return nil
}
