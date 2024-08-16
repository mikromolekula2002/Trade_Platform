package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	jwtKey := "test_secret_key"
	userLogin := "testuser"

	// Основной тест на создание токена
	tokenString, err := CreateToken(userLogin, jwtKey)
	assert.NoError(t, err, "Ошибка при создании токена")
	assert.NotEmpty(t, tokenString, "Созданный токен пуст")

	// Проверка длины токена
	assert.Greater(t, len(tokenString), 0, "Длина созданного токена равна нулю")

	// Проверка содержимого токена (три части, разделенные точками)
	parts := strings.Split(tokenString, ".")
	assert.Equal(t, 3, len(parts), "Токен должен содержать три части, разделенные точками")

	// Проверка срока действия токена
	claims := &Claims{UserLogin: userLogin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Trade-Platform",
			Subject:   "auth_token",
		},
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	assert.NoError(t, err, "Ошибка при парсинге токена")
	assert.True(t, token.Valid, "Токен недействителен")

	assert.WithinDuration(t, time.Now().Add(24*time.Hour), time.Unix(claims.ExpiresAt, 0), time.Minute, "Срок действия токена установлен неправильно")

	// Тест с пустым userLogin
	emptyLoginToken, err := CreateToken("", jwtKey)
	assert.NoError(t, err, "Ошибка при создании токена с пустым userLogin")
	assert.NotEmpty(t, emptyLoginToken, "Созданный токен с пустым userLogin пуст")

	// Тест с пустым ключом
	invalidKeyToken, err := CreateToken(userLogin, "")
	assert.Error(t, err, "Ожидалась ошибка при создании токена с пустым ключом")
	assert.Empty(t, invalidKeyToken, "Созданный токен с пустым ключом не пуст")
}

func TestVerifyToken(t *testing.T) {
	jwtKey := "test_secret_key"
	userLogin := "testuser"

	// Создаем валидный токен
	validTokenString, err := CreateToken(userLogin, jwtKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, validTokenString)

	// Создаем токен с истекшим сроком действия
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
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := expiredToken.SignedString([]byte(jwtKey))
	assert.NoError(t, err)
	assert.NotEmpty(t, expiredTokenString)

	// Модифицируем валидный токен
	modifiedTokenString := validTokenString[:len(validTokenString)-1] + "a"

	testTable := []struct {
		name      string
		token     string
		key       string
		expectErr bool
	}{
		{"ValidToken", validTokenString, jwtKey, false},
		{"InvalidToken", "invalidTokenString", jwtKey, true},
		{"WrongKey", validTokenString, "wrong_secret_key", true},
		{"ExpiredToken", expiredTokenString, jwtKey, true},
		{"MalformedToken", "malformed.token.string", jwtKey, true},
		{"EmptyToken", "", jwtKey, true},
		{"ModifiedToken", modifiedTokenString, jwtKey, true},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyToken(tt.token, tt.key)
			if tt.expectErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Did not expect error but got one")
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	jwtKey := "test_secret_key"
	userLogin := "testuser"

	// Создаем валидный токен
	expirationTimeValid := time.Now().Add(24 * time.Hour)

	claimsValid := &Claims{
		UserLogin: userLogin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeValid.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Trade-Platform",
			Subject:   "auth_token",
		},
	}

	tokenValid := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsValid)
	validTokenString, _ := tokenValid.SignedString([]byte(jwtKey))

	//....................................................................................
	// Создаем токен с истекшим сроком действия
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
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := expiredToken.SignedString([]byte(jwtKey))
	assert.NoError(t, err)
	assert.NotEmpty(t, expiredTokenString)

	// Модифицируем валидный токен
	modifiedTokenString := validTokenString[:len(validTokenString)-1] + "a"

	tests := []struct {
		name      string
		token     string
		key       string
		expectErr bool
		expected  *jwt.Claims
	}{
		{"ValidToken", validTokenString, jwtKey, false, &tokenValid.Claims},
		{"InvalidToken", "invalidTokenString", jwtKey, true, nil},
		{"WrongKey", validTokenString, "wrong_secret_key", true, nil},
		{"ExpiredToken", expiredTokenString, jwtKey, true, nil},
		{"MalformedToken", "malformed.token.string", jwtKey, true, nil},
		{"EmptyToken", "", jwtKey, true, nil},
		{"ModifiedToken", modifiedTokenString, jwtKey, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ExtractToken(tt.token, tt.key)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}
