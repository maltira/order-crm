package repository

import (
	"database/sql"
	"order-crm/internal/model"
)

type ClientRepository interface {
	GetAllClients() ([]model.Client, error)
	GetClientByID(id int) (*model.Client, error)
	CreateClient(label string) (*model.Client, error)
	UpdateClient(id int, label string) error
	DeleteClient(id int) error
}

type clientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) GetAllClients() ([]model.Client, error) {
	query := `SELECT * FROM clients`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []model.Client
	for rows.Next() {
		var c model.Client
		err = rows.Scan(&c.ID, &c.Label)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func (r *clientRepository) GetClientByID(id int) (*model.Client, error) {
	query := `SELECT * FROM clients WHERE id = $1`

	var c model.Client
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Label)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *clientRepository) CreateClient(label string) (*model.Client, error) {
	query := `INSERT INTO clients (label) VALUES ($1) RETURNING id, label`

	var c model.Client
	err := r.db.QueryRow(query, label).Scan(&c.ID, &c.Label)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *clientRepository) UpdateClient(id int, label string) error {
	query := `UPDATE clients SET label = $1 WHERE id = $2`

	_, err := r.db.Exec(query, label, id)
	return err
}

func (r *clientRepository) DeleteClient(id int) error {
	query := `DELETE FROM clients WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
