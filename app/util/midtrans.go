package util

import (
	"context"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"riz.it/domped/app/config"
	"riz.it/domped/app/domain"
)

type MidtransUtil struct {
	Config      *config.Midtrans
	Environment midtrans.EnvironmentType
}

func NewMidtransUtil(cnf *config.Config) domain.Midtrans {
	env := midtrans.Sandbox

	if cnf.Midtrans.IsProd {
		env = midtrans.Production
	}

	return &MidtransUtil{
		Environment: env,
		Config:      &cnf.Midtrans,
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

	var client snap.Client
	client.New(m.Config.Key, m.Environment)

	// 3. Request create Snap transaction to Midtrans
	snapResp, err := client.CreateTransaction(req)

	if err != nil {
		return err
	}

	t.SnapURL = snapResp.RedirectURL
	return nil
}

// VerifyPayment implements domain.Midtrans.
func (m *MidtransUtil) VerifyPayment(ctx context.Context, orderID string) (bool, error) {
	var client coreapi.Client
	client.New(m.Config.Key, m.Environment)

	transactionStatusResp, e := client.CheckTransaction(orderID)
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
