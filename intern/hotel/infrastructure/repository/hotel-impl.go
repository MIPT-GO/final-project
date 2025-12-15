package repository

import (
	"database/sql"
	"errors"
	"final-project/intern/hotel/domain/entity"
	"final-project/intern/hotel/domain/interfaces"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
	"log/slog"

	"github.com/lib/pq"
)

type HotelRepositoryImpl struct {
	DB  *sql.DB
	Log *slog.Logger
}

func NewHotelRepository(db *sql.DB, log *slog.Logger) interfaces.HotelRepository {
	return &HotelRepositoryImpl{
		DB:  db,
		Log: log,
	}
}

func (r *HotelRepositoryImpl) GetAll() ([]entity.Hotel, error) {
	query := "SELECT name, email FROM hotels"
	rows, err := r.DB.Query(query)

	if err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBQuery,
			"sql", query,
			logs.KeyError, err)
		return nil, custom_errors.ErrRepoQueryFailed
	}
	defer rows.Close()

	var hotels []entity.Hotel
	for rows.Next() {
		var h entity.Hotel
		if err := rows.Scan(&h.Name, &h.Email); err != nil {
			r.Log.Error(logs.MsgDatabaseQueryFailed,
				logs.KeyEvent, logs.EventDBScan,
				"scan_target", "Hotel",
				logs.KeyError, err)
			return nil, custom_errors.ErrRepoScanFailed
		}
		hotels = append(hotels, h)
	}

	if err := rows.Err(); err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBQuery,
			"detail", "rows iteration failure",
			logs.KeyError, err)
		return nil, custom_errors.ErrRepoQueryFailed
	}

	return hotels, nil
}

func (r *HotelRepositoryImpl) GetByName(name string) (entity.Hotel, error) {
	query := "SELECT name, email FROM hotels WHERE name = $1"
	row := r.DB.QueryRow(query, name)

	var h entity.Hotel
	err := row.Scan(&h.Name, &h.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return entity.Hotel{}, custom_errors.ErrEntityNotFound
	}
	if err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBScan,
			logs.KeyHotelName, name,
			logs.KeyError, err)
		return entity.Hotel{}, custom_errors.ErrRepoScanFailed
	}

	return h, nil
}

func (r *HotelRepositoryImpl) GetByEmail(email string) (string, error) {
	var foundEmail string
	query := "SELECT email FROM hotels WHERE email = $1"
	row := r.DB.QueryRow(query, email)

	err := row.Scan(&foundEmail)

	if errors.Is(err, sql.ErrNoRows) {
		return "", custom_errors.ErrEntityNotFound
	}
	if err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBScan,
			logs.KeyHotelEmail, email,
			logs.KeyError, err)
		return "", custom_errors.ErrRepoScanFailed
	}

	return foundEmail, nil
}

func (r *HotelRepositoryImpl) AddNewHotel(hotel entity.Hotel) error {
	query := "INSERT INTO hotels (name, email) VALUES ($1, $2)"
	_, err := r.DB.Exec(query, hotel.Name, hotel.Email)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return custom_errors.ErrEntityAlreadyExists
		}

		r.Log.Error(logs.MsgDatabaseWriteFailed,
			logs.KeyEvent, logs.EventDBWrite,
			logs.KeyHotelName, hotel.Name,
			logs.KeyError, err)
		return custom_errors.ErrRepoQueryFailed
	}

	return nil
}

func (r *HotelRepositoryImpl) ExistsByName(name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM hotels WHERE name = $1)"

	err := r.DB.QueryRow(query, name).Scan(&exists)

	if err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBQuery,
			logs.KeyHotelName, name,
			logs.KeyError, err)
		return false, custom_errors.ErrRepoQueryFailed
	}

	return exists, nil
}
