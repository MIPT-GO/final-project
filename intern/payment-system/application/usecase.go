package application

import (
	"context"
	"errors"
	"final-project/intern/payment-system/config"
	"final-project/intern/payment-system/domain"
	"final-project/pkg/payment-system/constants"
	"io"
	"log/slog"
	"time"
)

type PaymentsUseCase struct {
	Repository PaymentInfoRepository
	Generator  KeyGenerator
	Render     Render
	Sender     WebHookSender
	Config     *config.Config
}

func (manager *PaymentsUseCase) CreatePaymentLink(ctx context.Context, amount string, webhook string, message string) (string, error) {
	record := domain.PaymentInfo{
		Amount:  amount,
		WebHook: webhook,
		Message: message,
		Status:  domain.WAITING,
	}

	key := manager.Generator.Generate(&record)

	err := manager.Repository.CreateRecord(ctx, key, record)
	if err != nil {
		return "", err
	}

	go time.AfterFunc(time.Duration(manager.Config.PaymentTimeout)*time.Minute, func() {
		info, err := manager.Repository.GetRecord(context.Background(), key)
		if err != nil || info.Status != domain.WAITING {
			return
		}

		info.Status = domain.TIMEOUT
		err = manager.Repository.CreateRecord(context.Background(), key, info)
		if err != nil {
			slog.Error(constants.MsgFailedUpdateRecord, constants.KeyError, err.Error(), constants.KeyPaymentKey, key)
		}

		err = manager.Sender.Send(info.WebHook, info.Status)
		if err != nil {
			slog.Error(constants.MsgFailedSendRequest,
				constants.KeyError, err.Error(),
				constants.KeyPaymentKey, key,
				constants.KeyWebHook, webhook,
			)
		}
	})

	return manager.Config.LinkPrefix + "/" + key, nil
}

func (manager *PaymentsUseCase) RenderPaymentPage(ctx context.Context, key string, writer io.Writer) error {
	info, err := manager.Repository.GetRecord(ctx, key)
	if err != nil || info.Status != domain.WAITING {
		err = manager.Render.RenderMessagePage(writer, MessagePageData{constants.MsgRecordExpired})
		return err
	}

	data := PaymentPageData{
		Amount:  info.Amount,
		Message: info.Message,
		Key:     key,
		Prefix:  manager.Config.RouterPrefix,
	}

	err = manager.Render.RenderPaymentPage(writer, data)
	return err
}

func (manager *PaymentsUseCase) ConfirmPayment(ctx context.Context, key string) error {
	info, err := manager.Repository.GetRecord(ctx, key)
	if err != nil {
		return err
	}

	if info.Status != domain.WAITING {
		return errors.New(constants.MsgRecordExpired)
	}

	info.Status = domain.OK

	err = manager.Repository.CreateRecord(ctx, key, info)

	go func() {
		err := manager.Sender.Send(info.WebHook, info.Status)
		if err != nil {
			slog.Error(constants.MsgFailedSendRequest,
				constants.KeyError, err.Error(),
				constants.KeyPaymentKey, key,
				constants.KeyWebHook, info.WebHook,
			)
		}
	}()

	return err
}
