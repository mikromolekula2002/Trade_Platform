package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mikromolekula2002/Trade_Platform/internal/mocks"
	"github.com/stretchr/testify/require"
)

func TestUserService_SendToken_Table(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		expectedPath string
		expectedTime time.Duration
		expectError  bool
	}{
		{
			name:         "Valid Token",
			token:        "test-token",
			expectedPath: "/",
			expectedTime: time.Hour * 24,
			expectError:  false,
		},
		{
			name:         "Empty Token",
			token:        "",
			expectedPath: "/",
			expectedTime: time.Hour * 24,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация UserService
			userService := UserService{}

			// Создаем тестовый ответ и записываем туда куку
			w := httptest.NewRecorder()

			// Вызываем метод SendToken с тестовым токеном
			userService.SendToken(w, tt.token)

			// Получаем куки из записи ответа
			resp := w.Result()
			cookies := resp.Cookies()

			// Находим куку с именем "jwt"
			var foundCookie *http.Cookie
			for _, c := range cookies {
				if c.Name == "jwt" {
					foundCookie = c
					break
				}
			}

			if tt.expectError {
				require.Nil(t, foundCookie)
			} else {
				require.NotNil(t, foundCookie)
				require.Equal(t, tt.token, foundCookie.Value)
				require.True(t, foundCookie.HttpOnly)
				require.WithinDuration(t, time.Now().Add(tt.expectedTime), foundCookie.Expires, time.Minute)
				require.Equal(t, tt.expectedPath, foundCookie.Path)
			}
		})
	}
}

func TestUserService_CheckToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockJWT := mocks.NewMockJWTService(mockCtrl)
	userService := UserService{
		jwt:    mockJWT,
		jwtKey: "test-key",
	}

	// Тестовые данные
	validToken := "valid-token"
	invalidToken := "invalid-token"

	tests := []struct {
		name        string
		tokenString string
		mockFunc    func()
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			mockFunc: func() {
				mockJWT.EXPECT().VerifyToken(validToken, "test-key").Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "Invalid Token",
			tokenString: invalidToken,
			mockFunc: func() {
				mockJWT.EXPECT().VerifyToken(invalidToken, "test-key").Return(fmt.Errorf("invalid token"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := userService.CheckToken(tt.tokenString)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUserService_SendExpiredToken(t *testing.T) {
	userService := UserService{}

	// Создаем тестовый ответ и записываем туда куку
	w := httptest.NewRecorder()

	// Вызываем метод SendExpiredToken
	userService.SendExpiredToken(w, "test-token")

	// Получаем куку из записи ответа
	resp := w.Result()
	cookies := resp.Cookies()

	// Ищем куку с именем "jwt"
	var jwtCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "jwt" {
			jwtCookie = c
			break
		}
	}

	// Проверяем, что кука найдена и она просрочена
	require.NotNil(t, jwtCookie, "Cookie 'jwt' should not be nil")
	require.Equal(t, "", jwtCookie.Value)
	require.WithinDuration(t, time.Now().Add(-time.Hour), jwtCookie.Expires, time.Minute)
}
