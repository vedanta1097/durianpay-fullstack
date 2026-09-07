package entity

import "time"

const (
	PaymentStatusCompleted  = "completed"
	PaymentStatusProcessing = "processing"
	PaymentStatusFailed     = "failed"
)

type Payment struct {
	ID        string
	Merchant  string
	Status    string
	Amount    int64
	CreatedAt time.Time
}

func IsSupportedPaymentStatus(status string) bool {
	return status == PaymentStatusCompleted || status == PaymentStatusProcessing || status == PaymentStatusFailed
}
