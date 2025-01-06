package dto

// Request
type CreateNotificationRequest struct {
	UserID int64  `json:"user_id"`
	Status int8   `json:"status"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Response

// Data
type NotificationData struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Status    int8   `json:"status"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
