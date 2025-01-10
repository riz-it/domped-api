package util

import (
	"os"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ConfigureLogger mengatur logger dengan format dan output file
func ConfigureLogger(logDir, logFile string) logger.Config {
	// Tentukan path lengkap log file
	logPath := logDir + "/" + logFile

	// Pastikan direktori log ada
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755) // Membuat direktori jika belum ada
		if err != nil {
			panic(err)
		}
	}

	// Membuka atau membuat file log
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// Format log
	logFormat := `{"time": "${time}", "status": "${status}", "latency": "${latency}", "ip": "${ip}", "method": "${method}", "path": "${path}", "error": "${error}"}` + "\n"

	// Mengembalikan konfigurasi logger
	return logger.Config{
		Format: logFormat,
		Output: file,
	}
}
