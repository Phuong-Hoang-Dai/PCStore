package handler

func BuildResponse(status string, errorMsg string, data any) map[string]any {
	return map[string]any{
		"status":     status,
		"error_code": errorMsg,
		"data":       data,
	}
}
