package dto

import "order-crm/internal/model"

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string     `json:"access_token"`
	User        model.User `json:"user"`
}
