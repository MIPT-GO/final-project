package constants

const (
	EnvPostgresHost        = "PG_HOST"
	EnvPostgresPort        = "PG_PORT"
	EnvPostgresUser        = "PG_USER"
	EnvPostgresPassword    = "PG_PASSWORD"
	EnvPostgresDB          = "PG_DATABASE"
	EnvPostgresMaxOpen     = "PG_MAX_OPEN"
	EnvPostgresMaxIdle     = "PG_MAX_IDLE"
	EnvPostgresConnTimeout = "PG_CONN_TIMEOUT"
	DefaultPostgresPort    = "5432"
)

// Hotel service defaults and env keys
const (
	EnvHotelServiceHost     = "HOTEL_SERVICE_HOST"
	EnvHotelServicePort     = "HOTEL_SERVICE_PORT"
	EnvHotelServiceTimeout  = "HOTEL_SERVICE_TIMEOUT"
	DefaultHotelServicePort = "8080"
)

// Server
const (
	EnvServerPort     = "SERVER_PORT"
	DefaultServerPort = ":8080"
)

// Logging
const (
	EnvLogLevel = "LOG_LEVEL"
)

// DB schema constants
const (
	TableReservations = "reservations"
)

// SQL queries
const (
	QueryFindByEmail       = "SELECT id, email, start, \"end\", hotel, number FROM " + TableReservations + " WHERE email=$1"
	QueryFindByHotel       = "SELECT id, email, start, \"end\", hotel, number FROM " + TableReservations + " WHERE hotel=$1"
	QueryInsertReservation = "INSERT INTO " + TableReservations + " (email, start, \"end\", hotel, number) VALUES ($1, $2, $3, $4, $5)"
	QueryCountOverlapping  = "SELECT COUNT(1) FROM " + TableReservations + " WHERE hotel=$1 AND number=$2 AND NOT ($3 >= \"end\" OR $4 <= start)"
	QueryDeleteReservation = "DELETE FROM " + TableReservations + " WHERE email=$1 AND hotel=$2 AND number=$3 AND start=$4 AND \"end\"=$5"
)

// Hotel service endpoints (format strings)
const (
	HotelRoomCheckEndpoint      = "/hotels/%s/rooms/%d"
	HotelAvailableRoomsEndpoint = "/hotels/%s/rooms/"
	HotelRoomPriceEndpoint      = "/hotels/%s/rooms/%d"
	PaymentInitiateEndpoint     = "/payments/initiate"
)

// Payment service env keys and defaults
const (
	EnvPaymentServiceHost     = "PAYMENT_SERVICE_HOST"
	EnvPaymentServicePort     = "PAYMENT_SERVICE_PORT"
	EnvPaymentServiceTimeout  = "PAYMENT_SERVICE_TIMEOUT"
	DefaultPaymentServicePort = "8081"
)

// Messages
const (
	MsgHotelServiceStatusFmt = "hotel service returned status %d"
)
