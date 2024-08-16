package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mikromolekula2002/Trade_Platform/internal/config"
	"gorm.io/gorm"
)

func ApplyMigrations(db *gorm.DB, cfg *config.Config) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("database: migrations.go/ApplyMigrations - ошибка работы GORM.\n ERROR: %v", err)
	}
	defer sqlDB.Close()
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("database: migrations.go/ApplyMigrations - ошибка драйвера postgres в GORM.\n ERROR: %v", err)
	}
	defer driver.Close()
	// Создание источника миграции
	fileSource, err := (&file.File{}).Open("file://" + cfg.Migration.MigrationPath)
	if err != nil {
		return fmt.Errorf("database: migrations.go/ApplyMigrations - ошибка открытия источника миграции.\n ERROR: %v", err)
	}

	m, err := migrate.NewWithInstance("file", fileSource, cfg.Database.DBName, driver)
	if err != nil {
		return fmt.Errorf("database: migrations.go/ApplyMigrations - ошибка создания экземпляра миграции.\n ERROR: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("database: migrations.go/ApplyMigrations - ошибка применения миграции.\n ERROR: %v", err)
	}

	return nil
}
