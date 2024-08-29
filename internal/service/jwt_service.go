package service

import (
	"fmt"
	"net/http"
	"time"
)

// Функция создания cookie для отправки токена через cookie
func (s *UserService) SendToken(w http.ResponseWriter, jwtToken string) {
	// Отправка токена пользователю (например, через куку)
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24), // Например, токен будет действителен 24 часа
		Path:     "/",                            // Это значит, что токен дает доступ ко всем URL сайта.
	}
	// Отправка токена пользователю через cookie
	http.SetCookie(w, cookie)
}

func (s *UserService) CheckToken(tokenString string) error {
	err := s.jwt.VerifyToken(tokenString, s.jwtKey)
	if err != nil {
		return fmt.Errorf("CheckToken - ошибка проверки токена: %w", err)
	}
	return nil
}

func (s *UserService) SendExpiredToken(w http.ResponseWriter, jwtToken string) {
	cookie := &http.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now().Add(-time.Hour), // Кука уже просрочена
	}

	http.SetCookie(w, cookie)
}
