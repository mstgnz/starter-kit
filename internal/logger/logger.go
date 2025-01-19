package logger

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

func logToDB(level string, message string) {
	/* db := config.App().DB.Create()
	stmt, err := config.App().DB.Prepare(config.App().QUERY["APP_LOG_INSERT"])
	if err == nil && level == "WARNING" {
		_, _ = stmt.Exec(level, message)
	}
	defer func() {
		_ = stmt.Close()
	}() */
}

func Info(messages ...string) {
	messages = append(messages, getCallerInfo())
	message := strings.Join(messages, ", ")
	slog.Info(message)
	logToDB("INFO", message)
}

func Error(messages ...string) {
	messages = append(messages, getCallerInfo())
	message := strings.Join(messages, ", ")
	slog.Error(message)
	logToDB("ERROR", message)
}

func Warn(messages ...string) {
	messages = append(messages, getCallerInfo())
	message := strings.Join(messages, ", ")
	slog.Error(message)
	logToDB("WARNING", message)
}

func Debug(messages ...string) {
	messages = append(messages, getCallerInfo())
	message := strings.Join(messages, ", ")
	slog.Debug(message)
	logToDB("DEBUG", message)
}

// The getCallerInfo function returns the file and line from which the log was called
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}

	fileParts := strings.Split(file, "/")
	fileName := fileParts[len(fileParts)-1]

	return fmt.Sprintf("%s:%d", fileName, line)
}
