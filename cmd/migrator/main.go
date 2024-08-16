package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mikromolekula2002/Trade_Platform/internal/config"
)

func main() {
	//Загрузка конфига
	cfg := config.LoadConfig("./configs/config.yaml")
	// Получение пути к миграциям и параметров подключения к базе данных из окружения или конфигурации
	migrationPath := "file://" + cfg.Migration.MigrationPath
	databaseURL := "postgres://trader:tradertestpopa@localhost:5432/tradeplatform?sslmode=disable"

	// Создание мигратора
	m, err := migrate.New(migrationPath, databaseURL)
	if err != nil {
		log.Fatalf("Ошибка создания мигратора: %v", err)
	}

	// Применение миграций
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("миграции уже были применены")
			return
		}
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	fmt.Println("Миграции успешно применены")
}
