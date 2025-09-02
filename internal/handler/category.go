// Package handler provides HTTP handlers for category operations.
package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

// CategoryHandler handles HTTP requests for category operations.
type CategoryHandler struct {
	CategoryService service.Category
}

// NewCategoryHandler creates a new CategoryHandler with the given CategoryService.
// categoryService: Category service.
// Returns a pointer to CategoryHandler.
func NewCategoryHandler(categoryService service.Category) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: categoryService,
	}
}

// Create creates a new category.
// ctx: Gin context.
// No return value.
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

	logger.InfoInConsole(fmt.Sprintf("category %s created", category.Title))
	ctx.JSON(http.StatusOK, "")
}

// Update updates an existing category by id.
// ctx: Gin context.
// No return value.
func (h *CategoryHandler) Update(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Update(ctx, id, category.Title); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("category %d updated", id))
	ctx.JSON(http.StatusOK, "")
}

// Delete performs a soft delete of a category by id.
// ctx: Gin context.
// No return value.
func (h *CategoryHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Delete(ctx, id); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("category %d deleted", id))
	ctx.JSON(http.StatusOK, "")
}

// Restore restores a previously deleted category by id.
// ctx: Gin context.
// No return value.
func (h *CategoryHandler) Restore(ctx *gin.Context) {
	var category *model.Category
	if err := ctx.BindJSON(&category); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Restore(ctx, id); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("category %d restored", id))
	ctx.JSON(http.StatusOK, "")
}

// GetAll retrieves all categories, filtered by deleted status.
// ctx: Gin context.
// No return value.
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

// GetById retrieves a category by id.
// ctx: Gin context.
// No return value.
func (h *CategoryHandler) GetById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.CategoryService.GetById(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("category %d found", id))
	ctx.JSON(http.StatusOK, res)
}
