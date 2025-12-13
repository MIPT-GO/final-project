package repository

import (
	"database/sql"
	"errors"
	"log/slog"

	"final-project/intern/hotel/domain/entity"
	"final-project/pkg/custom_errors"
	"final-project/pkg/logs"
)

type RoomRepositoryImpl struct {
	DB  *sql.DB
	Log *slog.Logger
}

func (r *RoomRepositoryImpl) GetAll(hotelName string) ([]entity.Room, error) {
	query := "SELECT hotel_name, number, cost FROM Rooms WHERE hotel_name = $1"
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

	query := "SELECT cost FROM Rooms WHERE hotel_name = $1 AND number = $2"
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
