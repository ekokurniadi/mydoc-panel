package handler

import (
	"documentation/input"
	"documentation/service"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type prdHandler struct {
	prdService service.PrdService
}

func NewPRDHandler(prdService service.PrdService) *prdHandler {
	return &prdHandler{prdService}
}

func (h *prdHandler) Index(c *gin.Context) {
	session := sessions.Default(c)

	data := session.Get("userName")
	flash := session.Get("message")
	token := session.Get("token")
	session.Set("message", "")
	session.Save()
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "List Of PRD"})
	c.HTML(http.StatusOK, "list_prd", gin.H{"data": flash, "token": token})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *prdHandler) New(c *gin.Context) {
	session := sessions.Default(c)
	data := session.Get("userName")
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Create a new PRD"})
	c.HTML(http.StatusOK, "prd_create.html", nil)
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *prdHandler) Create(c *gin.Context) {
	var input input.PrdInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusOK, "prd_create.html", input)
		return
	}

	_, err = h.prdService.PrdServiceCreate(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Create PRD Success")
	session.Save()
	c.Redirect(http.StatusFound, "/prds")

}

func (h *prdHandler) Update(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)
	inputID := input.InputIDPrd{ID: id}
	featured, err := h.prdService.PrdServiceGetByID(inputID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "prd_edit.html", nil)
		return
	}
	input := input.PrdInput{}
	input.DocumentName = featured.DocumentName
	input.Description = featured.Description
	input.Link = featured.Link

	session := sessions.Default(c)
	data := session.Get("userName")
	token := session.Get("token")
	c.HTML(http.StatusOK, "header", gin.H{"nama": data, "title": "Update a PRD"})
	c.HTML(http.StatusOK, "prd_edit.html", gin.H{"input": input, "token": token, "ID": id})
	c.HTML(http.StatusOK, "footer", nil)
}

func (h *prdHandler) UpdateAction(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputData input.PrdInput

	var inputID input.InputIDPrd
	inputID.ID = id

	inputData.DocumentName = c.PostForm("document_name")
	inputData.Description = c.PostForm("description")
	inputData.Link = c.PostForm("link")
	session := sessions.Default(c)
	data := session.Get("userName")

	_, err := h.prdService.PrdServiceUpdate(inputID, inputData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "header", gin.H{"nama": data, "title": "Update PRD"})
		c.HTML(http.StatusInternalServerError, "prd_edit.html", inputData)
		c.HTML(http.StatusInternalServerError, "footer", nil)
		return
	}
	session.Set("message", "Update PRD Success")
	session.Save()
	c.Redirect(http.StatusFound, "/prds")
}

func (h *prdHandler) Delete(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.Atoi(param)

	var inputID input.InputIDPrd
	inputID.ID = id
	_, err := h.prdService.PrdServiceGetByID(inputID)
	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "PRD not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/prds")
		return
	}
	_, err = h.prdService.PrdServiceDeleteByID(inputID)

	if err != nil {
		session := sessions.Default(c)
		session.Set("message", "PRD not Found")
		session.Save()
		c.Redirect(http.StatusFound, "/prds")
		return
	}
	session := sessions.Default(c)
	session.Set("message", "Delete PRD Success")
	session.Save()
	c.Redirect(http.StatusFound, "/prds")

}
