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

////////////////////////////////////////////////////////
// REGISTER
////////////////////////////////////////////////////////

func (s *authService) Register(req *models.RegisterRequest) (*models.User, error) {

	// check email exists
	if existing, _ := s.userRepo.GetByEmail(req.Email); existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// check username exists
	if existing, _ := s.userRepo.GetByUsername(req.Username); existing != nil {
		return nil, fmt.Errorf("username already taken")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// create user struct (FIXED LastName)
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

////////////////////////////////////////////////////////
// LOGIN (FULLY FIXED)
////////////////////////////////////////////////////////

func (s *authService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {

	// unified lookup (email or username)
	user, err := s.userRepo.GetByEmailOrUsername(req.EmailOrUsername)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is disabled")
	}

	// compare password
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// update last_login_at
	_ = s.userRepo.UpdateLastLogin(user.ID)

	// create access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// create refresh token
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// save refresh token
	rt := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.config.JWT.RefreshExpiry)),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreateRefreshToken(rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// return login response
	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.config.JWT.AccessExpiry,
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

////////////////////////////////////////////////////////
// REFRESH TOKEN
////////////////////////////////////////////////////////

func (s *authService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {

	tokenModel, err := s.userRepo.GetRefreshToken(refreshToken)
	if err != nil || tokenModel == nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if tokenModel.RevokedAt != nil {
		return nil, fmt.Errorf("refresh token revoked")
	}

	if time.Now().After(tokenModel.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := s.userRepo.GetByID(tokenModel.UserID)
	if err != nil || user == nil {
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

	// revoke old
	_ = s.userRepo.RevokeRefreshToken(refreshToken)

	// save new token
	newRT := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Second * time.Duration(s.config.JWT.RefreshExpiry)),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreateRefreshToken(newRT); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.config.JWT.AccessExpiry,
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

////////////////////////////////////////////////////////
// FORGOT PASSWORD
////////////////////////////////////////////////////////

func (s *authService) ForgotPassword(email string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil || user == nil {
		return "", nil // do not reveal existence
	}

	token, err := generateRandomToken(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	reset := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreatePasswordResetToken(reset); err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	return token, nil
}

////////////////////////////////////////////////////////
// RESET PASSWORD
////////////////////////////////////////////////////////

func (s *authService) ResetPassword(token, newPassword string) error {

	prt, err := s.userRepo.GetPasswordResetToken(token)
	if err != nil || prt == nil {
		return fmt.Errorf("invalid reset token")
	}

	if prt.UsedAt != nil {
		return fmt.Errorf("reset token already used")
	}

	if time.Now().After(prt.ExpiresAt) {
		return fmt.Errorf("reset token expired")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.GetByID(prt.UserID)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}

	user.PasswordHash = string(newHash)

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	_ = s.userRepo.MarkPasswordResetTokenUsed(token)

	return nil
}

////////////////////////////////////////////////////////
// VALIDATE TOKEN
////////////////////////////////////////////////////////

func (s *authService) ValidateToken(tokenString string) (*utils.Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&utils.Claims{},
		func(token *jwt.Token) (interface{}, error) {
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

////////////////////////////////////////////////////////
// TOKEN HELPERS
////////////////////////////////////////////////////////

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
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
