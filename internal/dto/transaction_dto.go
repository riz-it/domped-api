package dto

// Request
type TransferInquiryRequest struct {
	AccountNumber string `json:"account_number"`
	Amount        int64  `json:"amount"`
}

type TransferExecuteRequest struct {
	InquiryKey string `json:"inquiry_key"`
}

// Response
type TransferInquiryResponse struct {
	InquiryKey string `json:"inquiry_key"`
}

type TransferExecuteResponse struct {
	InquiryKey  string       `json:"inquiry_key"`
	Information TransferData `json:"information"`
}

// Data
type TransactionData struct {
	ID              int64  `json:"id"`
	WalletID        int64  `json:"wallet_id"`
	SofNumber       string `json:"sof_number"`
	DofNumber       string `json:"dof_number"`
	Amount          int64  `json:"amount"`
	TransactionType string `json:"transaction_type"`
	TransactionAt   string `json:"transaction_at"`
}

type TransferData struct {
	SofNumber     string `json:"sof_number"`
	DofNumber     string `json:"dof_number"`
	Amount        int64  `json:"amount"`
	TransactionAt string `json:"transaction_at"`
	Status        string `json:"status"`
}
