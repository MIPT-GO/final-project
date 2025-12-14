package constants

// Log keys
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
	KeyMethod     = "method"
	KeyPath       = "path"
	KeyHandler    = "handler"
	KeyCount      = "count"
	KeyStatus     = "status"
	KeyAddr       = "addr"
	KeyQuery      = "query"
	KeyDSN        = "dsn"
	KeyURL        = "url"
	KeyStatusCode = "status_code"
)

// Log events
const (
	EventIncomingRequest     = "incoming_request"
	EventRequestCompleted    = "request_completed"
	EventGetByEmail          = "get_by_email"
	EventGetByHotel          = "get_by_hotel"
	EventGetAvailableNumbers = "get_available_numbers"
	EventBookRoomInHotel     = "book_room_in_hotel"
	EventCheckAccuracy       = "check_accuracy"
	EventServerStarted       = "server_started"
	EventServerShutdown      = "server_shutdown"
	EventFailedToEncode      = "failed_to_encode"
	EventFailedToFetch       = "failed_to_fetch"
	EventMethodNotAllowed    = "method_not_allowed_event"
	EventRepoInit            = "repo_init"
	EventDBQuery             = "db_query"
	EventDBInsert            = "db_insert"
	EventDBPing              = "db_ping"
	EventHasOverlap          = "has_overlap"
	EventHotelRequest        = "hotel_request"
	EventHotelResponse       = "hotel_response"
)

// Messages
const (
	MsgMethodNotAllowed    = "method not allowed"
	MsgBadRequest          = "bad request"
	MsgInternalServerError = "internal server error"
	MsgUnprocessableEntity = "unprocessable entity"
	MsgServerStarted       = "server started"
	MsgServerShutdown      = "graceful shutdown"
	MsgConfigLoadFailedFmt = "failed to load config: %v"
	MsgFailedOpenDB        = "failed to open db"
	MsgDBPingFailed        = "db ping failed"
)

// Service names
const (
	ServiceBooking = "booking_service"
)
