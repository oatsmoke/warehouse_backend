package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/dto"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
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
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	token, err := h.authService.AuthUser(ctx, req)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}
	setCookie(ctx, token.Access, token.Refresh)

	user, err := h.userService.Read(ctx, token.UserID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *AuthHandler) GetUser(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	user, err := h.userService.Read(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UserIdentity(ctx *gin.Context) {
	access, err := ctx.Cookie("access")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		logger.ErrResponse(ctx, err, http.StatusForbidden)
		return
	}

	refresh, err := ctx.Cookie("refresh")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		logger.ErrResponse(ctx, err, http.StatusForbidden)
		return
	}

	token, err := h.authService.CheckToken(ctx, &jwt_auth.Token{
		Access:  access,
		Refresh: refresh,
	})
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusForbidden)
		return
	}

	if token.Access != access || token.Refresh != refresh {
		setCookie(ctx, access, refresh)
	}

	ctx.Set("userId", token.UserID)
}

func getUserId(ctx *gin.Context) (int64, error) {
	if userId, ok := ctx.Get("userId"); !ok {
		return 0, errors.New("user id in context not found")
	} else {
		return userId.(int64), nil
	}
}

func setCookie(ctx *gin.Context, access, refresh string) {
	ctx.SetCookie("access", access, 3600, "/", "", true, true)
	ctx.SetCookie("refresh", refresh, 604800, "/", "", true, true)
}
