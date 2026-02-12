package log

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func CreateLogger() (*slog.Logger, *os.File, error) {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	log := &lumberjack.Logger{
		Filename:   "./app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default

	}

	multiWriter := io.MultiWriter(log, os.Stdout)

	handler := slog.NewJSONHandler(multiWriter, nil)

	logger := slog.New(handler)

	return logger, file, nil
}
