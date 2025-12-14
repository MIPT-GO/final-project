package logs

const (
	EventDBQuery = "db_query"
	EventDBWrite = "db_write"
	EventDBScan  = "db_scan_fail"
) //for repository

const (
	EventServiceInit    = "service_init"
	EventServerShutdown = "server_shutdown"
	EventRouterSetup    = "router_setup"
	EventRouteRegister  = "route_register"
) //for server

const (
	EventPermissionDenied = "auth_denied"
) //for service
