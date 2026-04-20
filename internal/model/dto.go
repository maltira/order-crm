package model

type CreateUserRequest struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	Fio       string `json:"fio"`
	IDRole    int    `json:"id_role"`
	IsBlocked bool   `json:"is_blocked"`
}

type UpdateUserRequest struct {
	ID        int     `json:"id"`
	Login     *string `json:"login"`
	Fio       *string `json:"fio"`
	IDRole    *int    `json:"id_role"`
	IsBlocked *bool   `json:"is_blocked"`
}

type ChangePasswordRequest struct {
	ID          int    `json:"id"`
	NewPassword string `json:"new_password"`
}
