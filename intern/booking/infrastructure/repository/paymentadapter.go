package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/constants"
)

type PaymentAdapter struct {
	BaseURL string
	Client  *http.Client
	Logger  interfaces.Logger
}

func NewPaymentAdapter(baseURL string, client *http.Client, logger interfaces.Logger) *PaymentAdapter {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &PaymentAdapter{BaseURL: baseURL, Client: client, Logger: logger}
}

type paymentRequest struct {
	Amount  string `json:"amount"`
	Webhook string `json:"webhook"`
	Message string `json:"message"`
}

type paymentResponse struct {
	PaymentLink string `json:"link"`
}

func (p *PaymentAdapter) InitiatePayment(amount string, webhook string, message string, extra map[string]string) (string, error) {
	endpoint := p.BaseURL + constants.PaymentInitiateEndpoint
	fullWebhook := webhook
	reqBody := paymentRequest{Amount: amount, Webhook: fullWebhook, Message: message}
	b, _ := json.Marshal(reqBody)
	p.Logger.Debug(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, endpoint)

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Client.Do(req)
	if err != nil {
		p.Logger.Error(constants.EventHotelRequest, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, endpoint)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		p.Logger.Error(constants.EventHotelResponse, fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode), constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse, constants.KeyStatusCode, resp.StatusCode)
		return "", fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode)
	}
	var pr paymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		p.Logger.Error(constants.EventHotelResponse, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse)
		return "", err
	}
	p.Logger.Info(constants.EventPaymentInitiated, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventPaymentInitiated, constants.KeyPaymentLink, pr.PaymentLink)
	return pr.PaymentLink, nil
}
