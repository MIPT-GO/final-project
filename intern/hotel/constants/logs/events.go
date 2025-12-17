package logs

const (
	EventDBQuery = "db_query"
	EventDBWrite = "db_write"
	EventDBScan  = "db_scan_fail"
)

const (
	EventServiceInit      = "service_init"
	EventServerShutdown   = "server_shutdown"
	EventRouterSetup      = "router_setup"
	EventRouteRegister    = "route_register"
	EventPermissionDenied = "auth_denied"
	EventDBConnect        = "db-connect"
	EventDBMigration      = "db-migration"
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
