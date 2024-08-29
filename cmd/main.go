package main

// добавить REDIS для подсчета просмотров(переделать методы под редис)
// Исправить логику (я написал код подсчета(инкремента) просмотров объяв, но не сделал получение инфы о просмотрах)
// На самом деле если РЕДИС не подключится, можно просто вывести варн и не выводить просмотры
// Предусмотреть на такой случай обработку ошибок от редиса и обход их

// добавить ЮНИТ ТЕСты
// перепроверить логирование и обработку ошибок(как оно по факту будет выглядеть)
// СКОРЕЕ ВСЕГО добавить в DEBUG логирование всех методов, и прочее (АБСОЛЮТНО ВСЕГО)
// ПЕРЕДАТЬ ЛОГЕР ВО ВСЕ СТРУКТУРЫ ДЛЯ ЛЕВЕЛА DEBUG
// Добавить фронт

import (
	"github.com/mikromolekula2002/Trade_Platform/internal/config"
	"github.com/mikromolekula2002/Trade_Platform/internal/handler"
	"github.com/mikromolekula2002/Trade_Platform/internal/jwt"
	"github.com/mikromolekula2002/Trade_Platform/internal/repository"
	"github.com/mikromolekula2002/Trade_Platform/internal/service"
	"github.com/mikromolekula2002/Trade_Platform/pkg/database"
	"github.com/mikromolekula2002/Trade_Platform/pkg/logger"
)

func main() {
	//Загрузка конфига
	cfg := config.LoadConfig("./configs/config.yaml")

	//Инициализация логера(скорее всего исправить и передавать просто логрус без структуры ебучей)
	loger := logger.Init(cfg)

	//Прописать инит PostgreSQL
	db, err := database.InitDB(cfg)
	if err != nil {
		loger.Logrus.Fatal(err)
	}
	//НАВЕРНОЕ НАХУЙ СНЕСТИ ОТ
	sqlDB, err := db.DB()
	if err != nil {
		loger.Logrus.Info("хуй")
	}
	defer sqlDB.Close()
	//НАВЕРНОЕ НАХУЙ СНЕСТИ ДО
	loger.Logrus.Debug("PostgreSQL успешно подключен.")

	//Применяем миграции Базы Данных
	//	err = database.ApplyMigrations(db, cfg)
	//	if err != nil {
	//		loger.Logrus.Fatal(err)
	//	}
	//	loger.Logrus.Debug("Миграции успешно выполнены.")
	jwtManager := jwt.InitJWT()
	//Прописать инит Юзер Сервиса(валидация, работа с бд и прочее)
	userSvc := service.NewUserService(repository.PostgreInit(db), jwtManager, cfg.Jwt.JwtKey)

	// Прописать init Хендлера и маршруты
	h := handler.Init(loger.Logrus, cfg, userSvc)
	h.Routes()

	// Запуск сервера
	h.Echo.Logger.Fatal(h.Echo.StartServer(h.Echo.Server))
}
