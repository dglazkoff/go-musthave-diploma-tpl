package models

type Withdrawals struct {
	ID          string  `json:"order"` // использовать как primary key
	UserId      string  `json:"-"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
