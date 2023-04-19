package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/pkg/model"
)

func (h *Handler) createCompany(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var company model.Company
	if err := c.BindJSON(&company); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Company.Create(company.Title); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) getByIdCompany(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var d model.Company
	if err := c.BindJSON(&d); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	company, err := h.service.Company.GetById(d.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, company)
}

func (h *Handler) getAllCompany(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	companies, err := h.service.Company.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, companies)
}

func (h *Handler) updateCompany(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var company model.Company
	if err := c.BindJSON(&company); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Company.Update(company.Id, company.Title); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}

func (h *Handler) deleteCompany(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var company model.Company
	if err := c.BindJSON(&company); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Company.Delete(company.Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(c)
	c.JSON(http.StatusOK, "")
}
