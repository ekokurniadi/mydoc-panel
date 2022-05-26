package formatter

import "documentation/entity"

type UserLoginFormatter struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func FormatUserLogin(user entity.User, token string) UserLoginFormatter {
	formatter := UserLoginFormatter{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Token:    token,
	}

	return formatter
}
