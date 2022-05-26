package handler

import (
	"documentation/auth"
	"documentation/entity"
	"documentation/formatter"
	"documentation/helper"
	"documentation/input"
	"documentation/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service     service.UserService
	authService auth.Service
}

func NewUserHandler(service service.UserService, authService auth.Service) *userHandler {
	return &userHandler{service, authService}
}
func (h *userHandler) GetUser(c *gin.Context) {
	var input input.InputIDUser
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to get User", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	userDetail, err := h.service.UserServiceGetByID(input)
	if err != nil {
		response := helper.ApiResponse("Failed to get User", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Detail User", http.StatusOK, "success", formatter.FormatUser(userDetail))
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) GetUsers(c *gin.Context) {
	users, err := h.service.UserServiceGetAll()
	if err != nil {
		response := helper.ApiResponse("Failed to get Users", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("List of Users", http.StatusOK, "success", formatter.FormatUsers(users))
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var input input.UserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Create User failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newUser, err := h.service.UserServiceCreate(input)
	if err != nil {
		response := helper.ApiResponse("Create User failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Create User", http.StatusOK, "success", formatter.FormatUser(newUser))
	c.JSON(http.StatusOK, response)
}
func (h *userHandler) UpdateUser(c *gin.Context) {
	var inputID input.InputIDUser
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Users", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var inputData input.UserInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Update User failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updatedUser, err := h.service.UserServiceUpdate(inputID, inputData)
	if err != nil {
		response := helper.ApiResponse("Failed to get Users", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Update User", http.StatusOK, "success", formatter.FormatUser(updatedUser))
	c.JSON(http.StatusOK, response)
}
func (h *userHandler) DeleteUser(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	var inputID input.InputIDUser
	inputID.ID = id
	_, err := h.service.UserServiceGetByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Users", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.service.UserServiceDeleteByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Users", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Delete User", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input input.LoginInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response :=
			helper.ApiResponse("Login Account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedInUser, err := h.service.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response :=
			helper.ApiResponse("Login Account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)

	if err != nil {
		response :=
			helper.ApiResponse("Login Account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := formatter.FormatUserLogin(loggedInUser, token)
	response := helper.ApiResponse("Login success", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(entity.User)

	formatter := formatter.FormatUser(currentUser)
	response := helper.ApiResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}
