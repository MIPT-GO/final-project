package constants

const (
	EnvPostgresHost        = "DB_HOST"
	EnvPostgresPort        = "DB_PORT"
	EnvPostgresUser        = "DB_USER"
	EnvPostgresPassword    = "DB_PASSWORD"
	EnvPostgresDB          = "DB_DATABASE"
	EnvPostgresMaxOpen     = "DB_MAX_OPEN"
	EnvPostgresMaxIdle     = "DB_MAX_IDLE"
	EnvPostgresConnTimeout = "DB_CONN_TIMEOUT"
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
	EnvServerHost     = "SERVER_HOST"
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
	QueryInsertReservation = "INSERT INTO " + TableReservations + " (email, start, \"end\", hotel, number) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	QueryFindById          = "SELECT id, email, start, \"end\", hotel, number FROM " + TableReservations + " WHERE id=$1"
	QueryCountOverlapping  = "SELECT COUNT(1) FROM " + TableReservations + " WHERE hotel=$1 AND number=$2 AND NOT ($3 >= \"end\" OR $4 <= start)"
	QueryDeleteReservation = "DELETE FROM " + TableReservations + " WHERE email=$1 AND hotel=$2 AND number=$3 AND start=$4 AND \"end\"=$5"
)

// Hotel service endpoints (format strings)
const (
	HotelRoomCheckEndpoint      = "/v1/hotel/%s/room/%d"
	HotelAvailableRoomsEndpoint = "/v1/hotel/%s/rooms/"
	HotelRoomPriceEndpoint      = "/v1/hotel/%s/room/%d"
	HotelOwnerEmailEndpoint     = "/v1/hotel/%s"
	PaymentInitiateEndpoint     = "/v1/payment/link"
)

// Payment service env keys and defaults
const (
	EnvPaymentServiceHost     = "PAYMENT_SERVICE_HOST"
	EnvPaymentServicePort     = "PAYMENT_SERVICE_PORT"
	EnvPaymentServiceTimeout  = "PAYMENT_SERVICE_TIMEOUT"
	DefaultPaymentServicePort = "8081"
)

// Kafka env keys and defaults
const (
	EnvKafkaHost      = "KAFKA_HOST"
	EnvKafkaTopic     = "KAFKA_TOPIC"
	DefaultKafkaHost  = "localhost:9092"
	DefaultKafkaTopic = "bookings"
)

// Messages
const (
	MsgHotelServiceStatusFmt = "hotel service returned status %d"
)
