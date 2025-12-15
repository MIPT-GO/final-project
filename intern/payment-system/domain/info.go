package domain

type PaymentStatus string

const (
	WAITING PaymentStatus = "waiting"
	OK      PaymentStatus = "ok"
	TIMEOUT PaymentStatus = "timeout"
)

type PaymentInfo struct {
	Amount  string
	WebHook string
	Message string
	Status  PaymentStatus
}
