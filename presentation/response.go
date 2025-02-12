package presentation

type ResponseBase struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (ResponseBase) Failed(status int, message string) ResponseBase {
	return ResponseBase{
		Status:  status,
		Message: message,
	}
}

func (ResponseBase) Success(data interface{}, message string, status int) ResponseBase {
	return ResponseBase{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
