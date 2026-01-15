package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/lib/role"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type AuthHandler struct {
	authService service.Auth
	userService service.User
}

func NewAuthHandler(authService service.Auth, userService service.User) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req *dto.UserLogin
	if err := ctx.BindJSON(&req); err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToParse, err, http.StatusBadRequest)
		return
	}

	token, err := h.authService.AuthUser(ctx, req)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgAuthenticationFailed, err, http.StatusUnauthorized)
		return
	}
	setCookie(ctx, token.Access, token.Refresh)

	user, err := h.userService.Read(ctx, token.UserID)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *AuthHandler) GetUser(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgAuthenticationFailed, err, http.StatusUnauthorized)
		return
	}

	user, err := h.userService.Read(ctx, id)
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgFailedToGet, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UserIdentity(ctx *gin.Context) {
	access, err := ctx.Cookie("access")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		logger.ResponseErr(ctx, logger.MsgAuthenticationFailed, err, http.StatusUnauthorized)
		return
	}

	refresh, err := ctx.Cookie("refresh")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		logger.ResponseErr(ctx, logger.MsgAuthenticationFailed, err, http.StatusUnauthorized)
		return
	}

	token, err := h.authService.CheckToken(ctx, &jwt_auth.Token{
		Access:  access,
		Refresh: refresh,
	})
	if err != nil {
		logger.ResponseErr(ctx, logger.MsgAuthenticationFailed, err, http.StatusUnauthorized)
		return
	}

	if token.Access != access || token.Refresh != refresh {
		setCookie(ctx, token.Access, token.Refresh)
	}

	ctx.Set("userId", token.UserID)
	ctx.Set("userRole", token.UserRole)
}

func (h *AuthHandler) RootAccess(ctx *gin.Context) {
	if ok, err := checkRole(ctx, role.RootRole); err != nil || !ok {
		logger.ResponseErr(ctx, logger.MsgAccessDenied, err, http.StatusForbidden)
		return
	}
}

func (h *AuthHandler) AdminAccess(ctx *gin.Context) {
	if ok, err := checkRole(ctx, role.AdminRole); err != nil || !ok {
		logger.ResponseErr(ctx, logger.MsgAccessDenied, err, http.StatusForbidden)
		return
	}
}

func (h *AuthHandler) GoverningAccess(ctx *gin.Context) {
	if ok, err := checkRole(ctx, role.GoverningRole); err != nil || !ok {
		logger.ResponseErr(ctx, logger.MsgAccessDenied, err, http.StatusForbidden)
		return
	}
}

func (h *AuthHandler) EmployeeAccess(ctx *gin.Context) {
	if ok, err := checkRole(ctx, role.EmployeeRole); err != nil || !ok {
		logger.ResponseErr(ctx, logger.MsgAccessDenied, err, http.StatusForbidden)
		return
	}
}

func getUserId(ctx *gin.Context) (int64, error) {
	userId, ok := ctx.Get("userId")
	if !ok {
		return 0, logger.Error(logger.MsgFailedToGet, logger.ErrUserIdNotFound)
	}

	return userId.(int64), nil
}

func checkRole(ctx *gin.Context, access role.Role) (bool, error) {
	if userRole, ok := ctx.Get("userRole"); !ok {
		return false, logger.Error(logger.MsgFailedToGet, logger.ErrUserRoleNotFound)
	} else {
		return userRole.(role.Role).CanAccess(access), nil
	}
}

func setCookie(ctx *gin.Context, access, refresh string) {
	ctx.SetCookie("access", access, 3600, "/", "", true, true)
	ctx.SetCookie("refresh", refresh, 604800, "/", "", true, true)
}
