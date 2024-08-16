package logger

import (
	"os"

	"github.com/mikromolekula2002/Trade_Platform/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Logrus *logrus.Logger
}

func Init(cfg *config.Config) *Logger {
	logger := logrus.New()

	// Устанавливаем уровень логирования
	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logger.Warnf("Failed to parse log level %s, defaulting to 'info'", cfg.Logger.Level)
		level = logrus.InfoLevel // По умолчанию устанавливаем info
	}
	logger.SetLevel(level)
	// Устанавливаем форматтер для вывода в stdout
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,       // Включаем цветной вывод
		FullTimestamp:          true,       // Выводим полную дату и время
		TimestampFormat:        "15:04:25", // Формат даты и времени
		DisableLevelTruncation: true,       // Отключаем усечение уровня логирования
		QuoteEmptyFields:       true,       // Кавычки для пустых полей
	})

	switch cfg.Logger.Output {
	case "file":
		file, err := os.OpenFile(cfg.Logger.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			logger.Info("Failed to open log file, using default stderr")
		}
	default:
		// По умолчанию вывод в stdout
		logger.SetOutput(os.Stdout)
	}
	// Установка вывода логов в stdout
	logger.SetOutput(os.Stdout)

	return &Logger{
		Logrus: logger,
	}
}
