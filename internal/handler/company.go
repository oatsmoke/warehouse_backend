package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
)

type CompanyHandler struct {
	serviceCompany service.Company
}

func NewCompanyHandler(serviceCompany service.Company) *CompanyHandler {
	return &CompanyHandler{
		serviceCompany: serviceCompany,
	}
}

func (h *CompanyHandler) createCompany(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceCompany.Create(ctx, company.Title); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *CompanyHandler) getByIdCompany(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var d *model.Company
	if err := ctx.BindJSON(&d); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	company, err := h.serviceCompany.GetById(ctx, d.ID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, company)
}

func (h *CompanyHandler) getAllCompany(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	companies, err := h.serviceCompany.GetAll(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, companies)
}

func (h *CompanyHandler) updateCompany(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceCompany.Update(ctx, company.ID, company.Title); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}

func (h *CompanyHandler) deleteCompany(ctx *gin.Context) {
	_, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.serviceCompany.Delete(ctx, company.ID); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	setCookie(ctx)
	ctx.JSON(http.StatusOK, "")
}
