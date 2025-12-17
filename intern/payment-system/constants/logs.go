package constants

// log keys
const (
	KeyError      = "error"
	KeyWebHook    = "webhook"
	KeyAmount     = "amount"
	KeyMessage    = "message"
	KeyPaymentKey = "key"
)

// Error Messages
const (
	MsgFailedConnectionDB     = "Cannot connect to database"
	MsgFailedShutdown         = "Error during shutdown system"
	MsgFailedLinkCreation     = "Error during create payment link"
	MsgFailedToCreateResponse = "Error during creating response"
	MsgFailedRenderPage       = "Error during render payment page"
	MsgFailedConfirm          = "Error during confirm payment"
	MsgFailedUpdateRecord     = "Error during update record"
	MsgFailedSendRequest      = "Error during send request to webhook"
)

// Info message
const (
	MsgStartSystem    = "Start payment system"
	MsgShutdownSystem = "Shutdown payment system"
	MsgCreateLink     = "Create payment link request"
	MsgRenderPage     = "Render payment page request"
	MsgConfirm        = "Confirm payment request"
	MsgRecordExpired  = "The payment has already been confirmed or the time has expired"
)
