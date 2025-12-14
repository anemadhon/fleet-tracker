package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogMessage struct {
	Level string
	Msg   string
	Time  time.Time
}

type AsyncLogger struct {
	Ch chan LogMessage
}

var globalLogger *AsyncLogger

func Init(logDir string) {
	os.MkdirAll(logDir, 0755)

	globalLogger = &AsyncLogger{
		Ch: make(chan LogMessage, 200),
	}

	go func() {
		for msg := range globalLogger.Ch {
			date := time.Now().Format("2006-01-02")
			filePath := filepath.Join(logDir, date+".log")
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("logger error: %v", err)
				continue
			}

			line := msg.Time.Format(time.RFC3339) + " [" + msg.Level + "] " + msg.Msg + "\n"

			f.WriteString(line)
			f.Close()
		}
	}()
}

func Info(msg string) {
	if globalLogger == nil {
		return
	}

	globalLogger.Ch <- LogMessage{"INFO", msg, time.Now()}
}

func Error(msg string) {
	if globalLogger == nil {
		return
	}

	globalLogger.Ch <- LogMessage{"ERROR", msg, time.Now()}
}
