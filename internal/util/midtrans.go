package util

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"riz.it/domped/internal/config"
	"riz.it/domped/internal/domain"
)

type MidtransUtil struct {
	Client         *snap.Client
	MidtransConfig *config.Midtrans
}

func NewMidtransUtil(cnf *config.Config) domain.Midtrans {
	var client snap.Client
	env := midtrans.Sandbox

	if cnf.Midtrans.IsProd {
		env = midtrans.Production
	}

	client.New(cnf.Midtrans.Key, env)

	return &MidtransUtil{
		Client:         &client,
		MidtransConfig: &cnf.Midtrans,
	}
}

// GenerateSnapURL implements domain.Midtrans.
func (m *MidtransUtil) GenerateSnapURL(ctx context.Context, t *domain.TopUpEntity) error {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  t.ID,
			GrossAmt: t.Amount,
		},
	}

	// 3. Request create Snap transaction to Midtrans
	snapResp, err := m.Client.CreateTransaction(req)

	if err != nil {
		return err
	}

	t.SnapURL = snapResp.RedirectURL
	return nil
}

// VerifyPayment implements domain.Midtrans.
func (m *MidtransUtil) VerifyPayment(ctx context.Context, data map[string]interface{}) (bool, error) {
	var client coreapi.Client
	env := midtrans.Sandbox

	if m.MidtransConfig.IsProd {
		env = midtrans.Production
	}

	client.New(m.MidtransConfig.Key, env)

	orderId, exists := data["order_id"].(string)
	if !exists {
		return false, domain.NewError(fiber.StatusBadRequest, "Invalid payload")
	}

	transactionStatusResp, e := client.CheckTransaction(orderId)
	if e != nil {
		return false, e
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				return true, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return false, nil
}
