package apiresponse

type Response struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func OK(data any, message string) Response {
	return Response{Success: true, Message: message, Data: data}
}

func OKWithMeta(data any, message string, meta Meta) Response {
	return Response{Success: true, Message: message, Data: data, Meta: &meta}
}

func Fail(code, message string, details any) Response {
	return Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message, Details: details},
	}
}

func NotFound(message string) Response {
	return Fail("NOT_FOUND", message, nil)
}

func Unauthorized(message string) Response {
	return Fail("UNAUTHORIZED", message, nil)
}

func BadRequest(message string, details any) Response {
	return Fail("BAD_REQUEST", message, details)
}

func InternalError(message string) Response {
	return Fail("INTERNAL_ERROR", message, nil)
}
