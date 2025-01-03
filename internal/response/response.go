package response

type Response struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

func (r *Response) SetSuccess(success bool) *Response {
	r.Success = success
	return r
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

func (r *Response) SetData(key string, value any) *Response {
	r.Data[key] = value
	return r
}
