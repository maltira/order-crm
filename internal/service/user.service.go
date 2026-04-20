package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"order-crm/internal/model"
	"order-crm/internal/repository"
)

type UserService interface {
	CreateUser(req *model.CreateUserRequest) (*model.User, error)

	GetUserByID(id int) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetAllUsers() ([]model.User, error)

	UpdateUser(req *model.UpdateUserRequest) error
	DeleteUser(id int) error
	ChangePassword(req *model.ChangePasswordRequest) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService - конструктор сервиса
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser - создание нового пользователя
func (s *userService) CreateUser(req *model.CreateUserRequest) (*model.User, error) {
	// 1. Проверяем, не существует ли пользователь с таким логином
	existingUser, _ := s.repo.GetUserByLogin(req.Login)
	if existingUser != nil {
		return nil, errors.New("пользователь с таким логином уже существует")
	}

	// 2. Валидация входных данных
	if req.Login == "" {
		return nil, errors.New("логин не может быть пустым")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("пароль должен содержать мин. 8 символов")
	}
	if req.Fio == "" {
		return nil, errors.New("ФИО не может быть пустым")
	}
	isBlocked := 0
	if req.IsBlocked {
		isBlocked = 1
	}

	// 3. Хэшируем пароль
	hash := md5.Sum([]byte(req.Password))
	passwordHash := hex.EncodeToString(hash[:])

	// 4. Создаём пользователя
	user := &model.User{
		Login:     req.Login,
		Pass:      passwordHash,
		Fio:       req.Fio,
		IDRole:    req.IDRole,
		IsBlocked: isBlocked,
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// 5. Формируем ответ
	return user, nil
}

// GetUserByID - получение пользователя по ID
func (s *userService) GetUserByID(id int) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("id должен быть > 0")
	}
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByLogin - получение пользователя по логину
func (s *userService) GetUserByLogin(login string) (*model.User, error) {
	user, err := s.repo.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers - получение всех пользователей
func (s *userService) GetAllUsers() ([]model.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser - обновление данных пользователя
func (s *userService) UpdateUser(req *model.UpdateUserRequest) error {
	// 1. Проверяем, существует ли пользователь
	user, err := s.repo.GetUserByID(req.ID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// 2. Валидация
	if req.Login != nil {
		if len(*req.Login) < 3 {
			return errors.New("логин должен содержать минимум 3 символа")
		}
		user.Login = *req.Login
	}
	if req.Fio != nil {
		if len(*req.Fio) < 3 {
			return errors.New("фио должно содержать минимум 3 символа")
		}
		user.Fio = *req.Fio
	}
	if req.IDRole != nil {
		user.IDRole = *req.IDRole
	}
	if req.IsBlocked != nil {
		isBlocked := 0
		if *req.IsBlocked {
			isBlocked = 1
		}
		user.IsBlocked = isBlocked
	}

	return s.repo.UpdateUser(user)
}

// DeleteUser - удаление пользователя
func (s *userService) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("ID должен быть положительным числом")
	}
	return s.repo.DeleteUser(id)
}

// ChangePassword - смена пароля пользователя
func (s *userService) ChangePassword(req *model.ChangePasswordRequest) error {
	// 1. Проверяем новый пароль
	if len(req.NewPassword) < 8 {
		return errors.New("пароль должен содержать минимум 8 символов")
	}

	// 2. Хэшируем новый пароль
	newHash := md5.Sum([]byte(req.NewPassword))
	newPasswordHash := hex.EncodeToString(newHash[:])

	// 3. Обновляем пароль
	return s.repo.UpdatePassword(req.ID, newPasswordHash)
}
