package handlers

import "final-project/intern/booking/application/service"

type handler struct {
	service service.ReservationService
}

func NewHandler(service service.ReservationService) handler {
	return handler{service: service}
}
