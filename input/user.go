package input

type InputIDUser struct {
	ID int `uri:"id" binding:"required"`
}

type UserInput struct {
	Username string `json:"username" form:"username" `
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name" `
	Role     string `json:"role" form:"role"`
}
type UserWebInput struct {
	ID       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username" `
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name" `
	Role     string `json:"role" form:"role"`
}

//Generated by Micagen at 25 Mei 2022