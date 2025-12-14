package logs

const (
	MsgEntityNotFound        = "Entity not found"
	MsgPermissionDenied      = "Permission check failed"
	MsgDatabaseQueryFailed   = "Database query failed"
	MsgDatabaseWriteFailed   = "Database write failed"
	MsgEntityAlreadyExists   = "Entity already exists"
	MsgUpdateFailure         = "Update operation failed"
	MsgStartOperation        = "Starting operation"
	MsgOperationSuccess      = "Operation successful"
	MsgStartServer           = "Starting HTTP server"
	MsgRouteRegistered       = "Registering route"
	MsgUnexpectedFail        = "Server failed unexpectedly"
	MsgServerStopped         = "HTTP Server stopped"
	MsgDatabaseConnectFailed = "Failed to connect to database"
	MsgDatabaseCloseFailed   = "Failed to close DB connection"

	MsgDatabaseMigrationFailed   = "Failed to run database migrations"
	MsgDatabaseMigrationNoChange = "Database schema is already up to date (no migrations applied)."

	MsgShuttingDown     = "Service shutting down"
	MsgForcedShutdown   = "Server forced to shutdown"
	MsgShutdownComplete = "Shutdown complete"
)
