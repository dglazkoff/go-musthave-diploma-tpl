package models

type Withdrawals struct {
	ID          string `json:"order"` // использовать как primary key
	UserId      string `json:"-"`
	Sum         uint   `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}
