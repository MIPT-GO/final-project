package room_service

import (
	"errors"
	"log/slog"
	"testing"

	"final-project/intern/hotel/domain/entity"
	"final-project/intern/hotel/domain/interfaces"
	"final-project/pkg/custom_errors"
	mock "final-project/tests/mock"
)

func setup() (*mock.RoomRepositoryMock, interfaces.RoomService) {
	rRepo := new(mock.RoomRepositoryMock)
	logger := slog.Default()
	svc := NewRoomService(logger, rRepo)
	return rRepo, svc
}

func TestCreate_Success(t *testing.T) {
	rRepo, svc := setup()
	rRepo.On("AddNewRoom", entity.Room{HotelName: "HotelA", Number: 101, Cost: 100.0}).Return(nil)

	err := svc.Create("HotelA", 101, 100.0)
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	rRepo.AssertExpectations(t)
}

func TestCreate_AlreadyExists(t *testing.T) {
	rRepo, svc := setup()
	rRepo.On("AddNewRoom", entity.Room{HotelName: "HotelA", Number: 101, Cost: 100.0}).Return(custom_errors.ErrEntityAlreadyExists)

	err := svc.Create("HotelA", 101, 100.0)
	if !errors.Is(err, custom_errors.ErrEntityAlreadyExists) {
		t.Fatalf("expected ErrEntityAlreadyExists, got %v", err)
	}
}

func TestCreate_RepoFailure(t *testing.T) {
	rRepo, svc := setup()

	rRepo.On("AddNewRoom", entity.Room{HotelName: "HotelA", Number: 1, Cost: 10.0}).Return(custom_errors.ErrDatabaseFailure)

	err := svc.Create("HotelA", 1, 10.0)
	if !errors.Is(err, custom_errors.ErrDatabaseFailure) {
		t.Fatalf("expected ErrDatabaseFailure, got %v", err)
	}
}

func TestUpdateCost_Success(t *testing.T) {
	rRepo, svc := setup()

	rRepo.On("UpdateCost", "HotelA", 101, float32(120.0)).Return(nil)

	err := svc.UpdateCost("HotelA", 101, 120.0)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	rRepo.AssertExpectations(t)
}

func TestUpdateCost_NotFound(t *testing.T) {
	rRepo, svc := setup()

	rRepo.On("UpdateCost", "HotelA", 999, float32(120.0)).Return(custom_errors.ErrEntityNotFound)

	err := svc.UpdateCost("HotelA", 999, 120.0)
	if !errors.Is(err, custom_errors.ErrEntityNotFound) {
		t.Fatalf("expected ErrEntityNotFound, got %v", err)
	}
}
