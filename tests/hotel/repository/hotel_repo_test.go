package repository

import (
	"database/sql"
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

var hotelTestLoger = slog.New(slog.NewTextHandler(nil, nil))

func setup(t *testing.T) (interfaces.HotelRepository, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewHotelRepository(db, hotelTestLoger)
	return repo, db, mock
}

func TestAddNewHotel(t *testing.T) {
	hotel := entity.Hotel{Name: "Test", Email: "test@mail.com"}

	t.Run("success", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		mock.ExpectExec(`INSERT INTO Hotels`).
			WithArgs(hotel.Name, hotel.Email).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddNewHotel(hotel)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		mock.ExpectExec(`INSERT INTO Hotels`).
			WithArgs(hotel.Name, hotel.Email).
			WillReturnError(&pq.Error{Code: "23505"})

		err := repo.AddNewHotel(hotel)
		assert.ErrorIs(t, err, custom_errors.ErrEntityAlreadyExists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name", "email"}).
			AddRow("Hotel", "hotel@mail.com")

		mock.ExpectQuery(`SELECT name, email FROM Hotels WHERE name = \$1`).
			WithArgs("Hotel").
			WillReturnRows(rows)

		h, err := repo.GetByName("Hotel")
		assert.NoError(t, err)
		assert.Equal(t, "Hotel", h.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		mock.ExpectQuery(`SELECT name, email FROM Hotels WHERE name = \$1`).
			WithArgs("Missing").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetByName("Missing")
		assert.ErrorIs(t, err, custom_errors.ErrEntityNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"email"}).
			AddRow("a@mail.com")

		mock.ExpectQuery(`SELECT email FROM Hotels WHERE email = \$1`).
			WithArgs("a@mail.com").
			WillReturnRows(rows)

		email, err := repo.GetByEmail("a@mail.com")
		assert.NoError(t, err)
		assert.Equal(t, "a@mail.com", email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		mock.ExpectQuery(`SELECT email FROM Hotels WHERE email = \$1`).
			WithArgs("missing@mail.com").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetByEmail("missing@mail.com")
		assert.ErrorIs(t, err, custom_errors.ErrEntityNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetAll(t *testing.T) {
	repo, db, mock := setup(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name", "email"}).
		AddRow("A", "a@mail.com").
		AddRow("B", "b@mail.com")

	mock.ExpectQuery(`SELECT name, email FROM Hotels`).
		WillReturnRows(rows)

	list, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExistsByName(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM Hotels WHERE name = \$1\)`).
			WithArgs("Hotel").
			WillReturnRows(rows)

		ok, err := repo.ExistsByName("Hotel")
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not exists", func(t *testing.T) {
		repo, db, mock := setup(t)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM Hotels WHERE name = \$1\)`).
			WithArgs("Missing").
			WillReturnRows(rows)

		ok, err := repo.ExistsByName("Missing")
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
