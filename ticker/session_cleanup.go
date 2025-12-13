package ticker

import (
	"log"
	"time"

	"cozeva.com/vault/pkg/helper"
)

func StartSessionCleanup() {
	ticker := time.NewTicker(12 * time.Hour) // triggers every 12 hour

	go func() { // runs in background goroutine
		for range ticker.C { // waits for ticker event
			if err := helper.CleanupExpiredSessions(); err != nil {
				log.Printf("session cleanup error: %v", err)
			}
		}
	}()
}
