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

type featuredetailHandler struct {
	service service.FeatureDetailService
}

func NewFeatureDetailHandler(service service.FeatureDetailService) *featuredetailHandler {
	return &featuredetailHandler{service}
}
func (h *featuredetailHandler) GetFeatureDetail(c *gin.Context) {
	var input input.InputIDFeatureDetail
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetail", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	featuredetailDetail, err := h.service.FeatureDetailServiceGetByID(input)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetail", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Detail FeatureDetail", http.StatusOK, "success", formatter.FormatFeatureDetail(featuredetailDetail))
	c.JSON(http.StatusOK, response)
}

func (h *featuredetailHandler) GetFeatureDetails(c *gin.Context) {
	var input input.InputIDFeatureDetail
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetails", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	featuredetails, err := h.service.FeatureDetailServiceGetAll(input)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetails", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("List of FeatureDetails", http.StatusOK, "success", formatter.FormatFeatureDetails(featuredetails))
	c.JSON(http.StatusOK, response)
}

func (h *featuredetailHandler) CreateFeatureDetail(c *gin.Context) {
	var input input.FeatureDetailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		print(err)

		errorMessage := gin.H{"errors": err}
		response := helper.ApiResponse("Create FeatureDetail failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newFeatureDetail, err := h.service.FeatureDetailServiceCreate(input)
	if err != nil {
		response := helper.ApiResponse("Create FeatureDetail failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Create FeatureDetail", http.StatusOK, "success", formatter.FormatFeatureDetail(newFeatureDetail))
	c.JSON(http.StatusOK, response)
}
func (h *featuredetailHandler) UpdateFeatureDetail(c *gin.Context) {
	var inputID input.InputIDFeatureDetail
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetails", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var inputData input.FeatureDetailInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Update FeatureDetail failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	updatedFeatureDetail, err := h.service.FeatureDetailServiceUpdate(inputID, inputData)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetails", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Update FeatureDetail", http.StatusOK, "success", formatter.FormatFeatureDetail(updatedFeatureDetail))
	c.JSON(http.StatusOK, response)
}
func (h *featuredetailHandler) DeleteFeatureDetail(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	var inputID input.InputIDFeatureDetail
	inputID.ID = id
	_, err := h.service.FeatureDetailServiceGetByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetails", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.service.FeatureDetailServiceDeleteByID(inputID)
	if err != nil {
		response := helper.ApiResponse("Failed to get FeatureDetails", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("Successfully Delete FeatureDetail", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *featuredetailHandler) GetFeaturesSSR(c *gin.Context) {
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

	response := helper.ServerSideResponses(totalData, int(totalPage), currentPage, formatter.FormatFeatureDetails(featureData))
	c.JSON(http.StatusOK, response)

}
