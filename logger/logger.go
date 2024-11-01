package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	LogsDir = "data/logs"
)

type LogEntry struct {
	Message    string         `json:"message"`
	Level      string         `json:"level"`
	Time       time.Time      `json:"time"`
	Properties map[string]any `json:"properties"`
}

func Init() {
	err := os.MkdirAll(LogsDir, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create logs directory: %v\n", err)
	}
}

func Log(e LogEntry) {
	e.Time = time.Now()

	data, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("failed to marshal log entry: %v\n", err)
		return
	}
	data = append(data, '\n')

	filePath := fmt.Sprintf("%s/%s.json", LogsDir, e.Time.Format("2006-01-02"))

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		fmt.Printf("failed to write log entry: %v\n", err)
		return
	}

	fmt.Println(string(data))
}

func Info(message string) {
	Log(LogEntry{
		Message:    message,
		Level:      "info",
		Properties: map[string]any{},
	})
}

func InfoProps(message string, properties map[string]any) {
	Log(LogEntry{
		Message:    message,
		Level:      "info",
		Properties: properties,
	})
}
