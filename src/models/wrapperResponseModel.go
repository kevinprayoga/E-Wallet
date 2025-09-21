package models

type Response struct {
	Code       int         `json:"code"`
	Content    interface{} `json:"content"`
	TotalItems int         `json:"totalItems"`
	Message    string      `json:"message"`
}
