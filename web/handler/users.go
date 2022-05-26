package handler

import (
	"documentation/input"
	"documentation/service"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.UserServiceGetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	data := session.Get("userName")
	flash := session.Get("message")
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "List Of Users"})
	c.HTML(http.StatusOK, "index.html", gin.H{"users": users, "data": flash})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *userHandler) New(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("userName")
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Create a new user"})
	c.HTML(http.StatusOK, "user_create.html", nil)
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *userHandler) Create(c *gin.Context) {
	var input input.UserInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusOK, "user_create.html", input)
		return
	}

	_, err = h.userService.UserServiceCreate(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Create User Success")
	session.Save()
	c.Redirect(http.StatusFound, "/users")

}

func (h *userHandler) Update(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	inputID := input.InputIDUser{ID: id}
	registeredUser, err := h.userService.UserServiceGetByID(inputID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "user_edit.html", nil)
		return
	}
	input := input.UserWebInput{}
	input.Name = registeredUser.Name
	input.ID = registeredUser.ID
	input.Username = registeredUser.Username
	input.Password = registeredUser.Password
	input.Role = registeredUser.Role

	session := sessions.Default(c)
	data := session.Get("userName")
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Create a new user"})
	c.HTML(http.StatusOK, "user_edit.html", input)
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *userHandler) UpdateAction(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputData input.UserInput

	var inputID input.InputIDUser
	inputID.ID = id

	inputData.Name = c.PostForm("name")
	inputData.Password = c.PostForm("password")
	inputData.Role = c.PostForm("role")
	session := sessions.Default(c)
	data := session.Get("userName")

	_, err := h.userService.UserServiceUpdate(inputID, inputData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "header", gin.H{"nama": data, "title": "Update User"})
		c.HTML(http.StatusInternalServerError, "user_edit.html", inputData)
		c.HTML(http.StatusInternalServerError, "footer", nil)
		return
	}
	session.Set("message", "Update User Success")
	session.Save()
	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputID input.InputIDUser
	inputID.ID = id
	_, err := h.userService.UserServiceGetByID(inputID)
	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "User not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/users")
		return
	}
	_, err = h.userService.UserServiceDeleteByID(inputID)

	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "User not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/users")
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Delete User Success")
	session.Save()
	c.Redirect(http.StatusFound, "/users")

}
