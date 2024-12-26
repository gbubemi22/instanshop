package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"instashop/internal/model"
	"instashop/internal/service"
)

type UserController struct {
	UserService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{UserService: userService}
}

// CreateUser handles the creation of a new user.
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.UserService.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// CreateAdminUser handles the creation of a new user.
func (ctrl *UserController) CreateAdmin(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := ctrl.UserService.CreateAdmin(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

type VerifyEmailRequest struct {
	Email    string `json:"email" binding:"required"`
	OtpToken string `json:"otpToken" binding:"required"`
}

// VerifyEmailHandler handles the verification of a user's email
func (ctrl *UserController) VerifyEmailHandler(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and OTP token are required"})
		return
	}
	ctrl.UserService.VerifyEmail(req.Email, req.OtpToken)

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

type SendEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func (ctrl *UserController) SendEmailHandler(c *gin.Context) {
	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	ctrl.UserService.SendMail(req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})

}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ctrl *UserController) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := ctrl.UserService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) 
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
