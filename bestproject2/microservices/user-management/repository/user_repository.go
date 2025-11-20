package repository

import (
	"database/sql"
	"fmt"
	"time"
	"user-management/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmailOrUsername(emailOrUsername string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	List(limit, offset int) ([]*models.User, error)
	GetStats() (*models.UserStats, error)
	UpdateLastLogin(userID string) error

	// Refresh tokens
	CreateRefreshToken(token *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	RevokeRefreshToken(token string) error
	DeleteExpiredRefreshTokens() error

	// Password reset
	CreatePasswordResetToken(token *models.PasswordResetToken) error
	GetPasswordResetToken(token string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(token string) error

	// Audit logs
	CreateAuditLog(log *models.AuditLog) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, username, password_hash, first_name, last_name, phone, role, is_active, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.PasswordHash, user.FirstName, user.LastName, user.Phone, user.Role, user.IsActive, user.IsVerified, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, phone, role, is_active, is_verified, avatar_url, created_at, updated_at, last_login_at, deleted_at
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.Role, &user.IsActive, &user.IsVerified,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, phone, role, is_active, is_verified, avatar_url, created_at, updated_at, last_login_at, deleted_at
		FROM users WHERE email = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.Role, &user.IsActive, &user.IsVerified,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, phone, role, is_active, is_verified, avatar_url, created_at, updated_at, last_login_at, deleted_at
		FROM users WHERE username = $1 AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.Role, &user.IsActive, &user.IsVerified,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *userRepository) GetByEmailOrUsername(emailOrUsername string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, first_name, last_name, phone, role, is_active, is_verified, avatar_url, created_at, updated_at, last_login_at, deleted_at
		FROM users WHERE (email = $1 OR username = $1) AND deleted_at IS NULL
	`
	err := r.db.QueryRow(query, emailOrUsername).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.Role, &user.IsActive, &user.IsVerified,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *userRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
		UPDATE users 
		SET email=$2, username=$3, first_name=$4, last_name=$5, phone=$6, role=$7, is_active=$8, is_verified=$9, avatar_url=$10, updated_at=$11
		WHERE id=$1
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.FirstName, user.LastName, user.Phone, user.Role, user.IsActive, user.IsVerified, user.AvatarURL, user.UpdatedAt)
	return err
}

func (r *userRepository) Delete(id string) error {
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *userRepository) List(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, email, username, first_name, last_name, phone, role, is_active, is_verified, avatar_url, created_at, updated_at, last_login_at
		FROM users WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Phone, &user.Role, &user.IsActive, &user.IsVerified, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) GetStats() (*models.UserStats, error) {
	stats := &models.UserStats{}
	query := `
		SELECT 
			COUNT(*) as total_users,
			COUNT(*) FILTER (WHERE is_active = true) as active_users,
			COUNT(*) FILTER (WHERE is_verified = true) as verified_users,
			COUNT(*) FILTER (WHERE role = 'admin') as admin_users
		FROM users WHERE deleted_at IS NULL
	`
	err := r.db.QueryRow(query).Scan(&stats.TotalUsers, &stats.ActiveUsers, &stats.VerifiedUsers, &stats.AdminUsers)
	return stats, err
}

func (r *userRepository) UpdateLastLogin(userID string) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

func (r *userRepository) CreateRefreshToken(token *models.RefreshToken) error {
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt, token.IPAddress, token.UserAgent)
	return err
}

func (r *userRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}
	query := `SELECT id, user_id, token, expires_at, created_at, revoked_at, ip_address, user_agent FROM refresh_tokens WHERE token = $1`
	err := r.db.QueryRow(query, token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt, &rt.RevokedAt, &rt.IPAddress, &rt.UserAgent)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found")
	}
	return rt, err
}

func (r *userRepository) RevokeRefreshToken(token string) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE token = $2`
	_, err := r.db.Exec(query, time.Now(), token)
	return err
}

func (r *userRepository) DeleteExpiredRefreshTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	_, err := r.db.Exec(query, time.Now())
	return err
}

func (r *userRepository) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	return err
}

func (r *userRepository) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	prt := &models.PasswordResetToken{}
	query := `SELECT id, user_id, token, expires_at, created_at, used_at FROM password_reset_tokens WHERE token = $1`
	err := r.db.QueryRow(query, token).Scan(&prt.ID, &prt.UserID, &prt.Token, &prt.ExpiresAt, &prt.CreatedAt, &prt.UsedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("password reset token not found")
	}
	return prt, err
}

func (r *userRepository) MarkPasswordResetTokenUsed(token string) error {
	query := `UPDATE password_reset_tokens SET used_at = $1 WHERE token = $2`
	_, err := r.db.Exec(query, time.Now(), token)
	return err
}

func (r *userRepository) CreateAuditLog(log *models.AuditLog) error {
	log.ID = uuid.New().String()
	log.CreatedAt = time.Now()
	query := `
		INSERT INTO audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query, log.ID, log.UserID, log.Action, log.Resource, log.ResourceID, log.Details, log.IPAddress, log.UserAgent, log.CreatedAt)
	return err
}
