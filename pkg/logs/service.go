package logs

const (
	KeyService    = "service"
	KeyEvent      = "event"
	KeyUserID     = "user_id"
	KeyUserEmail  = "user_email"
	KeyError      = "error"
	KeyHotelEmail = "hotel_email"

	KeyHotelName  = "hotel_name"
	KeyRoomNumber = "number"
	KeyCost       = "cost"
	KeyNewCost    = "new_cost"
)

const (
	EventDatabaseQuery    = "db_query"
	EventDatabaseWrite    = "db_write"
	EventPermissionDenied = "auth_denied"
	EventRoomCostQuery    = "room_cost_query"
	EventRoomCostUpdate   = "room_cost_update"
	EventRoomCreate       = "room_create"
	EventHotelOwnerCheck  = "hotel_owner_check"
	EventHotelGetAll      = "hotel_get_all"
	EventHotelGetEmail    = "hotel_get_email"
	EventHotelCreate      = "hotel_create"
)
