package dto

type BaseResponse struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}
