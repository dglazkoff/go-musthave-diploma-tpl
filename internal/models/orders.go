package models

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing             = "PROCESSING"
	Invalid                = "INVALID"
	Processed              = "PROCESSED"
)

type Order struct {
	ID         string      `json:"number"`
	UserID     string      `json:"userID"`
	Status     OrderStatus `json:"status"`
	UploadedAt string      `json:"uploaded_at"`
	Accrual    int         `json:"accrual"`
}
