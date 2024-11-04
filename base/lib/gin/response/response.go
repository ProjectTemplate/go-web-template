package response

// Response is the response struct
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace-id"`
}

func newResponse(data interface{}, reason Reason) Response {
	return Response{
		Code:    reason.Code,
		Message: reason.Message,
		Data:    data,
	}
}
