package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type CompanyHandler struct {
	CompanyService service.Company
}

func NewCompanyHandler(companyService service.Company) *CompanyHandler {
	return &CompanyHandler{
		CompanyService: companyService,
	}
}

// Create is company create
func (h *CompanyHandler) Create(ctx *gin.Context) {
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CompanyService.Create(ctx, company.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", company.Title))
	ctx.JSON(http.StatusOK, "")
}

// Update is company update
func (h *CompanyHandler) Update(ctx *gin.Context) {
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CompanyService.Update(ctx, company.ID, company.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", company.Title))
	ctx.JSON(http.StatusOK, "")
}

// Delete is company delete
func (h *CompanyHandler) Delete(ctx *gin.Context) {
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CompanyService.Delete(ctx, company.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", company.ID))
	ctx.JSON(http.StatusOK, "")
}

// Restore is a company restore
func (h *CompanyHandler) Restore(ctx *gin.Context) {
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CompanyService.Restore(ctx, company.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", company.ID))
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all companies
func (h *CompanyHandler) GetAll(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.CompanyService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("companies list sended (deleted = %t)", deleted))
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get company by id
func (h *CompanyHandler) GetById(ctx *gin.Context) {
	var company *model.Company
	if err := ctx.BindJSON(&company); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.CompanyService.GetById(ctx, company.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("company %s found", res.Title))
	ctx.JSON(http.StatusOK, res)
}
