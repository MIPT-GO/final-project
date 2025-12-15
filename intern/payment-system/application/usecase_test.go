package application_test

import (
	"bytes"
	"context"
	"errors"
	"final-project/intern/payment-system/application"
	"final-project/intern/payment-system/config"
	"final-project/intern/payment-system/domain"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockGenerator struct{}

func (generator *mockGenerator) Generate(info *domain.PaymentInfo) string {
	return "test-key"
}

type mockRepository struct {
	records map[string]domain.PaymentInfo
}

func (repository *mockRepository) New() {
	repository.records = make(map[string]domain.PaymentInfo)
}

func (repository *mockRepository) CreateRecord(ctx context.Context, key string, info domain.PaymentInfo) error {
	repository.records[key] = info
	return nil
}

func (repository *mockRepository) GetRecord(ctx context.Context, key string) (domain.PaymentInfo, error) {
	if info, ok := repository.records[key]; ok {
		return info, nil
	}
	return domain.PaymentInfo{}, errors.New("not found")
}

type mockRender struct {
	renderedPayment bool
	renderedMessage bool
	paymentData     application.PaymentPageData
	messageData     application.MessagePageData
	renderErr       error
}

func (render *mockRender) New(config *config.Config) {}

func (render *mockRender) RenderPaymentPage(writer io.Writer, data application.PaymentPageData) error {
	render.renderedPayment = true
	render.paymentData = data
	return nil
}

func (render *mockRender) RenderMessagePage(writer io.Writer, data application.MessagePageData) error {
	render.renderedMessage = true
	render.messageData = data
	return nil
}

type mockSender struct{}

func (sender *mockSender) Send(webhook string, status domain.PaymentStatus) error {
	return nil
}

func TestCreatePaymentLink(t *testing.T) {
	repository := mockRepository{}
	generator := mockGenerator{}
	render := mockRender{}
	sender := mockSender{}

	config := config.Config{
		LinkPrefix:     "http://test-url.ru/test/payment",
		PaymentTimeout: 1,
		RouterPrefix:   "/test/payment",
	}

	render.New(&config)
	repository.New()

	usecase := application.PaymentsUseCase{
		Repository: &repository,
		Generator:  &generator,
		Render:     &render,
		Sender:     &sender,
		Config:     &config,
	}

	link, err := usecase.CreatePaymentLink(t.Context(), "100", "http://test-webhook.com", "test")

	assert.NoError(t, err)
	assert.Equal(t, "http://test-url.ru/test/payment/test-key", link)

	assert.Equal(t, "100", repository.records["test-key"].Amount)
	assert.Equal(t, "http://test-webhook.com", repository.records["test-key"].WebHook)
	assert.Equal(t, "test", repository.records["test-key"].Message)
}

func TestRenderPaymentPage(t *testing.T) {
	repository := mockRepository{}
	generator := mockGenerator{}
	render := mockRender{}
	sender := mockSender{}

	config := config.Config{
		LinkPrefix:     "http://test-url.ru/test/payment",
		PaymentTimeout: 1,
		RouterPrefix:   "/test/payment",
	}

	render.New(&config)
	repository.New()

	usecase := application.PaymentsUseCase{
		Repository: &repository,
		Generator:  &generator,
		Render:     &render,
		Sender:     &sender,
		Config:     &config,
	}

	_, err := usecase.CreatePaymentLink(t.Context(), "100", "http://test-webhook.com", "test")
	assert.NoError(t, err)

	var buffer bytes.Buffer
	err = usecase.RenderPaymentPage(t.Context(), "test-key", &buffer)

	assert.NoError(t, err)

	assert.True(t, render.renderedPayment)
	assert.Equal(t, "100", render.paymentData.Amount)
	assert.Equal(t, "test", render.paymentData.Message)
}

func TestRenderPaymentPageTimeout(t *testing.T) {
	repository := mockRepository{}
	generator := mockGenerator{}
	render := mockRender{}
	sender := mockSender{}

	config := config.Config{
		LinkPrefix:     "http://test-url.ru/test/payment",
		PaymentTimeout: 0,
		RouterPrefix:   "/test/payment",
	}

	render.New(&config)
	repository.New()

	usecase := application.PaymentsUseCase{
		Repository: &repository,
		Generator:  &generator,
		Render:     &render,
		Sender:     &sender,
		Config:     &config,
	}

	repository.CreateRecord(t.Context(), "test-key", domain.PaymentInfo{
		Status: domain.TIMEOUT,
	})

	var buffer bytes.Buffer
	err := usecase.RenderPaymentPage(t.Context(), "test-key", &buffer)

	assert.NoError(t, err)

	assert.True(t, render.renderedMessage)
	assert.Equal(t, "The payment has already been confirmed or the time has expired", render.messageData.Message)
}

func TestConfirmPaymentSuccess(t *testing.T) {
	repository := mockRepository{}
	generator := mockGenerator{}
	render := mockRender{}
	sender := mockSender{}

	config := config.Config{
		LinkPrefix:     "http://test-url.ru/test/payment",
		PaymentTimeout: 1,
		RouterPrefix:   "/test/payment",
	}

	render.New(&config)
	repository.New()

	usecase := application.PaymentsUseCase{
		Repository: &repository,
		Generator:  &generator,
		Render:     &render,
		Sender:     &sender,
		Config:     &config,
	}

	repository.CreateRecord(t.Context(), "test-key", domain.PaymentInfo{
		Status: domain.WAITING,
	})

	err := usecase.ConfirmPayment(t.Context(), "test-key")

	assert.NoError(t, err)
	assert.Equal(t, domain.OK, repository.records["test-key"].Status)
}

func TestConfirmPaymentTimeout(t *testing.T) {
	repository := mockRepository{}
	generator := mockGenerator{}
	render := mockRender{}
	sender := mockSender{}

	config := config.Config{
		LinkPrefix:     "http://test-url.ru/test/payment",
		PaymentTimeout: 0,
		RouterPrefix:   "/test/payment",
	}

	render.New(&config)
	repository.New()

	usecase := application.PaymentsUseCase{
		Repository: &repository,
		Generator:  &generator,
		Render:     &render,
		Sender:     &sender,
		Config:     &config,
	}

	repository.CreateRecord(t.Context(), "test-key", domain.PaymentInfo{
		Status: domain.TIMEOUT,
	})

	err := usecase.ConfirmPayment(t.Context(), "test-key")

	assert.Error(t, err)
	assert.Equal(t, "The payment has already been confirmed or the time has expired", err.Error())
}
