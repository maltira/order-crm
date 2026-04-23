package repository

import (
	"database/sql"
	"errors"
	"order-crm/config"
	"order-crm/internal/model"
	"time"
)

type UserRepository interface {
	SaveRefreshToken(user *model.User, refreshToken string) error
	RevokeRefreshToken(token string) error
	GetTokenInfo(token string) (*model.RefreshToken, error)

	CreateUser(req *model.User) error

	GetUserByID(id int) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetAllUsers() ([]model.User, error)

	UpdateUser(user *model.User) error
	DeleteUser(id int) error
	UpdatePassword(id int, passwordHash string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveRefreshToken(user *model.User, refreshToken string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`

	err := r.db.QueryRow(query, user.ID, refreshToken, time.Now().Add(config.Env.RefreshTokenDuration)).Err()
	if err != nil {
		return errors.New("ошибка сохранения токена в бд: " + err.Error())
	}

	return nil
}

func (r *userRepository) RevokeRefreshToken(token string) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE token = $1
	`

	result, err := r.db.Exec(query, token)
	if err != nil {
		return errors.New("ошибка инвалидации токена: " + err.Error())
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("токен не найден или уже отозван")
	}

	return nil
}

func (r *userRepository) GetTokenInfo(token string) (*model.RefreshToken, error) {
	query := `
		SELECT user_id, token, expires_at, created_at 
		FROM refresh_tokens WHERE token = $1
	`
	var refresh model.RefreshToken
	err := r.db.QueryRow(query, token).Scan(
		&refresh.UserID,
		&refresh.Token,
		&refresh.ExpiresAt,
		&refresh.CreatedAt,
	)
	if err != nil {
		return nil, errors.New("Ошибка поиска токена: " + err.Error())
	}
	return &refresh, nil
}

func (r *userRepository) CreateUser(req *model.User) error {
	query := `
		INSERT INTO users (login, pass, fio, id_role, is_blocked)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRow(query, req.Login, req.Pass, req.Fio, req.IDRole, req.IsBlocked).Scan(&req.ID)
	if err != nil {
		return errors.New("ошибка создания пользователя: " + err.Error())
	}

	return nil
}

func (r *userRepository) GetUserByID(id int) (*model.User, error) {
	query := `
		SELECT 
			u.id, 
			u.login, 
			u.fio, 
			u.id_role, 
			u.is_blocked,
			r.code, 
			r.label
		FROM users u
		INNER JOIN roles r ON r.id = u.id_role
		WHERE u.id = $1
		ORDER BY u.id
	`

	var user model.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Fio,
		&user.IDRole,
		&user.IsBlocked,
		&user.Role.Code,
		&user.Role.Label,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("пользователь не найден")
	}
	if err != nil {
		return nil, errors.New("ошибка получения пользователя: " + err.Error())
	}

	return &user, nil
}

func (r *userRepository) GetUserByLogin(login string) (*model.User, error) {
	query := `
		SELECT 
			u.id, 
			u.login, 
			u.pass,
			u.fio, 
			u.id_role, 
			u.is_blocked,
			r.code,
			r.label
		FROM users u
		INNER JOIN roles r ON r.id = u.id_role
		WHERE u.login = $1
		ORDER BY u.id
	`

	var user model.User
	err := r.db.QueryRow(query, login).Scan(
		&user.ID,
		&user.Login,
		&user.Pass,
		&user.Fio,
		&user.IDRole,
		&user.IsBlocked,
		&user.Role.Code,
		&user.Role.Label,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("ошибка получения пользователя: " + err.Error())
	}

	return &user, nil
}

func (r *userRepository) GetAllUsers() ([]model.User, error) {
	query := `
		SELECT 
			u.id, 
			u.login, 
			u.fio, 
			u.id_role, 
			u.is_blocked,
			r.code as role_code, 
			r.label as role_label
		FROM users u
		INNER JOIN roles r ON r.id = u.id_role
		ORDER BY u.id
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.New("ошибка получения списка пользователей: " + err.Error())
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err = rows.Scan(
			&u.ID,
			&u.Login,
			&u.Fio,
			&u.IDRole,
			&u.IsBlocked,
			&u.Role.Code,
			&u.Role.Label,
		)
		if err != nil {
			return nil, errors.New("ошибка сканирования данных: " + err.Error())
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	query := `
		UPDATE users 
		SET login = $1, 
		    fio = $2, 
		    id_role = $3, 
		    is_blocked = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(query,
		user.Login,
		user.Fio,
		user.IDRole,
		user.IsBlocked,
		user.ID,
	)

	if err != nil {
		return errors.New("ошибка обновления пользователя: " + err.Error())
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}

func (r *userRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return errors.New("ошибка удаления пользователя: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}

func (r *userRepository) UpdatePassword(id int, passwordHash string) error {
	query := `
		UPDATE users 
		SET pass = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, passwordHash, id)
	if err != nil {
		return errors.New("ошибка обновления пароля: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("пользователь не найден")
	}

	return nil
}
