package domain

import "context"

type Midtrans interface {
	GenerateSnapURL(ctx context.Context, t *TopUpEntity) error
	VerifyPayment(ctx context.Context, data map[string]interface{}) (bool, error)
}
