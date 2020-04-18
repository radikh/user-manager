package server

import "github.com/lvl484/user-manager/model"

type responseCreateAccount struct {
	ID        string `json:"id"`
	Username  string `json:"user_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type responseAccountInfo struct {
	Username  string `json:"user_name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

func convertToResponseCreateAccount(u *model.User) responseCreateAccount {
	return responseCreateAccount{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
	}
}

func convertToResponseAccountInfo(u *model.User) responseAccountInfo {
	return responseAccountInfo{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Email:     u.Email,
	}
}
