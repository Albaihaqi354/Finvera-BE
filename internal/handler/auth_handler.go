package handler

import (
	"net/http"
	"strings"
	"time"

	"finvera-be/internal/config"
	"finvera-be/internal/dto"
	"finvera-be/internal/service"
	"finvera-be/pkg/blacklist"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService, cfg}
}


// Register godoc
// @Summary      Registrasi user baru
// @Description  Membuat akun user baru dengan username, email, dan password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Data registrasi"
// @Success      201   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
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
// @Param        body  body      dto.LoginRequest  true  "Kredensial login"
// @Success      200   {object}  dto.Response
// @Failure      400   {object}  dto.Response
// @Failure      401   {object}  dto.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			
			// Parse without validation just to extract expiration claim
			token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			var expTime time.Time
			if token != nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if exp, ok := claims["exp"].(float64); ok {
						expTime = time.Unix(int64(exp), 0)
					}
				}
			}
			if expTime.IsZero() {
				expTime = time.Now().Add(24 * time.Hour)
			}
			
			blacklist.Add(tokenString, expTime)
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Logout successful", nil))
}
