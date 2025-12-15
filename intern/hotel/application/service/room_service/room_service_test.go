package room_service

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

func setupRoomService() (*mock.HotelRepositoryMock, *mock.RoomRepositoryMock, interfaces.RoomService) {
	hotelRepo := new(mock.HotelRepositoryMock)
	roomRepo := new(mock.RoomRepositoryMock)
	service := NewRoomService(testLogger, hotelRepo, roomRepo)
	return hotelRepo, roomRepo, service
}

func TestRoomService_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()
		_ = hotelRepo

		roomRepo.On("GetAll", "HotelA").Return([]entity.Room{
			{Number: 101}, {Number: 102},
		}, nil)

		result, err := service.GetAll("HotelA")
		assert.NoError(t, err)
		assert.Equal(t, []int{101, 102}, result)
		roomRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()
		_ = hotelRepo

		roomRepo.On("GetAll", "HotelA").Return(nil, custom_errors.ErrDatabaseFailure)

		result, err := service.GetAll("HotelA")
		assert.ErrorIs(t, err, custom_errors.ErrDatabaseFailure)
		assert.Nil(t, result)
		roomRepo.AssertExpectations(t)
	})
}

func TestRoomService_GetCost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()
		_ = hotelRepo

		roomRepo.On("GetCost", "HotelA", 101).Return(float32(150), nil)

		cost, err := service.GetCost("HotelA", 101)
		assert.NoError(t, err)
		assert.Equal(t, float32(150), cost)
		roomRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()
		_ = hotelRepo

		roomRepo.On("GetCost", "HotelA", 101).Return(float32(0), custom_errors.ErrEntityNotFound)

		cost, err := service.GetCost("HotelA", 101)
		assert.ErrorIs(t, err, custom_errors.ErrEntityNotFound)
		assert.Equal(t, float32(0), cost)
		roomRepo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()
		_ = hotelRepo

		roomRepo.On("GetCost", "HotelA", 101).Return(float32(0), custom_errors.ErrDatabaseFailure)

		cost, err := service.GetCost("HotelA", 101)
		assert.ErrorIs(t, err, custom_errors.ErrDatabaseFailure)
		assert.Equal(t, float32(0), cost)
		roomRepo.AssertExpectations(t)
	})
}

func TestRoomService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()

		hotelRepo.On("GetByEmail", "HotelA").Return("owner@mail.com", nil)
		roomRepo.On("AddNewRoom", entity.Room{
			HotelName: "HotelA",
			Number:    101,
			Cost:      200,
		}).Return(nil)

		err := service.Create("HotelA", 101, 200, "owner@mail.com")
		assert.NoError(t, err)
		hotelRepo.AssertExpectations(t)
		roomRepo.AssertExpectations(t)
	})

	t.Run("permission denied", func(t *testing.T) {
		hotelRepo, _, service := setupRoomService()

		hotelRepo.On("GetByEmail", "HotelA").Return("owner@mail.com", nil)

		err := service.Create("HotelA", 101, 200, "other@mail.com")
		assert.ErrorIs(t, err, custom_errors.ErrPermissionDenied)
		hotelRepo.AssertExpectations(t)
	})

	t.Run("room already exists", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()

		hotelRepo.On("GetByEmail", "HotelA").Return("owner@mail.com", nil)
		roomRepo.On("AddNewRoom", entity.Room{
			HotelName: "HotelA",
			Number:    101,
			Cost:      200,
		}).Return(custom_errors.ErrEntityAlreadyExists)

		err := service.Create("HotelA", 101, 200, "owner@mail.com")
		assert.ErrorIs(t, err, custom_errors.ErrEntityAlreadyExists)
		hotelRepo.AssertExpectations(t)
		roomRepo.AssertExpectations(t)
	})
}

func TestRoomService_UpdateCost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()

		hotelRepo.On("GetByEmail", "HotelA").Return("owner@mail.com", nil)
		roomRepo.On("UpdateCost", "HotelA", 101, float32(250)).Return(nil)

		err := service.UpdateCost("HotelA", 101, 250, "owner@mail.com")
		assert.NoError(t, err)
		hotelRepo.AssertExpectations(t)
		roomRepo.AssertExpectations(t)
	})

	t.Run("permission denied", func(t *testing.T) {
		hotelRepo, _, service := setupRoomService()

		hotelRepo.On("GetByEmail", "HotelA").Return("owner@mail.com", nil)

		err := service.UpdateCost("HotelA", 101, 250, "other@mail.com")
		assert.ErrorIs(t, err, custom_errors.ErrPermissionDenied)
		hotelRepo.AssertExpectations(t)
	})

	t.Run("room not found", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()

		hotelRepo.On("GetByEmail", "HotelA").Return("owner@mail.com", nil)
		roomRepo.On("UpdateCost", "HotelA", 101, float32(250)).Return(custom_errors.ErrEntityNotFound)

		err := service.UpdateCost("HotelA", 101, 250, "owner@mail.com")
		assert.ErrorIs(t, err, custom_errors.ErrEntityNotFound)
		hotelRepo.AssertExpectations(t)
		roomRepo.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		hotelRepo, roomRepo, service := setupRoomService()

		hotelRepo.On("GetByEmail", "HotelA").Return("owner@mail.com", nil)
		roomRepo.On("UpdateCost", "HotelA", 101, float32(250)).Return(custom_errors.ErrDatabaseFailure)

		err := service.UpdateCost("HotelA", 101, 250, "owner@mail.com")
		assert.ErrorIs(t, err, custom_errors.ErrDatabaseFailure)
		hotelRepo.AssertExpectations(t)
		roomRepo.AssertExpectations(t)
	})
}
