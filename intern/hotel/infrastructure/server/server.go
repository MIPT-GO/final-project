package server

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"

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

func NewServer(log *slog.Logger, db *sql.DB, addr string, readTimeout, writeTimeout, idleTimeout time.Duration) *HotelServer {
	hRepo := repository.NewHotelRepository(db, log)
	rRepo := repository.NewRoomRepository(db, log)

	hService := hotel_service.NewHotelService(log, hRepo)
	rService := room_service.NewRoomService(log, hRepo, rRepo)

	handlerImpl := handlers.NewHotelHandler(log, hService, rService)

	router := setupRoutes(handlerImpl, log)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return &HotelServer{
		Log:        log,
		HTTPServer: httpServer,
		DB:         db,
	}
}

func setupRoutes(h interfaces.HotelHandler, log *slog.Logger) *mux.Router {
	log.Info(logs.MsgStartOperation, logs.KeyEvent, logs.EventRouterSetup)

	router := mux.NewRouter()

	logRoute := func(method, path string) {
		log.Debug(logs.MsgRouteRegistered,
			logs.KeyEvent, logs.EventRouteRegister,
			slog.String("method", method),
			slog.String("path", path))
	}

	router.HandleFunc("/v1/hotel/all/", h.GetAllHotels).Methods("GET")
	logRoute("GET", "/v1/hotel/all/")

	router.HandleFunc("/v1/hotel/", h.CreateHotel).Methods("POST")
	logRoute("POST", "/v1/hotel/")

	router.HandleFunc("/v1/hotel/{hotel}/", h.GetOneHotel).Methods("GET")
	logRoute("GET", "/v1/hotel/{hotel}/")

	router.HandleFunc("/v1/hotel/{hotel}/rooms/", h.GetAllRooms).Methods("GET")
	logRoute("GET", "/v1/hotel/{hotel}/rooms/")

	router.HandleFunc("/v1/hotel/{hotel}/room/", h.UpdateRoom).Methods("PUT")
	logRoute("PUT", "/v1/hotel/{hotel}/room/")

	router.HandleFunc("/v1/hotel/{hotel}/room/{number}", h.GetRoom).Methods("GET")
	logRoute("GET", "/v1/hotel/{hotel}/room/{number}")

	router.HandleFunc("/v1/hotel/{hotel}/room/", h.CreateRoom).Methods("POST")
	logRoute("POST", "/v1/hotel/{hotel}/room/")

	log.Info(logs.MsgOperationSuccess, logs.KeyEvent, logs.EventRouterSetup)
	return router
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
