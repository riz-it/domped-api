package dto

// Request
type SetupWalletPINRequest struct {
	WalletID            int64  `json:"wallet_id" validate:"required,numeric"`
	PinCode             string `json:"pin_code" validate:"required,numeric"`
	PinCodeConfirmation string `json:"pin_code_confirmation" validate:"required,numeric"`
}

// Response
