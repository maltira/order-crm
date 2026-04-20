package repository

import (
	"database/sql"
	"errors"
	"order-crm/internal/model"
)

type RoleRepository interface {
	GetByID(id int) (*model.Role, error)
	GetByCode(code string) (*model.Role, error)
	GetAll() ([]model.Role, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(id int) (*model.Role, error) {
	query := `SELECT * FROM roles WHERE id = $1`

	var role model.Role
	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Code, &role.Label)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("роль не найдена")
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetByCode(code string) (*model.Role, error) {
	query := `SELECT * FROM roles WHERE code = $1`

	var role model.Role
	err := r.db.QueryRow(query, code).Scan(&role.ID, &role.Code, &role.Label)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("роль не найдена")
	}
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetAll() ([]model.Role, error) {
	query := `SELECT * FROM roles ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		err := rows.Scan(&role.ID, &role.Code, &role.Label)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
