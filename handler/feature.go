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

type featureHandler struct {
	service service.FeatureService
}

func NewFeatureHandler(service service.FeatureService) *featureHandler {
	return &featureHandler{service}
}
func (h *featureHandler) GetFeature(c *gin.Context) {
	var input input.InputIDFeature
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to get Feature", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	featureDetail, err := h.service.FeatureServiceGetByID(input)
	if err != nil {
		response := helper.ApiResponse("Failed to get Feature", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Detail Feature", http.StatusOK, "success", formatter.FormatFeature(featureDetail))
	c.JSON(http.StatusOK, response)
}

func (h *featureHandler) GetFeatures(c *gin.Context) {
	features, err := h.service.FeatureServiceGetAll()
	if err != nil {
		response := helper.ApiResponse("Failed to get Features", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("List of Features", http.StatusOK, "success", formatter.FormatFeatures(features))
	c.JSON(http.StatusOK, response)
}

func (h *featureHandler) CreateFeature(c *gin.Context) {
	var input input.FeatureInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Create Feature failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newFeature, err := h.service.FeatureServiceCreate(input)
	if err != nil {
		response := helper.ApiResponse("Create Feature failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Create Feature", http.StatusOK, "success", formatter.FormatFeature(newFeature))
	c.JSON(http.StatusOK, response)
}
func (h *featureHandler) UpdateFeature(c *gin.Context) {
	var inputID input.InputIDFeature
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Features", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var inputData input.FeatureInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Update Feature failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updatedFeature, err := h.service.FeatureServiceUpdate(inputID, inputData)
	if err != nil {
		response := helper.ApiResponse("Failed to get Features", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Update Feature", http.StatusOK, "success", formatter.FormatFeature(updatedFeature))
	c.JSON(http.StatusOK, response)
}
func (h *featureHandler) DeleteFeature(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	var inputID input.InputIDFeature
	inputID.ID = id
	_, err := h.service.FeatureServiceGetByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Features", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.service.FeatureServiceDeleteByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get Features", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Delete Feature", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *featureHandler) GetFeaturesSSR(c *gin.Context) {
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

	response := helper.ServerSideResponses(totalData, int(totalPage), currentPage, formatter.FormatFeatures(featureData))
	c.JSON(http.StatusOK, response)

}
