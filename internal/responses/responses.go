package responses

type ResponseBase struct {
	Result  bool        `json:"result"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}