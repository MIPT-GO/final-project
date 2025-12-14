package server

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"final-project/intern/hotel/application/service/hotel_service"
	"final-project/intern/hotel/application/service/room_service"
	"final-project/intern/hotel/domain/interfaces"
	"final-project/intern/hotel/infrastructure/repository"
	"final-project/intern/hotel/infrastructure/server/handlers"
	"final-project/pkg/logs"
)

type HotelServer struct {
	Log        *slog.Logger
	HTTPServer *http.Server
	DB         *sql.DB
}

func NewServer(log *slog.Logger, db *sql.DB, addr string) *HotelServer {
	hRepo := repository.NewHotelRepository(db, log)
	rRepo := repository.NewRoomRepository(db, log)

	hService := hotel_service.NewHotelService(log, hRepo)
	rService := room_service.NewRoomService(log, hRepo, rRepo)

	handlerImpl := handlers.NewHotelHandler(log, hService, rService)

	router := setupRoutes(handlerImpl, log)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &HotelServer{
		Log:        log,
		HTTPServer: httpServer,
		DB:         db,
	}
}

func setupRoutes(h interfaces.HotelHandler, log *slog.Logger) *http.ServeMux {
	log.Info(logs.MsgStartOperation, logs.KeyEvent, logs.EventRouterSetup)

	mux := http.NewServeMux()

	logRoute := func(method, path string) {
		log.Debug(logs.MsgRouteRegistered,
			logs.KeyEvent, logs.EventRouteRegister,
			slog.String("method", method),
			slog.String("path", path))
	}

	mux.HandleFunc("GET /v1/hotel/all/", h.GetAllHotels)
	logRoute("GET", "/v1/hotel/all/")

	mux.HandleFunc("POST /v1/hotel/", h.CreateHotel)
	logRoute("POST", "/v1/hotel/")

	mux.HandleFunc("GET /v1/hotel/{hotel}/", h.GetOneHotel)
	logRoute("GET", "/v1/hotel/{hotel}/")

	mux.HandleFunc("GET /v1/hotel/{hotel}/rooms/", h.GetAllRooms)
	logRoute("GET", "/v1/hotel/{hotel}/rooms/")

	mux.HandleFunc("PUT /v1/hotel/{hotel}/room/", h.UpdateRoom)
	logRoute("PUT", "/v1/hotel/{hotel}/room/")

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hotel Service is UP"))
	})
	logRoute("GET", "/health")

	log.Info(logs.MsgOperationSuccess, logs.KeyEvent, logs.EventRouterSetup)
	return mux
}

func (s *HotelServer) Run() error {
	addr := s.HTTPServer.Addr

	s.Log.Info(logs.MsgStartServer,
		logs.KeyEvent, logs.EventServiceInit,
		slog.String("addr", addr))

	err := s.HTTPServer.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Log.Error(logs.MsgUnexpectedFail,
			logs.KeyEvent, logs.EventServiceInit,
			logs.KeyError, err)
		return err
	}

	s.Log.Info(logs.MsgServerStopped,
		logs.KeyEvent, logs.EventServerShutdown,
		slog.String("addr", addr))
	return nil
}
