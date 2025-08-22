package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type ProfileHandler struct {
	ProfileService service.Profile
}

func NewProfileHandler(profileService service.Profile) *ProfileHandler {
	return &ProfileHandler{
		ProfileService: profileService,
	}
}

// Create is profile create
func (h *ProfileHandler) Create(ctx *gin.Context) {
	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ProfileService.Create(ctx, profile.Title, profile.Category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s created", profile.Title))
	ctx.JSON(http.StatusOK, "")
}

// Update is profile update
func (h *ProfileHandler) Update(ctx *gin.Context) {
	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ProfileService.Update(ctx, profile.ID, profile.Title, profile.Category.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s updated", profile.Title))
	ctx.JSON(http.StatusOK, "")
}

// Delete is profile delete
func (h *ProfileHandler) Delete(ctx *gin.Context) {
	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ProfileService.Delete(ctx, profile.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d deleted", profile.ID))
	ctx.JSON(http.StatusOK, "")
}

// Restore is profile restore
func (h *ProfileHandler) Restore(ctx *gin.Context) {
	var profile *model.Category
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := h.ProfileService.Restore(ctx, profile.ID); err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%d restored", profile.ID))
	ctx.JSON(http.StatusOK, "")
}

// GetAll is to get all profiles
func (h *ProfileHandler) GetAll(ctx *gin.Context) {
	var deleted bool
	if err := ctx.BindJSON(&deleted); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.ProfileService.GetAll(ctx, deleted)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("profiles list sended (deleted = %t)", deleted))
	ctx.JSON(http.StatusOK, res)
}

// GetById is to get profile by id
func (h *ProfileHandler) GetById(ctx *gin.Context) {
	var profile *model.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	res, err := h.ProfileService.GetById(ctx, profile.ID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("profile %s found", res.Title))
	ctx.JSON(http.StatusOK, res)
}
