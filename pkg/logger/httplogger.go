package httplogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type WinstonLog struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

type Logger struct {
	endpoint string
	client   *http.Client
	service  string
}

func NewLogger(endpoint, service string) *Logger {
	return &Logger{
		endpoint: endpoint,
		service:  service,
		client:   &http.Client{Timeout: 3 * time.Second},
	}
}

func (l *Logger) send(level, msg string, meta map[string]any) {
	logData := WinstonLog{
		Level:     level,
		Message:   fmt.Sprintf("%s", map[string]any{"msg": msg, "meta": meta}),
		Service:   os.Getenv("service-id"),
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	data, _ := json.Marshal(logData)
	req, err := http.NewRequest("POST", l.endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Println("failed to create log request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	go func() { // async send
		resp, err := l.client.Do(req)
		if err != nil {
			log.Printf("[WARN] Failed to send log to Winston: %v\n", err)
			return
		}
		defer resp.Body.Close()
	}()
}

func (l *Logger) Info(msg string, meta map[string]any)  { l.send("info", msg, meta) }
func (l *Logger) Error(msg string, meta map[string]any) { l.send("error", msg, meta) }
func (l *Logger) Warn(msg string, meta map[string]any)  { l.send("warn", msg, meta) }
