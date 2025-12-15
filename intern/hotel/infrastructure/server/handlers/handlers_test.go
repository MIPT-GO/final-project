package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"final-project/intern/hotel/infrastructure/server/dto"
	"final-project/pkg/custom_errors"
	"final-project/tests/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetAllHotels(t *testing.T) {
	cases := []struct {
		name         string
		getAllFunc   func() ([]string, error)
		expectedCode int
	}{
		{"Success", func() ([]string, error) { return []string{"HotelA"}, nil }, http.StatusOK},
		{"DBError", func() ([]string, error) { return nil, errors.New("db fail") }, http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := NewHotelHandler(slog.Default(), &mock.MockHotelService{GetAllFunc: c.getAllFunc}, &mock.MockRoomService{})
			req := httptest.NewRequest(http.MethodGet, "/v1/hotel/all", nil)
			w := httptest.NewRecorder()

			handler.(*HotelHandlerImpl).GetAllHotels(w, req)

			if w.Code != c.expectedCode {
				t.Fatalf("expected %d, got %d", c.expectedCode, w.Code)
			}
		})
	}
}

func TestCreateHotel(t *testing.T) {
	cases := []struct {
		name         string
		body         dto.CreateHotelRequest
		createFunc   func(name, email string) error
		expectedCode int
	}{
		{"Success", dto.CreateHotelRequest{Name: "HotelA", Email: "a@example.com"}, func(name, email string) error { return nil }, http.StatusCreated},
		{"Conflict", dto.CreateHotelRequest{Name: "HotelA", Email: "a@example.com"}, func(name, email string) error { return custom_errors.ErrEntityAlreadyExists }, http.StatusConflict},
		{"DBError", dto.CreateHotelRequest{Name: "HotelA", Email: "a@example.com"}, func(name, email string) error { return custom_errors.ErrDatabaseFailure }, http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := NewHotelHandler(slog.Default(), &mock.MockHotelService{CreateFunc: c.createFunc}, &mock.MockRoomService{})
			bodyJSON, _ := json.Marshal(c.body)
			req := httptest.NewRequest(http.MethodPost, "/v1/hotel", bytes.NewReader(bodyJSON))
			w := httptest.NewRecorder()

			handler.(*HotelHandlerImpl).CreateHotel(w, req)

			if w.Code != c.expectedCode {
				t.Fatalf("expected %d, got %d", c.expectedCode, w.Code)
			}
		})
	}
}

func TestGetOneHotel(t *testing.T) {
	cases := []struct {
		name         string
		getEmailFunc func(name string) (string, error)
		expectedCode int
	}{
		{"Success", func(name string) (string, error) { return "a@example.com", nil }, http.StatusOK},
		{"NotFound", func(name string) (string, error) { return "", custom_errors.ErrEntityNotFound }, http.StatusNotFound},
		{"DBError", func(name string) (string, error) { return "", custom_errors.ErrDatabaseFailure }, http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := NewHotelHandler(slog.Default(), &mock.MockHotelService{GetEmailFunc: c.getEmailFunc}, &mock.MockRoomService{})
			req := httptest.NewRequest(http.MethodGet, "/v1/hotel/HotelA", nil)
			req = mux.SetURLVars(req, map[string]string{"hotel": "HotelA"})
			w := httptest.NewRecorder()

			handler.(*HotelHandlerImpl).GetOneHotel(w, req)

			if w.Code != c.expectedCode {
				t.Fatalf("expected %d, got %d", c.expectedCode, w.Code)
			}
		})
	}
}

func TestGetAllRooms(t *testing.T) {
	cases := []struct {
		name         string
		getAllFunc   func(name string) ([]int, error)
		expectedCode int
	}{
		{"Success", func(name string) ([]int, error) { return []int{101, 102}, nil }, http.StatusOK},
		{"DBError", func(name string) ([]int, error) { return nil, errors.New("db fail") }, http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := NewHotelHandler(slog.Default(), &mock.MockHotelService{}, &mock.MockRoomService{GetAllFunc: c.getAllFunc})
			req := httptest.NewRequest(http.MethodGet, "/v1/hotel/HotelA/rooms", nil)
			req = mux.SetURLVars(req, map[string]string{"hotel": "HotelA"})
			w := httptest.NewRecorder()

			handler.(*HotelHandlerImpl).GetAllRooms(w, req)

			if w.Code != c.expectedCode {
				t.Fatalf("expected %d, got %d", c.expectedCode, w.Code)
			}
		})
	}
}

func TestUpdateRoom(t *testing.T) {
	cases := []struct {
		name           string
		reqBody        dto.UpdateRoomRequest
		updateCostFunc func(hotelName string, number int, newCost float32) error
		expectedCode   int
	}{
		{"Success", dto.UpdateRoomRequest{Number: 101, NewCost: 150.5}, func(hotelName string, number int, newCost float32) error { return nil }, http.StatusOK},
		{"NotFound", dto.UpdateRoomRequest{Number: 101, NewCost: 150.5}, func(hotelName string, number int, newCost float32) error {
			return custom_errors.ErrEntityNotFound
		}, http.StatusNotFound},
		{"DBError", dto.UpdateRoomRequest{Number: 101, NewCost: 150.5}, func(hotelName string, number int, newCost float32) error {
			return custom_errors.ErrDatabaseFailure
		}, http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := NewHotelHandler(slog.Default(), &mock.MockHotelService{}, &mock.MockRoomService{UpdateCostFunc: c.updateCostFunc})
			bodyJSON, _ := json.Marshal(c.reqBody)
			req := httptest.NewRequest(http.MethodPut, "/v1/hotel/HotelA/rooms", bytes.NewReader(bodyJSON))
			req = mux.SetURLVars(req, map[string]string{"hotel": "HotelA"})
			w := httptest.NewRecorder()

			handler.(*HotelHandlerImpl).UpdateRoom(w, req)

			if w.Code != c.expectedCode {
				t.Fatalf("expected %d, got %d", c.expectedCode, w.Code)
			}
		})
	}
}

func TestGetRoom(t *testing.T) {
	cases := []struct {
		name         string
		numberStr    string
		getCostFunc  func(name string, number int) (float32, error)
		expectedCode int
	}{
		{"Success", "101", func(name string, number int) (float32, error) { return 200.5, nil }, http.StatusOK},
		{"InvalidNumber", "abc", nil, http.StatusBadRequest},
		{"NotFound", "101", func(name string, number int) (float32, error) { return 0, custom_errors.ErrEntityNotFound }, http.StatusNotFound},
		{"DBError", "101", func(name string, number int) (float32, error) { return 0, custom_errors.ErrDatabaseFailure }, http.StatusInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := NewHotelHandler(slog.Default(), &mock.MockHotelService{}, &mock.MockRoomService{GetCostFunc: c.getCostFunc})
			req := httptest.NewRequest(http.MethodGet, "/v1/hotel/HotelA/rooms/"+c.numberStr, nil)
			req = mux.SetURLVars(req, map[string]string{"hotel": "HotelA", "number": c.numberStr})
			w := httptest.NewRecorder()

			handler.(*HotelHandlerImpl).GetRoom(w, req)

			if w.Code != c.expectedCode {
				t.Fatalf("expected %d, got %d", c.expectedCode, w.Code)
			}
		})
	}
}
