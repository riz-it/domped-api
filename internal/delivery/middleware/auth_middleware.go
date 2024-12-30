package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"riz.it/domped/internal/domain"
)

func NewAuthMiddleware(auth domain.JWT) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		splitToken := strings.Split(ctx.Get("Authorization"), "Bearer ")
		if len(splitToken) < 2 {
			return fiber.ErrUnauthorized
		}

		accessToken := splitToken[1]

		userId, err := auth.ValidateAccessToken(accessToken)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		ctx.Locals("userId", userId)

		return ctx.Next()
	}
}
