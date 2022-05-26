package handler

import (
	"documentation/input"
	"documentation/service"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type featureHandler struct {
	featureService service.FeatureService
}

func NewFeatureHandler(featureService service.FeatureService) *featureHandler {
	return &featureHandler{featureService}
}

func (h *featureHandler) Index(c *gin.Context) {
	session := sessions.Default(c)

	data := session.Get("userName")
	flash := session.Get("message")
	token := session.Get("token")
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "List Of Features"})
	c.HTML(http.StatusOK, "list_features", gin.H{"data": flash, "token": token})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *featureHandler) New(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("userName")
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Create a new feature"})
	c.HTML(http.StatusOK, "feature_create.html", nil)
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *featureHandler) Create(c *gin.Context) {
	var input input.FeatureInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusOK, "feature_create.html", input)
		return
	}

	_, err = h.featureService.FeatureServiceCreate(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Create Feature Success")
	session.Save()
	c.Redirect(http.StatusFound, "/features")

}

func (h *featureHandler) Update(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	inputID := input.InputIDFeature{ID: id}
	featured, err := h.featureService.FeatureServiceGetByID(inputID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "feature_edit.html", nil)
		return
	}
	input := input.FeatureInput{}
	input.FeatureName = featured.FeatureName
	input.FeatureDescription = featured.FeatureDescription

	session := sessions.Default(c)
	data := session.Get("userName")
	token := session.Get("token")
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Update a Feature"})
	c.HTML(http.StatusOK, "feature_edit.html", gin.H{"input": input, "token": token, "ID": id})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *featureHandler) UpdateAction(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputData input.FeatureInput

	var inputID input.InputIDFeature
	inputID.ID = id

	inputData.FeatureName = c.PostForm("feature_name")
	inputData.FeatureDescription = c.PostForm("feature_description")
	session := sessions.Default(c)
	data := session.Get("userName")

	_, err := h.featureService.FeatureServiceUpdate(inputID, inputData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "header", gin.H{"nama": data, "title": "Update Feature"})
		c.HTML(http.StatusInternalServerError, "feature_edit.html", inputData)
		c.HTML(http.StatusInternalServerError, "footer", nil)
		return
	}
	session.Set("message", "Update Feature Success")
	session.Save()
	c.Redirect(http.StatusFound, "/features")
}

func (h *featureHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputID input.InputIDFeature
	inputID.ID = id
	_, err := h.featureService.FeatureServiceGetByID(inputID)
	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "Feature not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/features")
		return
	}
	_, err = h.featureService.FeatureServiceDeleteByID(inputID)

	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "Feature not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/features")
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Delete Feature Success")
	session.Save()
	c.Redirect(http.StatusFound, "/features")

}
