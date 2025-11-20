package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"user-management/config"
	"user-management/models"
	"user-management/repository"
	"user-management/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(refreshToken string) (*models.LoginResponse, error)
	ForgotPassword(email string) (string, error)
	ResetPassword(token, newPassword string) error
	ValidateToken(tokenString string) (*utils.Claims, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.User, error) {
	if existingUser, _ := s.userRepo.GetByEmail(req.Email); existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}
	
	if existingUser, _ := s.userRepo.GetByUsername(req.Username); existingUser != nil {
		return nil, fmt.Errorf("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         "user",
		IsActive:     true,
		IsVerified:   false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmailOrUsername(req.EmailOrUsername)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	s.userRepo.UpdateLastLogin(user.ID)

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.config.JWT.RefreshExpiry)),
	}

	if err := s.userRepo.CreateRefreshToken(refreshTokenModel); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.config.JWT.AccessExpiry,
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	tokenModel, err := s.userRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if tokenModel.RevokedAt != nil {
		return nil, fmt.Errorf("refresh token has been revoked")
	}

	if time.Now().After(tokenModel.ExpiresAt) {
		return nil, fmt.Errorf("refresh token has expired")
	}

	user, err := s.userRepo.GetByID(tokenModel.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	s.userRepo.RevokeRefreshToken(refreshToken)

	newTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.config.JWT.RefreshExpiry)),
	}
	s.userRepo.CreateRefreshToken(newTokenModel)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.config.JWT.AccessExpiry,
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

func (s *authService) ForgotPassword(email string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.userRepo.CreatePasswordResetToken(resetToken); err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	return token, nil
}

func (s *authService) ResetPassword(token, newPassword string) error {
	resetToken, err := s.userRepo.GetPasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid reset token")
	}

	if resetToken.UsedAt != nil {
		return fmt.Errorf("reset token has already been used")
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return fmt.Errorf("reset token has expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.GetByID(resetToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.userRepo.MarkPasswordResetTokenUsed(token)

	return nil
}

func (s *authService) ValidateToken(tokenString string) (*utils.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*utils.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *authService) generateAccessToken(user *models.User) (string, error) {
	claims := &utils.Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.config.JWT.AccessExpiry))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-management-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *authService) generateRefreshToken(user *models.User) (string, error) {
	return generateRandomToken(32)
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
