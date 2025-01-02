package delivery

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"riz.it/domped/internal/delivery/controller"
)

type RouterConfig struct {
	App *fiber.App
}

func NewRouter(r *fiber.App, auth fiber.Handler, authController *controller.AuthController, transactionController *controller.TransactionController) *RouterConfig {
	// Logger configure
	logConfig := configureLogger("./logs", "access_log.json")
	r.Use(logger.New(logConfig))

	// Middleware CORS
	r.Use(cors.New())

	// Route
	/// Auth
	r.Post("/auth/login", authController.Login)
	r.Post("/auth/register", authController.Register)
	r.Post("/auth/refresh", authController.Refresh)
	r.Delete("/auth/logout", auth, authController.Logout)
	r.Post("/auth/verify", authController.EmailVerification)

	/// Transaction
	r.Post("/transaction/transfer/inquiry", auth, transactionController.Inquiry)
	r.Post("/transaction/transfer/execute", auth, transactionController.Execute)

	// Mengembalikan RouterConfig
	return &RouterConfig{
		App: r,
	}
}

// ConfigureLogger output file
func configureLogger(logDir, logFile string) logger.Config {
	logPath := logDir + "/" + logFile

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logFormat := `{"time": "${time}", "status": "${status}", "latency": "${latency}", "ip": "${ip}", "method": "${method}", "path": "${path}", "error": "${error}"}` + "\n"

	return logger.Config{
		Format: logFormat,
		Output: file,
	}
}
