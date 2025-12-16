package repository

type PaymentPort interface {
	InitiatePayment(amount string, webhook string, message string, extra map[string]string) (string, error)
}
