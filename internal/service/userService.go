package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"instashop/internal/common"
	"instashop/internal/model"
	"instashop/internal/utils"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) validateUserInput(user *model.User) error {
	if user.Email == "" || user.Username == "" || user.Password == "" {
		utils.NewBadRequestError("email, username, and password are required")

	}
	return nil
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	// Validate input
	if err := s.validateUserInput(user); err != nil {
		return err
	}

	// Check for existing user with the same email or username
	var existingUser model.User
	err := s.DB.WithContext(ctx).
		Where("email = ? OR username = ?", user.Email, user.Username).
		First(&existingUser).Error

	if err == nil {
		utils.NewBadRequestError("user with given email or username already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("internal server error: %w", err)
	}

	// Validate password
	passwordValidation, err := common.ValidatePasswordString(user.Password)
	if err != nil {
		return err
	}
	if !passwordValidation.IsValid {
		return errors.New("password is not valid")
	}

	// Hash password
	user.Password = utils.HashPassword(user.Password)
	user.VerifiedEmail = false

	// Generate OTP token
	otpToken, err := utils.GenerateRandomNumber()
	if err != nil {
		return fmt.Errorf("failed to generate OTP token: %w", err)
	}
	user.OtpToken = otpToken
	user.ExpiredAt = utils.GetOtpExpiryTime()

	// Set createdAt and updatedAt timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Send OTP email asynchronously
	go func() {
		to := []string{user.Email}
		subject := "Test Email"
		body := fmt.Sprintf("<h1>Hello from Mailtrap! Here is your OTP: %s</h1>", otpToken)

		if err := utils.SendMail(subject, body, to); err != nil {
			log.Printf("Could not send email: %v", err)
		} else {
			fmt.Println("Email sent successfully!")
		}
	}()

	// Insert the new user into the database
	result := s.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *UserService) CreateAdmin(ctx context.Context, user *model.User) error {
	// Validate input
	if err := s.validateUserInput(user); err != nil {
		return err
	}

	// Check for existing user with the same email or username
	var existingUser model.User
	err := s.DB.WithContext(ctx).
		Where("email = ? OR username = ?", user.Email, user.Username).
		First(&existingUser).Error

	if err == nil {
		return errors.New("user with given email or username already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("internal server error: %w", err)
	}

	// Validate password
	passwordValidation, err := common.ValidatePasswordString(user.Password)
	if err != nil {
		return err
	}
	if !passwordValidation.IsValid {
		return errors.New("password is not valid")
	}

	// Hash password
	user.Password = utils.HashPassword(user.Password)
	user.VerifiedEmail = false

	// Generate OTP token
	otpToken, err := utils.GenerateRandomNumber()
	if err != nil {
		return fmt.Errorf("failed to generate OTP token: %w", err)
	}
	user.OtpToken = otpToken
	user.ExpiredAt = utils.GetOtpExpiryTime()
	user.Role = "admin"

	// Set createdAt and updatedAt timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Send OTP email asynchronously
	go func() {
		to := []string{user.Email}
		subject := "Test Email"
		body := fmt.Sprintf("<h1>Hello from Mailtrap! Here is your OTP: %s</h1>", otpToken)

		if err := utils.SendMail(subject, body, to); err != nil {
			log.Printf("Could not send email: %v", err)
		} else {
			fmt.Println("Email sent successfully!")
		}
	}()

	// Insert the new user into the database
	result := s.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *UserService) VerifyEmail(email string, otpToken string) error {
	// Find user by email
	var user model.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("error finding user: %w", err)
	}

	// Check if OTP token matches and is not expired
	if user.OtpToken != otpToken {
		utils.NewBadRequestError("invalid OTP token")
	}

	if user.ExpiredAt.Before(time.Now()) {
		return errors.New("OTP token has expired")
	}

	// Update user record to mark email as verified
	user.VerifiedEmail = true
	user.OtpToken = ""
	user.ExpiredAt = time.Time{}

	if err := s.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserService) SendMail(email string) (string, error) {
	// Check if user with the given email exists
	var user model.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user does not exist")
		}
		return "", fmt.Errorf("error finding user: %w", err)
	}

	// Generate OTP token
	otpToken, err := utils.GenerateRandomNumber()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP token: %w", err)
	}

	// Update user's OTP token and expiry time
	user.OtpToken = otpToken
	user.ExpiredAt = utils.GetOtpExpiryTime()

	if err := s.DB.Save(&user).Error; err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	// Send email with OTP token
	to := []string{email}
	subject := "OTP for account verification"
	body := fmt.Sprintf("Your OTP is: %s", otpToken)

	if err := utils.SendMail(subject, body, to); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s\n", email)

	// Return success message
	successMessage := fmt.Sprintf("Email sent successfully to %s", email)
	return successMessage, nil
}

func (s *UserService) Login( email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Find user by email
	var user model.User
	err := s.DB.WithContext(ctx).
		Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user not found: %w", err)
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	// Check if password is correct
	if !utils.VerifyPassword(password, user.Password) {
		return "", utils.NewBadRequestError("invalid credentials")
	}

	// Check if email is verified
	if !user.VerifiedEmail {
		return "", utils.NewBadRequestError("please verify your email")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(fmt.Sprintf("%d", user.ID), user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
