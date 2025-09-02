package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oatsmoke/warehouse_backend/internal/lib/jwt_auth"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
	"github.com/oatsmoke/warehouse_backend/internal/service"
)

type AuthHandler struct {
	AuthService     service.Auth
	EmployeeService service.Employee
}

func NewAuthHandler(authService service.Auth, employeeService service.Employee) *AuthHandler {
	return &AuthHandler{
		AuthService:     authService,
		EmployeeService: employeeService,
	}
}

// SignIn is user authentication for login
func (h *AuthHandler) SignIn(ctx *gin.Context) {
	req := make(map[string]string)
	if err := ctx.BindJSON(&req); err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.AuthUser(ctx, req["login"], req["password"])
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	//user, err := user.GenerateToken(id)
	//if err != nil {
	//	logger.ErrResponse(ctx, err, http.StatusInternalServerError)
	//	return
	//}
	//
	//hash, err := h.AuthService.GenerateHash(ctx, id)
	//if err != nil {
	//	logger.ErrResponse(ctx, err, http.StatusInternalServerError)
	//	return
	//}

	employee, err := h.EmployeeService.GetById(ctx, user.UserID)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	setCookie(ctx, user.Access, user.Refresh)

	logger.InfoInConsole(fmt.Sprintf("%s authenticated", employee.Name))
	ctx.JSON(http.StatusOK, employee)
}

// GetUser is getting user data
func (h *AuthHandler) GetUser(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
		return
	}

	employee, err := h.EmployeeService.GetById(ctx, id)
	if err != nil {
		logger.ErrResponse(ctx, err, http.StatusBadRequest)
		return
	}

	logger.InfoInConsole(fmt.Sprintf("%s authenticated", employee.Name))
	ctx.JSON(http.StatusOK, employee)
}

// UserIdentity is user authentication check
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

	token, err := h.AuthService.Check(ctx, &jwt_auth.Token{
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

	//userId, err := token.ParseToken(token)
	//if err != nil {
	//	logger.ErrInConsole(err)
	//	h.Refresh(ctx)
	//	return
	//}

	ctx.Set("userId", token.UserID)
	logger.InfoInConsole(fmt.Sprintf("id set to %b", token.UserID))
}

// Refresh is a user authentication update
//func (h *AuthHandler) Refresh(ctx *gin.Context) {
//	hash, err := ctx.Cookie("hash")
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
//		return
//	}
//
//	userId, err := h.AuthService.FindByHash(ctx, hash)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusUnauthorized)
//		return
//	}
//
//	newToken, err := jwt_auth.GenerateToken(userId)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	newHash, err := h.AuthService.GenerateHash(ctx, userId)
//	if err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	ctx.Set("token", newToken)
//	ctx.Set("hash", newHash)
//	ctx.Set("userId", userId)
//
//	if err := setCookie(ctx); err != nil {
//		logger.ErrResponse(ctx, err, http.StatusInternalServerError)
//		return
//	}
//
//	logger.InfoInConsole(fmt.Sprintf("token and hash for id %b updated", userId))
//}

// getUserId is to get the user id from the context
func getUserId(ctx *gin.Context) (int64, error) {
	if userId, ok := ctx.Get("userId"); !ok {
		return 0, errors.New("user id in context not found")
	} else {
		return userId.(int64), nil
	}
}

// setCookie is to set cookies
func setCookie(ctx *gin.Context, access, refresh string) {
	//if access, ok := ctx.Get("access"); ok {
	ctx.SetCookie("access", access, 3600, "/", "", true, true)
	//} else {
	//	return errors.New("access token in context not found")
	//}

	//if refresh, ok := ctx.Get("refresh"); ok {
	ctx.SetCookie("refresh", refresh, 604800, "/", "", true, true)
	//} else {
	//	return errors.New("refresh token in context not found")
	//}

	//return nil
}
