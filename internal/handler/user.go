package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/model"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type UserHandler struct {
	userService service.User
}

func NewUserHandler(userService service.User) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Create(ctx *gin.Context) {
	var req *dto.UserRequest
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if req != nil && !req.Role.IsValid() {
		logger.ResponseErr(ctx, logger.MsgFailedToValidate, logger.ErrInvalidRole, http.StatusBadRequest)
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Employee: &model.Employee{
			ID: req.EmployeeID,
		},
	}

	err := h.userService.Create(ctx, user)
	if err != nil {
		if errors.Is(err, logger.ErrAlreadyExists) {
			logger.ResponseErr(ctx, logger.ErrAlreadyExists.Error(), err, http.StatusConflict)
			return
		}
		logger.ResponseErr(ctx, logger.MsgFailedToInsert, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, "")
}

func (h *UserHandler) Read(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	res, err := h.userService.Read(ctx, id)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	var req *dto.UserRequest
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	user := &model.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Employee: &model.Employee{
			ID: req.EmployeeID,
		},
	}

	err = h.userService.Update(ctx, user)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	if err := h.userService.Delete(ctx, id); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToDelete, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *UserHandler) List(ctx *gin.Context) {
	res, err := h.userService.List(ctx)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) SetPassword(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	var req *dto.UserPasswordUpdate
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	err = h.userService.SetPassword(ctx, id, req.OldPassword, req.NewPassword)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	err = h.userService.ResetPassword(ctx, id)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h *UserHandler) SetEnabled(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	var req *dto.UserEnabledUpdate
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	err = h.userService.SetEnabled(ctx, id, req.Enabled)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToUpdate, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, "")
}
