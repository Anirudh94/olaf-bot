package error

type ServiceFailureError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func NewServiceFailureError(message string) ServiceFailureError {
	return ServiceFailureError{
		"ServiceFailureError",
		message,
	}
}
