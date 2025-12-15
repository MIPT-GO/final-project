package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"final-project/intern/hotel/application/service/hotel_service"
	"final-project/intern/hotel/domain/entity"
	"final-project/pkg/custom_errors"

	mocks "final-project/tests/hotel/repository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var testLogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

func TestHotelService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockHotelRepository(ctrl)
	service := hotel_service.NewHotelService(testLogger, mockRepo)

	expectedHotels := []entity.Hotel{
		{Name: "Hilton", Email: "hilton@test.com"},
		{Name: "Marriott", Email: "marriott@test.com"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetAll().Return(expectedHotels, nil)

		names, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Hilton", "Marriott"}, names)
	})

	t.Run("Repository_Failure", func(t *testing.T) {
		repoError := errors.New("db error")
		mockRepo.EXPECT().GetAll().Return(nil, repoError)

		names, err := service.GetAll()

		assert.Error(t, err)
		assert.Nil(t, names)
	})
}

func TestHotelService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockHotelRepository(ctrl)
	service := hotel_service.NewHotelService(testLogger, mockRepo)

	hotelName := "New Hotel"
	hotelEmail := "new@hotel.com"
	newHotel := entity.Hotel{Name: hotelName, Email: hotelEmail}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().AddNewHotel(newHotel).Return(nil)

		err := service.Create(hotelName, hotelEmail)

		assert.NoError(t, err)
	})

	t.Run("Repository_Conflict", func(t *testing.T) {
		mockRepo.EXPECT().AddNewHotel(newHotel).Return(custom_errors.ErrEntityAlreadyExists)

		err := service.Create(hotelName, hotelEmail)

		assert.ErrorIs(t, err, custom_errors.ErrEntityAlreadyExists)
	})
}
