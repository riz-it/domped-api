package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"riz.it/domped/app/delivery/controller"
)

type RouterConfig struct {
	App *fiber.App
}

func NewRouter(r *fiber.App, auth fiber.Handler, authController *controller.AuthController, transactionController *controller.TransactionController, pinRecoveryController *controller.PinRecoveryController, notificationController *controller.NotificationController, topUpController *controller.TopUpController) *RouterConfig {
	// Logger configure
	// logConfig := configureLogger("./logs", "access_log.json")
	// r.Use(logger.New(logConfig))

	// Middleware CORS
	r.Use(cors.New())

	// Route
	/// Auth
	r.Post("/auth/login", authController.Login)
	r.Post("/auth/register", authController.Register)
	r.Post("/auth/refresh", authController.Refresh)
	r.Delete("/auth/logout", auth, authController.Logout)
	r.Post("/auth/verify", authController.EmailVerification)

	/// Pin Recovery
	r.Post("/wallet/pin/recovery", auth, pinRecoveryController.SetupPin)

	/// Notification
	r.Get("/notifications", auth, notificationController.GetUserNotifications)

	/// Transaction
	r.Post("/transaction/transfer/inquiry", auth, transactionController.Inquiry)
	r.Post("/transaction/transfer/execute", auth, transactionController.Execute)

	/// TopUp
	r.Post("/topup/initialize", auth, topUpController.Initialize)
	r.Post("/topup/callback", auth, topUpController.Verify)

	// Mengembalikan RouterConfig
	return &RouterConfig{
		App: r,
	}
}

// ConfigureLogger configures the logger based on environment
// func configureLogger(logDir, logFile string) logger.Config {
// 	environment := os.Getenv("Environment") // Get the environment variable

// 	logFormat := `{"time": "${time}", "status": "${status}", "latency": "${latency}", "ip": "${ip}", "method": "${method}", "path": "${path}", "error": "${error}"}` + "\n"

// 	if environment == "development" || environment == "" {
// 		// Development environment: Write logs to file
// 		logPath := logDir + "/" + logFile

// 		if _, err := os.Stat(logDir); os.IsNotExist(err) {
// 			err := os.MkdirAll(logDir, 0755)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}

// 		file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 		if err != nil {
// 			panic(err)
// 		}

// 		return logger.Config{
// 			Format: logFormat,
// 			Output: file,
// 		}
// 	}

// 	// Default: Use stdout for other environments
// 	return logger.Config{
// 		Format: logFormat,
// 		Output: os.Stdout,
// 	}
// }
