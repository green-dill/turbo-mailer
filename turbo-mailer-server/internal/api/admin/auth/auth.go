package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// @Summary		User login
// @Description	Authenticate a user and return a JWT token
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		Credentials	true	"User credentials"
// @Success		200			{object}	map[string]string
// @Router			/api/v1/auth/login [post]
func Login(ctx echo.Context) error {
	var credentials Credentials
	if err := ctx.Bind(&credentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate credentials (replace with your actual authentication logic)
	// For example, check against database
	user, err := validateUser(ctx.Request().Context(), credentials.Username, credentials.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return err
	}

	// Set JWT token as a cookie
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = t
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = strings.EqualFold(ctx.Request().URL.Scheme, "https")
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

// @Summary		User logout
// @Description	Invalidate the user's JWT token
// @Tags			Auth
// @Produce		json
// @Success		200	{object}	map[string]string
// @Router			/api/v1/auth/logout [post]
func Logout(ctx echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = strings.EqualFold(ctx.Request().URL.Scheme, "https")
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// @Summary		Change user password
// @Description	Change the authenticated user's password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		ChangePasswordRequest	true	"Change password request"
// @Success		200		{object}	map[string]string
// @Security		JWT
// @Router			/api/v1/auth/change-password [post]
func ChangePassword(ctx echo.Context) error {
	var req ChangePasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := ctx.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get the current user from the JWT token
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	// Validate the old password
	dbUser, err := validateUser(ctx.Request().Context(), username, req.OldPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid old password")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash new password")
	}

	// Update the password in the database
	dbUser.Password = string(hashedPassword)
	if err := query.User.WithContext(ctx.Request().Context()).Save(dbUser); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}

type ProfileResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// @Summary		Get user profile
// @Description	Retrieve the profile of the currently authenticated user
// @Tags			Auth
// @Produce		json
// @Success		200	{object}	ProfileResponse
// @Security		JWT
// @Router			/api/v1/auth/profile [get]
func Profile(ctx echo.Context) error {
	// 使用 echojwt 获取用户信息
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	// Retrieve user information from the database
	dbUser, err := query.User.WithContext(ctx.Request().Context()).Where(query.User.Username.Eq(username)).First()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user information")
	}

	// Prepare the response
	response := ProfileResponse{
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, response)
}

// validateUser is an internal function and doesn't need Swagger documentation
func validateUser(ctx context.Context, username, password string) (*models.User, error) {
	user, err := query.User.WithContext(ctx).Where(query.User.Username.Eq(username)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
