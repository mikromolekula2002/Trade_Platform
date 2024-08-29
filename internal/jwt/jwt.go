package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct{}

type JWTService interface {
	ExtractToken(tokenString, jwtKey string) (*Claims, error)
	VerifyToken(tokenString, jwtKey string) error
	CreateToken(userLogin string, jwtKey string) (string, error)
	DeleteToken(userLogin string, jwtKey string) (string, error)
}

// Claims структура для хранения данных в токене
type Claims struct {
	UserLogin string `json:"user_login"`
	jwt.StandardClaims
}

func InitJWT() *JWTManager {
	return &JWTManager{} // Возвращаем конкретную реализацию интерфейса
}

// CreateToken создает новый JWT токен
func (j *JWTManager) CreateToken(userLogin string, jwtKey string) (string, error) {
	if jwtKey == "" {
		return "", fmt.Errorf("jwt.go: CreateToken: пустой jwt ключ")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserLogin: userLogin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Trade-Platform",
			Subject:   "auth_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", fmt.Errorf("jwt.go: CreateToken:\n%v", err)
	}

	return signedToken, nil
}

// DeleteToken удаляет куку, точнее говоря обновляет токен на просроченный
func (j *JWTManager) DeleteToken(userLogin string, jwtKey string) (string, error) {
	if jwtKey == "" {
		return "", fmt.Errorf("jwt.go: CreateToken: пустой jwt ключ")
	}

	// Установить время истечения в прошлое, чтобы токен был немедленно просрочен
	expirationTime := time.Now().Add(-1 * time.Hour)

	claims := &Claims{
		UserLogin: userLogin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Trade-Platform",
			Subject:   "auth_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", fmt.Errorf("jwt.go: CreateToken:\n%v", err)
	}

	return signedToken, nil
}

// VerifyToken проверяет действительность JWT токена
func (j *JWTManager) VerifyToken(tokenString, jwtKey string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return fmt.Errorf("jwt.go: VerifyToken: ERROR SIGNATURE INVALID: \n%v", err)
		}
		return fmt.Errorf("jwt.go: VerifyToken:\n%v", err)
	}

	if !token.Valid {
		return fmt.Errorf("jwt.go: VerifyToken: Токен не валиден: \n%v", err)
	}
	return nil
}

// ExtractToken извлекает данные из JWT токена
func (j *JWTManager) ExtractToken(tokenString, jwtKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.go: ExtractToken:\n%v", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("jwt.go: ExtractToken:\n%v", err)
	}
	return claims, nil
}
