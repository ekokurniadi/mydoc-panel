package handler

import (
	"documentation/input"
	"documentation/service"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type featureDetailHandler struct {
	featureDetailService service.FeatureDetailService
	featureService       service.FeatureService
}

func NewFeatureDetailHandler(featureDetailService service.FeatureDetailService, featureService service.FeatureService) *featureDetailHandler {
	return &featureDetailHandler{featureDetailService, featureService}
}

func (h *featureDetailHandler) Index(c *gin.Context) {
	session := sessions.Default(c)

	data := session.Get("userName")
	flash := session.Get("message")
	token := session.Get("token")
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "List Of Features Detail"})
	c.HTML(http.StatusOK, "list_features_detail", gin.H{"data": flash, "token": token})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *featureDetailHandler) New(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("userName")
	token := session.Get("token")
	features, _ := h.featureService.FeatureServiceGetAll()
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Create a new feature detail"})
	c.HTML(http.StatusOK, "feature_detail_create.html", gin.H{"token": token, "features": features})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *featureDetailHandler) Create(c *gin.Context) {
	var input input.FeatureDetailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.HTML(http.StatusOK, "feature_detail_create.html", input)
		return
	}

	_, err = h.featureDetailService.FeatureDetailServiceCreate(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Create Feature Detail Success")
	session.Save()
	c.Redirect(http.StatusFound, "/features_detail")

}

func (h *featureDetailHandler) Update(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	inputID := input.InputIDFeatureDetail{ID: id}
	featured, err := h.featureDetailService.FeatureDetailServiceGetByID(inputID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "feature_edit.html", nil)
		return
	}
	var inputHeader input.InputIDFeature
	inputHeader.ID = featured.FeatureID
	features, _ := h.featureService.FeatureServiceGetByID(inputHeader)

	input := input.FeatureDetailInput{}
	input.FeatureID = featured.FeatureID
	input.Title = featured.Title
	input.PathOfFile = featured.PathOfFile
	input.Code = featured.Code
	input.AuthorName = featured.AuthorName
	input.Description = featured.Description
	listOffeatures, _ := h.featureService.FeatureServiceGetAll()
	session := sessions.Default(c)
	data := session.Get("userName")
	token := session.Get("token")
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Update a Feature Detail"})
	c.HTML(http.StatusOK, "feature_detail_edit.html", gin.H{"input": input, "token": token, "ID": id, "featureName": features.FeatureName, "features": listOffeatures})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *featureDetailHandler) UpdateAction(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputData input.FeatureDetailInput

	var inputID input.InputIDFeatureDetail
	inputID.ID = id

	inputData.FeatureID, _ = strconv.Atoi(c.PostForm("feature_id"))
	inputData.Title = c.PostForm("title")
	inputData.PathOfFile = c.PostForm("path_of_file")
	inputData.Code = c.PostForm("code")
	inputData.AuthorName = c.PostForm("author_name")
	inputData.Description = c.PostForm("description")
	session := sessions.Default(c)
	data := session.Get("userName")

	_, err := h.featureDetailService.FeatureDetailServiceUpdate(inputID, inputData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "header", gin.H{"nama": data, "title": "Update Feature Detail"})
		c.HTML(http.StatusInternalServerError, "feature_detail_edit.html", inputData)
		c.HTML(http.StatusInternalServerError, "footer", nil)
		return
	}
	session.Set("message", "Update Feature Detail Success")
	session.Save()
	c.Redirect(http.StatusFound, "/features_detail")
}

func (h *featureDetailHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputID input.InputIDFeatureDetail
	inputID.ID = id
	_, err := h.featureDetailService.FeatureDetailServiceGetByID(inputID)
	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "Feature Detail not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/features_detail")
		return
	}
	_, err = h.featureDetailService.FeatureDetailServiceDeleteByID(inputID)

	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "Feature Detail not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/features_detail")
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Delete Feature Detail Success")
	session.Save()
	c.Redirect(http.StatusFound, "/features_detail")

}
