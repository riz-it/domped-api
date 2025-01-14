package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"riz.it/domped/app/delivery/controller"
)

type RouterConfig struct {
	App *fiber.App
}

func NewRouter(r *fiber.App, auth fiber.Handler, authController *controller.AuthController, transactionController *controller.TransactionController, pinRecoveryController *controller.PinRecoveryController, notificationController *controller.NotificationController, topUpController *controller.TopUpController, mainController *controller.MainController) *RouterConfig {
	// Logger configure
	logFormat := `{"time": "${time}", "status": "${status}", "latency": "${latency}", "ip": "${ip}", "method": "${method}", "path": "${path}", "error": "${error}"}` + "\n"

	logConfig := logger.Config{
		Format: logFormat,
	}

	r.Use(logger.New(logConfig))

	// Middleware CORS
	r.Use(cors.New())

	// Route
	r.Get("/", mainController.Main)
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
	r.Post("/topup/callback", topUpController.Verify)

	// Mengembalikan RouterConfig
	return &RouterConfig{
		App: r,
	}
}
