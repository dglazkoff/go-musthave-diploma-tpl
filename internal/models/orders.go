package models

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing OrderStatus = "PROCESSING"
	Invalid    OrderStatus = "INVALID"
	Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         string      `json:"number"`
	UserID     string      `json:"userID"`
	Status     OrderStatus `json:"status"`
	UploadedAt string      `json:"uploaded_at"`
	Accrual    float64     `json:"accrual"`
}
