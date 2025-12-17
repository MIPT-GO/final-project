package constants

import "errors"

var (
	ErrHotelNotFound       = errors.New("hotel not found")
	ErrHotelOrRoomNotFound = errors.New("hotel or room number not found")
	ErrRoomAlreadyBooked   = errors.New("room already booked in this period")
	ErrRepoNotSpecified    = errors.New("repository not specified")
)

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
)
