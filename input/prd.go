package input

type InputIDPrd struct {
	ID int `uri:"id" binding:"required"`
}

type PrdInput struct {
	DocumentName string `json:"document_name" form:"document_name"`
	Description  string `json:"description" form:"description"`
	Link         string `json:"link" form:"link"`
}

//Generated by Micagen at 25 Mei 2022
