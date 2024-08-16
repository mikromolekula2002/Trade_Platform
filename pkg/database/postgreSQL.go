package database

import (
	"fmt"

	"github.com/mikromolekula2002/Trade_Platform/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.Sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("database: postgreSQL.go/InitDB - ошибка работы GORM.\n%v", err)
	}

	// Проверка подключения
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: postgreSQL.go/InitDB - ошибка подключения к postgreSQL.\n%v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		defer sqlDB.Close()
		return nil, fmt.Errorf("database: postgreSQL.go/InitDB - ошибка проверки подключения postgreSQL.\n%v", err)
	}

	return db, nil
}
