package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// глобальный логгер
var log *slog.Logger

// Init инициализирует логгер
func Init(exportLevel string) {
	pathToFile := os.Getenv("LOG_FILE")

	// Собираем разветвитель (куда писать)
	writers := []io.Writer{os.Stdout} // в консоль

	if pathToFile != "" {
		// Открываем файл для дозаписи
		file, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			writers = append(writers, file)
		} else {
			// Если файл не открылся хотя бы предупредим в консоль
			fmt.Println("Невозможно открыть файл:", "error", err)
		}
	}

	multi := io.MultiWriter(writers...)

	level := levelChange(exportLevel)
	var handler slog.Handler
	handler = slog.NewTextHandler(multi, &slog.HandlerOptions{Level: level})

	// Сохраняем в глобальную переменную пакета
	log = slog.New(handler)
}

// Смена уровня логгирования
func levelChange(exportLevel string) slog.Level {
	if exportLevel == "" {
		fmt.Println("DEBUG: Не задан уровень логгирования")
	}
	switch exportLevel {
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func Debug(msg string, args ...any) {
	if log != nil {
		log.Debug(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if log != nil {
		log.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if log != nil {
		log.Warn(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if log != nil {
		log.Error(msg, args...)
	}
}
