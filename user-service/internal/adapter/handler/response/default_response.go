package response

type DefaultReponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
