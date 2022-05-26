package handler

import (
	"documentation/formatter"
	"documentation/helper"
	"documentation/input"
	"documentation/service"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type prdHandler struct {
	service service.PrdService
}

func NewPrdHandler(service service.PrdService) *prdHandler {
	return &prdHandler{service}
}
func (h *prdHandler) GetPrd(c *gin.Context) {
	var input input.InputIDPrd
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to get Prd", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	prdDetail, err := h.service.PrdServiceGetByID(input)
	if err != nil {
		response := helper.ApiResponse("Failed to get Prd", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Detail Prd", http.StatusOK, "success", formatter.FormatPrd(prdDetail))
	c.JSON(http.StatusOK, response)
}

func (h *prdHandler) GetPrds(c *gin.Context) {
	prds, err := h.service.PrdServiceGetAll()
	if err != nil {
		response := helper.ApiResponse("Failed to get Prds", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("List of Prds", http.StatusOK, "success", formatter.FormatPrds(prds))
	c.JSON(http.StatusOK, response)
}

func (h *prdHandler) CreatePrd(c *gin.Context) {
	var input input.PrdInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Create Prd failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newPrd, err := h.service.PrdServiceCreate(input)
	if err != nil {
		response := helper.ApiResponse("Create Prd failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Create Prd", http.StatusOK, "success", formatter.FormatPrd(newPrd))
	c.JSON(http.StatusOK, response)
}
func (h *prdHandler) UpdatePrd(c *gin.Context) {
	var inputID input.InputIDPrd
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Prds", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var inputData input.PrdInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Update Prd failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updatedPrd, err := h.service.PrdServiceUpdate(inputID, inputData)
	if err != nil {
		response := helper.ApiResponse("Failed to get Prds", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Update Prd", http.StatusOK, "success", formatter.FormatPrd(updatedPrd))
	c.JSON(http.StatusOK, response)
}
func (h *prdHandler) DeletePrd(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	var inputID input.InputIDPrd
	inputID.ID = id
	_, err := h.service.PrdServiceGetByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Prds", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.service.PrdServiceDeleteByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Prds", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Delete Prd", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *prdHandler) GetPRDSR(c *gin.Context) {
	search := c.Query("search")
	page := c.Query("page")
	size := c.Query("size")

	convertedPage, _ := strconv.Atoi(page)
	convertedSize, _ := strconv.Atoi(size)

	featureData, err := h.service.FindAll(search, convertedPage, convertedSize)
	if err != nil {
		errorMessage := gin.H{"error_message": err.Error()}
		response := helper.ApiResponse("Gagal mendapatkan data", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	totalData, err := h.service.TotalFetchData(search, convertedPage, convertedSize)
	if err != nil {
		errorMessage := gin.H{"error_message": err.Error()}
		response := helper.ApiResponse("Gagal mendapatkan data", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	totalPage := math.Ceil(float64(totalData) / float64(convertedSize))
	currentPage := convertedPage + 1

	response := helper.ServerSideResponses(totalData, int(totalPage), currentPage, formatter.FormatPrds(featureData))
	c.JSON(http.StatusOK, response)

}
