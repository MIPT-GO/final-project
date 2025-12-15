package repository

import (
	"database/sql"
	"errors"
	"log/slog"

	"final-project/intern/hotel/domain/entity"
	"final-project/intern/hotel/domain/interfaces"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"

	"github.com/lib/pq"
)

type RoomRepositoryImpl struct {
	DB  *sql.DB
	Log *slog.Logger
}

func NewRoomRepository(db *sql.DB, log *slog.Logger) interfaces.RoomRepository {
	return &RoomRepositoryImpl{
		DB:  db,
		Log: log,
	}
}

func (r *RoomRepositoryImpl) GetAll(hotelName string) ([]entity.Room, error) {
	query := "SELECT hotel_name, number, cost FROM rooms WHERE hotel_name = $1"
	rows, err := r.DB.Query(query, hotelName)

	if err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBQuery,
			logs.KeyHotelName, hotelName,
			logs.KeyError, err)
		return nil, custom_errors.ErrRepoQueryFailed
	}
	defer rows.Close()

	var rooms []entity.Room
	for rows.Next() {
		var rm entity.Room
		if err := rows.Scan(&rm.HotelName, &rm.Number, &rm.Cost); err != nil {
			r.Log.Error(logs.MsgDatabaseQueryFailed,
				logs.KeyEvent, logs.EventDBScan,
				logs.KeyHotelName, hotelName,
				"scan_target", "Room",
				logs.KeyError, err)
			return nil, custom_errors.ErrRepoScanFailed
		}
		rooms = append(rooms, rm)
	}

	if err := rows.Err(); err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBQuery,
			logs.KeyHotelName, hotelName,
			"detail", "rows iteration failure",
			logs.KeyError, err)
		return nil, custom_errors.ErrRepoQueryFailed
	}

	return rooms, nil
}

func (r *RoomRepositoryImpl) GetCost(hotelName string, number int) (float32, error) {
	var cost float32

	query := "SELECT cost FROM rooms WHERE hotel_name = $1 AND number = $2"
	row := r.DB.QueryRow(query, hotelName, number)

	err := row.Scan(&cost)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, custom_errors.ErrEntityNotFound
	}
	if err != nil {
		r.Log.Error(logs.MsgDatabaseQueryFailed,
			logs.KeyEvent, logs.EventDBScan,
			logs.KeyHotelName, hotelName,
			slog.Int(logs.KeyRoomNumber, number),
			logs.KeyError, err)
		return 0, custom_errors.ErrRepoScanFailed
	}

	return cost, nil
}

func (r *RoomRepositoryImpl) AddNewRoom(room entity.Room) error {
	r.Log.Debug(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventDBWrite,
		logs.KeyHotelName, room.HotelName)

	query := `
        INSERT INTO rooms (hotel_name, number, cost) 
        VALUES ($1, $2, $3)`

	_, err := r.DB.Exec(query, room.HotelName, room.Number, room.Cost)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return custom_errors.ErrEntityAlreadyExists
		}
		r.Log.Error(logs.MsgDatabaseWriteFailed,
			logs.KeyEvent, logs.EventDBWrite,
			logs.KeyHotelName, room.HotelName,
			slog.Int(logs.KeyRoomNumber, room.Number),
			logs.KeyError, err)

		return custom_errors.ErrDatabaseFailure
	}

	r.Log.Debug(logs.MsgOperationSuccess, logs.KeyEvent, logs.EventDBWrite)
	return nil
}

func (r *RoomRepositoryImpl) UpdateCost(hotelName string, number int, newCost float32) error {
	r.Log.Debug(logs.MsgStartOperation,
		logs.KeyEvent, logs.EventDBWrite,
		logs.KeyHotelName, hotelName)

	query := `
        UPDATE rooms 
        SET cost = $3 
        WHERE hotel_name = $1 AND number = $2`

	result, err := r.DB.Exec(query, hotelName, number, newCost)

	if err != nil {
		r.Log.Error(logs.MsgDatabaseWriteFailed,
			logs.KeyEvent, logs.EventDBWrite,
			logs.KeyHotelName, hotelName,
			slog.Int(logs.KeyRoomNumber, number),
			logs.KeyError, err)
		return custom_errors.ErrDatabaseFailure
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.Log.Error(logs.MsgDatabaseWriteFailed,
			logs.KeyEvent, logs.EventDBWrite,
			logs.KeyHotelName, hotelName,
			slog.Int(logs.KeyRoomNumber, number),
			logs.KeyError, err)
		return custom_errors.ErrDatabaseFailure
	}

	if rowsAffected == 0 {
		r.Log.Warn(logs.MsgEntityNotFound,
			logs.KeyHotelName, hotelName,
			slog.Int(logs.KeyRoomNumber, number))
		return custom_errors.ErrEntityNotFound
	}

	r.Log.Debug(logs.MsgOperationSuccess, logs.KeyEvent, logs.EventDBWrite)
	return nil
}
