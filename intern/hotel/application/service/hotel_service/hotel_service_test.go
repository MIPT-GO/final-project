package hotel_service

import (
	"io"
	"log/slog"
	"testing"

	"final-project/intern/hotel/domain/entity"
	"final-project/intern/hotel/domain/interfaces"
	"final-project/pkg/custom_errors"
	"final-project/tests/mock"

	"github.com/stretchr/testify/assert"
)

var testLogger = slog.New(
	slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}),
)

func setupService() (*mock.HotelRepositoryMock, interfaces.HotelService) {
	repoMock := new(mock.HotelRepositoryMock)
	service := NewHotelService(testLogger, repoMock)
	return repoMock, service
}

func TestHotelService_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, service := setupService()

		repo.On("GetAll").Return([]entity.Hotel{
			{Name: "A"},
			{Name: "B"},
		}, nil)

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, []string{"A", "B"}, result)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo, service := setupService()

		repo.On("GetAll").
			Return(nil, assert.AnError)

		result, err := service.GetAll()

		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestHotelService_GetEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, service := setupService()

		repo.On("GetByName", "Hotel").
			Return(entity.Hotel{Name: "Hotel", Email: "a@b.com"}, nil)

		email, err := service.GetEmail("Hotel")

		assert.NoError(t, err)
		assert.Equal(t, "a@b.com", email)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		repo, service := setupService()

		repo.On("GetByName", "Hotel").
			Return(nil, custom_errors.ErrEntityNotFound)

		email, err := service.GetEmail("Hotel")

		assert.ErrorIs(t, err, custom_errors.ErrEntityNotFound)
		assert.Empty(t, email)
		repo.AssertExpectations(t)
	})
}

func TestHotelService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, service := setupService()
		name := "Hotel"
		email := "a@b.com"
		newHotel := entity.Hotel{Name: name, Email: email}

		repo.On("AddNewHotel", newHotel).Return(nil)

		err := service.Create(name, email)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("already exists", func(t *testing.T) {
		repo, service := setupService()
		name := "Hotel"
		email := "a@b.com"
		newHotel := entity.Hotel{Name: name, Email: email}

		repo.On("AddNewHotel", newHotel).Return(custom_errors.ErrEntityAlreadyExists)

		err := service.Create(name, email)

		assert.ErrorIs(t, err, custom_errors.ErrEntityAlreadyExists)
		repo.AssertExpectations(t)
	})

	t.Run("db failure", func(t *testing.T) {
		repo, service := setupService()
		name := "Hotel"
		email := "a@b.com"
		newHotel := entity.Hotel{Name: name, Email: email}

		repo.On("AddNewHotel", newHotel).Return(custom_errors.ErrDatabaseFailure)

		err := service.Create(name, email)

		assert.ErrorIs(t, err, custom_errors.ErrDatabaseFailure)
		repo.AssertExpectations(t)
	})
}
