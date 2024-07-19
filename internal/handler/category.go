package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse_backend/internal/lib/logger"
	"warehouse_backend/internal/model"
	"warehouse_backend/internal/service"
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
	const fn = "handler.Category.Create"

	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.CategoryService.Create(ctx, category.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", category.Title), fn)
	ctx.JSON(http.StatusOK, "")
}

// Update is category update
func (h *CategoryHandler) Update(ctx *gin.Context) {
	const fn = "handler.Category.Update"

	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.CategoryService.Update(ctx, category.ID, category.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", category.Title), fn)
	ctx.JSON(http.StatusOK, "")
}

// Delete is category delete
func (h *CategoryHandler) Delete(ctx *gin.Context) {
	const fn = "handler.Category.Delete"

	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.CategoryService.Delete(ctx, category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", category.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// Restore is category restore
func (h *CategoryHandler) Restore(ctx *gin.Context) {
	const fn = "handler.Category.Restore"

	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	if err := h.CategoryService.Restore(ctx, category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", category.ID), fn)
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all categories
func (h *CategoryHandler) GetAll(ctx *gin.Context) {
	const fn = "handler.Category.GetAll"

	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.CategoryService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("categories list sended (deleted = %t)", deleted), fn)
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get category by id
func (h *CategoryHandler) GetById(ctx *gin.Context) {
	const fn = "handler.Category.GetById"

	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest, fn)
		return
	}

	res, err := h.CategoryService.GetById(ctx, category.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError, fn)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("category %s found", res.Title), fn)
	ctx.JSON(http.StatusOK, res)
}
