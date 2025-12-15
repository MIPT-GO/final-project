package repository

import (
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"final-project/intern/hotel/domain/entity"
	"final-project/intern/hotel/domain/interfaces"
	"final-project/intern/hotel/infrastructure/repository"
	"final-project/pkg/custom_errors"
)

var roomTestLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func setupRoomRepo(t *testing.T) (interfaces.RoomRepository, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewRoomRepository(db, roomTestLogger)
	return repo, db, mock
}

func TestRoomGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"hotel_name", "number", "cost"}).
			AddRow("Hotel", 1, 100).
			AddRow("Hotel", 2, 200)

		mock.ExpectQuery(`SELECT hotel_name, number, cost FROM Rooms WHERE hotel_name = \$1`).
			WithArgs("Hotel").
			WillReturnRows(rows)

		rooms, err := repo.GetAll("Hotel")

		assert.NoError(t, err)
		assert.Len(t, rooms, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		mock.ExpectQuery(`SELECT hotel_name, number, cost FROM Rooms WHERE hotel_name = \$1`).
			WithArgs("Hotel").
			WillReturnError(errors.New("db error"))

		rooms, err := repo.GetAll("Hotel")

		assert.ErrorIs(t, err, custom_errors.ErrRepoQueryFailed)
		assert.Nil(t, rooms)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoomGetCost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"cost"}).
			AddRow(float32(200))

		mock.ExpectQuery(`SELECT cost FROM Rooms WHERE hotel_name = \$1 AND number = \$2`).
			WithArgs("Hotel", 10).
			WillReturnRows(rows)

		cost, err := repo.GetCost("Hotel", 10)
		assert.NoError(t, err)
		assert.Equal(t, float32(200), cost)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		mock.ExpectQuery(`SELECT cost FROM Rooms WHERE hotel_name = \$1 AND number = \$2`).
			WithArgs("Hotel", 99).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetCost("Hotel", 99)
		assert.ErrorIs(t, err, custom_errors.ErrEntityNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoomAddNewRoom(t *testing.T) {
	room := entity.Room{
		HotelName: "Hotel",
		Number:    1,
		Cost:      120,
	}

	t.Run("success", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		mock.ExpectExec(`INSERT INTO Rooms`).
			WithArgs(room.HotelName, room.Number, room.Cost).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddNewRoom(room)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		mock.ExpectExec(`INSERT INTO Rooms`).
			WithArgs(room.HotelName, room.Number, room.Cost).
			WillReturnError(&pq.Error{Code: "23505"})

		err := repo.AddNewRoom(room)
		assert.ErrorIs(t, err, custom_errors.ErrEntityAlreadyExists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoomUpdateCost(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		mock.ExpectExec(`UPDATE Rooms`).
			WithArgs("Hotel", 1, float32(300)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateCost("Hotel", 1, 300)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		repo, db, mock := setupRoomRepo(t)
		defer db.Close()

		mock.ExpectExec(`UPDATE Rooms`).
			WithArgs("Hotel", 99, float32(300)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateCost("Hotel", 99, 300)
		assert.ErrorIs(t, err, custom_errors.ErrEntityNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
