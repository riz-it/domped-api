package dto

// Request
type SetupWalletPINRequest struct {
	PinCode             string `json:"pin_code" validate:"required,numeric"`
	PinCodeConfirmation string `json:"pin_code_confirmation" validate:"required,numeric"`
}

// Response
