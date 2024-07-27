package handlers

import (
	"auth-service/service"
)

type HTTPHandler struct {
	US *service.UserService
}

func NewHandler(us *service.UserService) *HTTPHandler {
	return &HTTPHandler{US: us}
}
