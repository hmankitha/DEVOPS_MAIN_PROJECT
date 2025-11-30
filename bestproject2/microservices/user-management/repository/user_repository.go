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
	GetByEmailOrUsername(credential string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	List(limit, offset int) ([]*models.User, error)
	GetStats() (*models.UserStats, error)
	UpdateLastLogin(userID string) error

	CreateRefreshToken(token *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	RevokeRefreshToken(token string) error
	DeleteExpiredRefreshTokens() error

	CreatePasswordResetToken(token *models.PasswordResetToken) error
	GetPasswordResetToken(token string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenUsed(token string) error

	CreateAuditLog(log *models.AuditLog) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

/////////////////////////////////////////
// Safe NULL scan helper
/////////////////////////////////////////

func scanUser(row *sql.Row) (*models.User, error) {
	user := &models.User{}
	var lastLogin sql.NullTime
	var deletedAt sql.NullTime

	err := row.Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.Role,
		&user.IsActive, &user.IsVerified, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
		&lastLogin, &deletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}
	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return user, nil
}

/////////////////////////////////////////
// CRUD Implementations
/////////////////////////////////////////

func (r *userRepository) Create(user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
        INSERT INTO users (
            id, email, username, password_hash,
            first_name, last_name, phone, role,
            is_active, is_verified, avatar_url,
            created_at, updated_at
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
    `
	_, err := r.db.Exec(query,
		user.ID, user.Email, user.Username, user.PasswordHash,
		user.FirstName, user.LastName, user.Phone, user.Role,
		user.IsActive, user.IsVerified, user.AvatarURL,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	row := r.db.QueryRow(`
        SELECT id, email, username, password_hash,
               first_name, last_name, phone, role,
               is_active, is_verified, avatar_url,
               created_at, updated_at, last_login_at, deleted_at
        FROM users WHERE id=$1 AND deleted_at IS NULL`, id)

	return scanUser(row)
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow(`
        SELECT id, email, username, password_hash,
               first_name, last_name, phone, role,
               is_active, is_verified, avatar_url,
               created_at, updated_at, last_login_at, deleted_at
        FROM users WHERE email=$1 AND deleted_at IS NULL`, email)

	return scanUser(row)
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(`
        SELECT id, email, username, password_hash,
               first_name, last_name, phone, role,
               is_active, is_verified, avatar_url,
               created_at, updated_at, last_login_at, deleted_at
        FROM users WHERE username=$1 AND deleted_at IS NULL`, username)

	return scanUser(row)
}

func (r *userRepository) GetByEmailOrUsername(c string) (*models.User, error) {
	row := r.db.QueryRow(`
        SELECT id, email, username, password_hash,
               first_name, last_name, phone, role,
               is_active, is_verified, avatar_url,
               created_at, updated_at, last_login_at, deleted_at
        FROM users 
        WHERE (email=$1 OR username=$1) AND deleted_at IS NULL`, c)

	return scanUser(row)
}

func (r *userRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
        UPDATE users SET
            email=$2, username=$3, first_name=$4, last_name=$5,
            phone=$6, role=$7, is_active=$8, is_verified=$9,
            avatar_url=$10, updated_at=$11
        WHERE id=$1
    `
	_, err := r.db.Exec(query,
		user.ID, user.Email, user.Username, user.FirstName, user.LastName,
		user.Phone, user.Role, user.IsActive, user.IsVerified,
		user.AvatarURL, user.UpdatedAt,
	)
	return err
}

func (r *userRepository) Delete(id string) error {
	_, err := r.db.Exec(`UPDATE users SET deleted_at=$1 WHERE id=$2`, time.Now(), id)
	return err
}

func (r *userRepository) List(limit, offset int) ([]*models.User, error) {
	rows, err := r.db.Query(`
        SELECT id, email, username, first_name, last_name,
               phone, role, is_active, is_verified, avatar_url,
               created_at, updated_at, last_login_at
        FROM users WHERE deleted_at IS NULL
        ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Email, &user.Username,
			&user.FirstName, &user.LastName, &user.Phone, &user.Role,
			&user.IsActive, &user.IsVerified, &user.AvatarURL,
			&user.CreatedAt, &user.UpdatedAt, &lastLogin,
		)
		if err != nil {
			return nil, err
		}

		if lastLogin.Valid {
			user.LastLoginAt = &lastLogin.Time
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetStats() (*models.UserStats, error) {
	stats := &models.UserStats{}
	err := r.db.QueryRow(`
        SELECT 
            COUNT(*) AS total_users,
            COUNT(*) FILTER (WHERE is_active=true) AS active_users,
            COUNT(*) FILTER (WHERE is_verified=true) AS verified_users,
            COUNT(*) FILTER (WHERE role='admin') AS admin_users
        FROM users WHERE deleted_at IS NULL`,
	).Scan(
		&stats.TotalUsers,
		&stats.ActiveUsers,
		&stats.VerifiedUsers,
		&stats.AdminUsers,
	)
	return stats, err
}

func (r *userRepository) UpdateLastLogin(userID string) error {
	_, err := r.db.Exec(`UPDATE users SET last_login_at=$1 WHERE id=$2`, time.Now(), userID)
	return err
}

/////////////////////////////////////////
// Refresh Tokens
/////////////////////////////////////////

func (r *userRepository) CreateRefreshToken(t *models.RefreshToken) error {
	t.ID = uuid.New().String()
	_, err := r.db.Exec(`
        INSERT INTO refresh_tokens (id,user_id,token,expires_at,created_at,ip_address,user_agent)
        VALUES ($1,$2,$3,$4,$5,$6,$7)
    `, t.ID, t.UserID, t.Token, t.ExpiresAt, t.CreatedAt, t.IPAddress, t.UserAgent)
	return err
}

func (r *userRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}
	var revoked sql.NullTime

	err := r.db.QueryRow(`
        SELECT id, user_id, token, expires_at, created_at, revoked_at, ip_address, user_agent
        FROM refresh_tokens WHERE token=$1`, token,
	).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt,
		&rt.CreatedAt, &revoked, &rt.IPAddress, &rt.UserAgent,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found")
	}
	if err != nil {
		return nil, err
	}

	if revoked.Valid {
		rt.RevokedAt = &revoked.Time
	}

	return rt, nil
}

func (r *userRepository) RevokeRefreshToken(token string) error {
	_, err := r.db.Exec(`UPDATE refresh_tokens SET revoked_at=$1 WHERE token=$2`, time.Now(), token)
	return err
}

func (r *userRepository) DeleteExpiredRefreshTokens() error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
	return err
}

/////////////////////////////////////////
// Password Reset
/////////////////////////////////////////

func (r *userRepository) CreatePasswordResetToken(t *models.PasswordResetToken) error {
	t.ID = uuid.New().String()
	_, err := r.db.Exec(`
        INSERT INTO password_reset_tokens (id,user_id,token,expires_at,created_at)
        VALUES ($1,$2,$3,$4,$5)
    `, t.ID, t.UserID, t.Token, t.ExpiresAt, t.CreatedAt)
	return err
}

func (r *userRepository) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	pr := &models.PasswordResetToken{}
	var used sql.NullTime

	err := r.db.QueryRow(`
        SELECT id,user_id,token,expires_at,created_at,used_at
        FROM password_reset_tokens WHERE token=$1`, token,
	).Scan(
		&pr.ID, &pr.UserID, &pr.Token, &pr.ExpiresAt,
		&pr.CreatedAt, &used,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("password reset token not found")
	}
	if err != nil {
		return nil, err
	}

	if used.Valid {
		pr.UsedAt = &used.Time
	}

	return pr, nil
}

func (r *userRepository) MarkPasswordResetTokenUsed(token string) error {
	_, err := r.db.Exec(`UPDATE password_reset_tokens SET used_at=$1 WHERE token=$2`, time.Now(), token)
	return err
}

/////////////////////////////////////////
// Audit Logs
/////////////////////////////////////////

func (r *userRepository) CreateAuditLog(log *models.AuditLog) error {
	log.ID = uuid.New().String()
	_, err := r.db.Exec(`
        INSERT INTO audit_logs (id,user_id,action,resource,resource_id,details,ip_address,user_agent,created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
    `, log.ID, log.UserID, log.Action, log.Resource, log.ResourceID,
		log.Details, log.IPAddress, log.UserAgent, time.Now())

	return err
}
