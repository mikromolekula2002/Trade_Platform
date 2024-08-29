package service

import (
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/mikromolekula2002/Trade_Platform/internal/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser_CreateTokenError(t *testing.T) {
	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	// Инициализируем UserService
	s := &UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "testkey",
	}

	// Тестовые данные
	login := "testuser"
	password := "testpass"

	// Настройка ожиданий на моки
	mockRepo.EXPECT().SaveUser(gomock.Any()).Return(nil)
	mockRepo.EXPECT().SaveUserData(gomock.Any()).Return(nil)
	mockJWT.EXPECT().CreateToken(login, s.jwtKey).Return("", errors.New("token creation error"))

	// Вызов тестируемого метода
	_, err := s.RegisterUser(login, password)

	// Проверка ошибок
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token creation error")
}

func TestRegisterUser_Success(t *testing.T) {
	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	s := &UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "testkey",
	}

	login := "testuser"
	password := "testpass"

	// Set up mocks
	mockRepo.EXPECT().SaveUser(gomock.Any()).Return(nil)
	mockRepo.EXPECT().SaveUserData(gomock.Any()).Return(nil)
	mockJWT.EXPECT().CreateToken(login, s.jwtKey).Return("mockToken", nil)

	// Call method
	token, err := s.RegisterUser(login, password)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "mockToken", token)
}

func TestRegisterUser_ValidateUsernameError(t *testing.T) {
	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	s := &UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "testkey",
	}

	// Подмена bcrypt.GenerateFromPassword, чтобы он не выполнялся в тестах
	patch := monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
		return nil, errors.New("bcrypt не должен вызываться")
	})
	defer patch.Unpatch()

	// Test cases for validation errors
	tests := []struct {
		login    string
		password string
		expected string
	}{
		{"", "testpass", "логин не может быть пустым"},
		{"testuser!", "testpass", "логин не может содержать спец. символы"},
		{"testuserlonglogin", "testpass", "логин не может превышать длину в 20 символов"},
		{"testuser", "", "пароль не может быть пустым"},
		{"testuser", "testpass!", "пароль не может содержать спец. символы"},
		{"testuser", "longpasswordlongpassword", "пароль не может превышать длину в 20 символов"},
	}

	for _, test := range tests {
		t.Run(test.login+"_"+test.password, func(t *testing.T) {

			// Подмена bcrypt.GenerateFromPassword, чтобы он не выполнялся в тестах
			patch := monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
				return nil, errors.New(test.expected)
			})
			defer patch.Unpatch()

			_, err := s.RegisterUser(test.login, test.password)

			// Проверяем, что произошла ошибка валидации
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.expected)
		})
	}
}

func TestRegisterUser_HashPasswordError(t *testing.T) {
	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	s := &UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "testkey",
	}

	login := "testuser"
	password := "testpass"

	// Подмена bcrypt.GenerateFromPassword для имитации ошибки хеширования пароля

	patch := monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
		return nil, errors.New("ошибка хеширования пароля")
	})
	defer patch.Unpatch()

	// Вызов метода
	_, err := s.RegisterUser(login, password)

	// Проверки
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка хеширования пароля")
}

func TestRegisterUser_SaveUserError(t *testing.T) {
	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	s := &UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "testkey",
	}

	login := "testuser"
	password := "testpass"

	// Set up mocks
	mockRepo.EXPECT().SaveUser(gomock.Any()).Return(errors.New("save user error"))

	// Call method
	_, err := s.RegisterUser(login, password)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save user error")
}

func TestRegisterUser_SaveUserDataError(t *testing.T) {
	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем моки
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	s := &UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "testkey",
	}

	login := "testuser"
	password := "testpass"

	// Set up mocks
	mockRepo.EXPECT().SaveUser(gomock.Any()).Return(nil)
	mockRepo.EXPECT().SaveUserData(gomock.Any()).Return(errors.New("save user data error"))

	// Call method
	_, err := s.RegisterUser(login, password)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "save user data error")
}
