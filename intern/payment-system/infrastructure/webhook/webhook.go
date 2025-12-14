package webhook

import (
	"bytes"
	"encoding/json"
	"final-project/intern/payment-system/domain"
	"fmt"
	"net/http"
)

type WebHookMessage struct {
	Status string `json:"status"`
}

type NetWebHookSender struct{}

func (sender *NetWebHookSender) Send(webhook string, status domain.PaymentStatus) error {
	message := WebHookMessage{
		Status: string(status),
	}

	data, err := json.Marshal(message)

	if err != nil {
		return err
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("WebHook: %s answer with status code: %d", webhook, resp.StatusCode)
	}

	return nil
}
