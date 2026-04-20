package service

import (
	"errors"
	"order-crm/internal/model"
	"order-crm/internal/model/dto"
	"order-crm/internal/repository"
)

type ClientService interface {
	GetAllClients() ([]model.Client, error)
	GetClientByID(id int) (*model.Client, error)
	CreateClient(req *dto.ClientRequest) (*model.Client, error)
	UpdateClient(id int, req *dto.ClientRequest) error
	DeleteClient(id int) error
}

type clientService struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientService{repo: repo}
}

func (sc *clientService) GetAllClients() ([]model.Client, error) {
	return sc.repo.GetAllClients()
}

func (sc *clientService) GetClientByID(id int) (*model.Client, error) {
	return sc.repo.GetClientByID(id)
}

func (sc *clientService) CreateClient(req *dto.ClientRequest) (*model.Client, error) {
	if len(req.Label) <= 0 {
		return nil, errors.New("label can't be empty")
	}
	return sc.repo.CreateClient(req.Label)
}

func (sc *clientService) UpdateClient(id int, req *dto.ClientRequest) error {
	if len(req.Label) <= 0 {
		return errors.New("label can't be empty")
	}
	return sc.repo.UpdateClient(id, req.Label)
}

func (sc *clientService) DeleteClient(id int) error {
	return sc.repo.DeleteClient(id)
}
