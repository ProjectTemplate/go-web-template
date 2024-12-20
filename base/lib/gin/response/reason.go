package response

var (
	// ReasonSuccess 成功
	ReasonSuccess = Reason{SuccessCode, "success"}

	// AdminInternalErrorReason 内部错误
	AdminInternalErrorReason = NewReason(AdminCodeInternalError)
	// AdminParamErrorReason 参数错误
	AdminParamErrorReason = NewReason(AdminCodeParamError)
)

type Reason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewReason(code Code) Reason {
	message := codeMessageMap[code.Code]
	return Reason{
		Code:    code.Encode(),
		Message: message,
	}
}
