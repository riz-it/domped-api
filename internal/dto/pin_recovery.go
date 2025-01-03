package dto

// Request
type SetupWalletPINRequest struct {
	WalletID int64  `json:"wallet_id"`
	PinCode  string `json:"pin_code"`
}

// Response
