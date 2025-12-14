package logs

const (
	KeyDBName     = "Hotel"
	KeyHotelEmail = "hotel_email"
	KeyHotelName  = "hotel_name"
	KeyRoomNumber = "number"
	KeyCost       = "cost"
	KeyNewCost    = "new_cost"
	KeyUserEmail  = "email"
)

const (
	EventRoomCostQuery   = "room_cost_query"
	EventRoomCostUpdate  = "room_cost_update"
	EventRoomCreate      = "room_create"
	EventHotelOwnerCheck = "hotel_owner_check"
	EventHotelGetAll     = "hotel_get_all"
	EventHotelGetEmail   = "hotel_get_email"
	EventHotelCreate     = "hotel_create"
)
