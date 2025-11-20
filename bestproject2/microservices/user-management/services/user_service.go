package services

import (
	"fmt"
	"user-management/models"
	"user-management/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(id string) (*models.User, error)
	UpdateProfile(id string, req *models.UpdateProfileRequest) (*models.User, error)
	ChangePassword(id string, req *models.ChangePasswordRequest) error
	DeleteAccount(id string) error
	ListUsers(limit, offset int) ([]*models.User, error)
	UpdateUserRole(id string, role string) error
	GetStats() (*models.UserStats, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) UpdateProfile(id string, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return user, nil
}

func (s *userService) ChangePassword(id string, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *userService) DeleteAccount(id string) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ListUsers(limit, offset int) ([]*models.User, error) {
	return s.userRepo.List(limit, offset)
}

func (s *userService) UpdateUserRole(id string, role string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	user.Role = role
	return s.userRepo.Update(user)
}

func (s *userService) GetStats() (*models.UserStats, error) {
	return s.userRepo.GetStats()
}
