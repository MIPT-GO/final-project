package api

type PaymentRequest struct {
	Amount  string `json:"amount"`
	WebHook string `json:"webhook"`
	Message string `json:"message"`
}
