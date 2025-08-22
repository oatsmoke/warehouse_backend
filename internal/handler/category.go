package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type CategoryHandler struct {
	CategoryService service.Category
}

func NewCategoryHandler(categoryService service.Category) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: categoryService,
	}
}

// Create is category create
func (h *CategoryHandler) Create(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Create(ctx, category.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", category.Title))
	ctx.JSON(http.StatusOK, "")
}

// Update is category update
func (h *CategoryHandler) Update(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Update(ctx, category.ID, category.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", category.Title))
	ctx.JSON(http.StatusOK, "")
}

// Delete is category delete
func (h *CategoryHandler) Delete(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Delete(ctx, category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", category.ID))
	ctx.JSON(http.StatusOK, "")
}

// Restore is category restore
func (h *CategoryHandler) Restore(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Restore(ctx, category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", category.ID))
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all categories
func (h *CategoryHandler) GetAll(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.CategoryService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("categories list sended (deleted = %t)", deleted))
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get category by id
func (h *CategoryHandler) GetById(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.CategoryService.GetById(ctx, category.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("category %s found", res.Title))
	ctx.JSON(http.StatusOK, res)
}
