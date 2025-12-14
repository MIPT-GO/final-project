package repository

import (
	"database/sql"

	"final-project/intern/booking/application/interfaces"
	"final-project/intern/booking/domain/models/reservation"
	"final-project/pkg/booking/constants"
)

type PostgresAdapter struct {
	DB     *sql.DB
	Logger interfaces.Logger
}

func NewPostgresAdapter(db *sql.DB, logger interfaces.Logger) *PostgresAdapter {
	return &PostgresAdapter{DB: db, Logger: logger}
}

func (p *PostgresAdapter) FindByEmail(email string) ([]reservation.Reserve, error) {
	if p.Logger != nil {
		p.Logger.Debug(constants.EventDBQuery, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyQuery, constants.QueryFindByEmail, constants.KeyUserEmail, email)
	}
	rows, err := p.DB.Query(constants.QueryFindByEmail, email)
	if err != nil {
		if p.Logger != nil {
			p.Logger.Error(constants.EventDBQuery, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyQuery, constants.QueryFindByEmail)
		}
		return nil, err
	}
	defer rows.Close()

	var res []reservation.Reserve
	for rows.Next() {
		var r reservation.Reserve
		if err := rows.Scan(&r.Id, &r.Email, &r.Start, &r.End, &r.Hotel, &r.Number); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (p *PostgresAdapter) FindByHotel(hotel string) ([]reservation.Reserve, error) {
	if p.Logger != nil {
		p.Logger.Debug(constants.EventDBQuery, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyQuery, constants.QueryFindByHotel, constants.KeyHotelName, hotel)
	}
	rows, err := p.DB.Query(constants.QueryFindByHotel, hotel)
	if err != nil {
		if p.Logger != nil {
			p.Logger.Error(constants.EventDBQuery, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyQuery, constants.QueryFindByHotel)
		}
		return nil, err
	}
	defer rows.Close()

	var res []reservation.Reserve
	for rows.Next() {
		var r reservation.Reserve
		if err := rows.Scan(&r.Id, &r.Email, &r.Start, &r.End, &r.Hotel, &r.Number); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (p *PostgresAdapter) AddNewReservation(reserv reservation.Reserve) error {
	if p.Logger != nil {
		p.Logger.Info(constants.EventDBInsert, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBInsert, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number, constants.KeyUserEmail, reserv.Email)
	}
	_, err := p.DB.Exec(constants.QueryInsertReservation, reserv.Email, reserv.Start, reserv.End, reserv.Hotel, reserv.Number)
	if err != nil && p.Logger != nil {
		p.Logger.Error(constants.EventDBInsert, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBInsert, constants.KeyQuery, constants.QueryInsertReservation)
	}
	return err
}

func (p *PostgresAdapter) HasOverlap(reserv reservation.Reserve) (bool, error) {
	if p.Logger != nil {
		p.Logger.Debug(constants.EventHasOverlap, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHasOverlap, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number)
	}
	var count int
	row := p.DB.QueryRow(constants.QueryCountOverlapping, reserv.Hotel, reserv.Number, reserv.Start, reserv.End)
	if err := row.Scan(&count); err != nil {
		if p.Logger != nil {
			p.Logger.Error(constants.EventHasOverlap, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHasOverlap)
		}
		return false, err
	}
	if p.Logger != nil {
		p.Logger.Debug(constants.EventHasOverlap, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHasOverlap, constants.KeyCount, count)
	}
	return count > 0, nil
}

func (p *PostgresAdapter) DeleteReservation(reserv reservation.Reserve) error {
	if p.Logger != nil {
		p.Logger.Info(constants.EventDBQuery, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyHotelName, reserv.Hotel, constants.KeyRoomNumber, reserv.Number, constants.KeyUserEmail, reserv.Email)
	}
	_, err := p.DB.Exec(constants.QueryDeleteReservation, reserv.Email, reserv.Hotel, reserv.Number, reserv.Start, reserv.End)
	if err != nil && p.Logger != nil {
		p.Logger.Error(constants.EventDBQuery, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventDBQuery, constants.KeyQuery, constants.QueryDeleteReservation)
	}
	return err
}
