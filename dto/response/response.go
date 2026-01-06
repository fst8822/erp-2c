package response

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(code int, err error) ErrResponse {
	return ErrResponse{
		Code:    code,
		Message: err.Error(),
	}
}
