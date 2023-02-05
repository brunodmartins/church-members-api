package apierrors

type Error struct {
	message      string
	statusCode   int
	customFields map[string]interface{}
}

// NewApiError creates a custom error
func NewApiError(message string, statusCode int) Error {
	return Error{
		message:    message,
		statusCode: statusCode,
	}
}

func (error Error) Error() string {
	return error.message
}

func (error Error) StatusCode() int {
	return error.statusCode
}

func (error *Error) AddField(key string, value interface{}) {
	if error.customFields == nil {
		error.customFields = make(map[string]interface{})
	}
	error.customFields[key] = value
}

func (error Error) GetField(key string) interface{} {
	return error.customFields[key]
}
