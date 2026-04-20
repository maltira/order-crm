package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"order-crm/internal/model"
	"order-crm/internal/model/dto"
	"order-crm/internal/repository"
	"order-crm/pkg/utils"
)

type UserService interface {
	Login(login string, password string) (*model.User, error)
	GenerateAndSaveTokens(user *model.User) (string, string, error)
	RevokeRefreshToken(refreshToken string) error
	RefreshToken(refreshToken string) (*model.User, error)

	CreateUser(req *dto.CreateUserRequest) (*model.User, error)

	GetUserByID(id int) (*model.User, error)
	GetUserByLogin(login string) (*model.User, error)
	GetAllUsers() ([]model.User, error)

	UpdateUser(req *dto.UpdateUserRequest) error
	DeleteUser(id int) error
	ChangePassword(req *dto.ChangePasswordRequest) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService - конструктор сервиса
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Login(login string, password string) (*model.User, error) {
	user, err := s.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	// Проверка пароля
	hash := md5.Sum([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])
	if user.Pass != passwordHash {
		return nil, errors.New("invalid credentials")
	}

	// проверка блокировки
	if user.IsBlocked == 1 {
		return nil, errors.New("user is blocked")
	}

	return user, nil
}

func (s *userService) GenerateAndSaveTokens(user *model.User) (string, string, error) {
	// генерация access и refresh
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Role.Code, user.Role.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	// сохранить токен в бд
	err = s.repo.SaveRefreshToken(user, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) RevokeRefreshToken(refreshToken string) error {
	_, err := s.repo.GetTokenInfo(refreshToken)
	if err != nil {
		return err
	}
	return s.repo.RevokeRefreshToken(refreshToken)
}

func (s *userService) RefreshToken(refreshToken string) (*model.User, error) {
	// проверка валидности токена
	userID, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// проверка пользователя
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user.IsBlocked == 1 {
		return nil, errors.New("user is blocked")
	}

	// инвалидация прошлого токена
	err = s.RevokeRefreshToken(refreshToken)
	if err != nil {
		log.Println("Revoke refresh token error:", err)
	}

	return user, nil
}

// CreateUser - создание нового пользователя
func (s *userService) CreateUser(req *dto.CreateUserRequest) (*model.User, error) {
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
func (s *userService) UpdateUser(req *dto.UpdateUserRequest) error {
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
	return s.repo.DeleteUser(id)
}

// ChangePassword - смена пароля пользователя
func (s *userService) ChangePassword(req *dto.ChangePasswordRequest) error {
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
