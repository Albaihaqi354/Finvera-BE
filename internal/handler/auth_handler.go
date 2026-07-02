package handler

import (
	"net/http"

	"finvera-be/internal/config"
	"finvera-be/internal/dto"
	"finvera-be/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService, cfg}
}

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Email    string `json:"email" binding:"required,email,max=255" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=8,max=128" example:"StrongP@ss123"`
}

// LoginRequest is the request body for login.
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"StrongP@ss123"`
}

// Register godoc
// @Summary      Registrasi user baru
// @Description  Membuat akun user baru dengan username, email, dan password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Data registrasi"
// @Success      201   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("User registered successfully", gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	}))
}

// Login godoc
// @Summary      Login user
// @Description  Autentikasi user dan mendapatkan JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Kredensial login"
// @Success      200   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Failure      401   {object}  dto.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	token, err := h.authService.Login(req.Username, req.Password, h.cfg)
	if err != nil {
		// Always return 401 for invalid credentials — don't expose whether user exists
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("invalid credentials"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Login successful", gin.H{
		"token": token,
	}))
}

// Logout godoc
// @Summary      Logout user
// @Description  Client-side logout — instructs client to discard the JWT token
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  dto.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT is stateless. Actual token invalidation must be handled on the client
	// by removing the token from storage (localStorage / cookie).
	// TODO: Implement server-side token blacklist with Redis for enhanced security.
	c.JSON(http.StatusOK, dto.SuccessResponse("Logout successful. Please discard your token.", nil))
}
