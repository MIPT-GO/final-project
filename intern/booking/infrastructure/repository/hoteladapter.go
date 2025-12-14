package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"final-project/intern/booking/application/interfaces"
	"final-project/pkg/booking/constants"
)

type HotelAdapter struct {
	BaseURL string
	Client  *http.Client
	Logger  interfaces.Logger
}

func NewHotelAdapter(baseURL string, client *http.Client, logger interfaces.Logger) *HotelAdapter {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &HotelAdapter{BaseURL: baseURL, Client: client, Logger: logger}
}

func (h *HotelAdapter) CheckAccuracy(hotel string, roomNumber uint64) (bool, error) {
	url := fmt.Sprintf(h.BaseURL+constants.HotelRoomCheckEndpoint, hotel, roomNumber)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if h.Logger != nil {
		h.Logger.Debug(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, url, constants.KeyHotelName, hotel, constants.KeyRoomNumber, roomNumber)
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelRequest, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, url)
		}
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode >= 400 {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelResponse, fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode), constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse, constants.KeyStatusCode, resp.StatusCode)
		}
		return false, fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode)
	}
	return true, nil
}

func (h *HotelAdapter) GetAllRoomsInHotel(hotel string) ([]uint64, error) {
	url := fmt.Sprintf(h.BaseURL+constants.HotelAvailableRoomsEndpoint, hotel)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if h.Logger != nil {
		h.Logger.Debug(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, url, constants.KeyHotelName, hotel)
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelRequest, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, url)
		}
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, constants.ErrHotelNotFound
	}
	if resp.StatusCode >= 400 {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelResponse, fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode), constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse, constants.KeyStatusCode, resp.StatusCode)
		}
		return nil, fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode)
	}
	//узнать формат комнат и нормально дешифровать
	var rooms []uint64
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelResponse, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse)
		}
		return nil, err
	}
	if h.Logger != nil {
		h.Logger.Info(constants.EventHotelResponse, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse, constants.KeyHotelName, hotel, constants.KeyCount, len(rooms))
	}
	return rooms, nil
}

func (h *HotelAdapter) GetRoomPrice(hotel string, roomNumber uint64) (string, error) {
	url := fmt.Sprintf(h.BaseURL+constants.HotelRoomPriceEndpoint, hotel, roomNumber)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	if h.Logger != nil {
		h.Logger.Debug(constants.EventHotelRequest, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, url, constants.KeyHotelName, hotel, constants.KeyRoomNumber, roomNumber)
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelRequest, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelRequest, constants.KeyURL, url)
		}
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", constants.ErrHotelNotFound
	}
	if resp.StatusCode >= 400 {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelResponse, fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode), constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse, constants.KeyStatusCode, resp.StatusCode)
		}
		return "", fmt.Errorf(constants.MsgHotelServiceStatusFmt, resp.StatusCode)
	}
	// узнать формат цены и нормально дешифровать
	var pr struct {
		Price string `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		if h.Logger != nil {
			h.Logger.Error(constants.EventHotelResponse, err, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse)
		}
		return "", err
	}
	if h.Logger != nil {
		h.Logger.Info(constants.EventHotelResponse, constants.KeyService, constants.ServiceBooking, constants.KeyEvent, constants.EventHotelResponse, constants.KeyHotelName, hotel)
	}
	return pr.Price, nil
}
