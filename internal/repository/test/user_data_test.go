package test

import (
	"testing"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserData(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}
	userData := &models.UserData{
		Login:        "testuser",
		Name:         "Test User",
		FirstName:    "Test",
		PhoneNumber:  "1234567890",
		ImageProfile: "test_image.jpg",
	}

	// Сначала создаем пользователя
	err = repo.SaveUserData(userData)
	require.NoError(t, err, "не удалось создать пользователя")

	// Обновляем данные пользователя
	userData.Name = "Updated User"
	err = repo.UpdateUserData(userData)
	require.NoError(t, err, "не удалось обновить данные пользователя")

	// Проверяем обновленные данные
	updatedUserData, err := repo.GetUserData(userData.Login)
	require.NoError(t, err, "не удалось получить данные пользователя")
	assert.Equal(t, "Updated User", updatedUserData.Name, "имя пользователя не обновлено")
}

func TestSaveUserData(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}
	userData := &models.UserData{
		Login:        "testuser",
		Name:         "Test User",
		FirstName:    "Test",
		PhoneNumber:  "1234567890",
		ImageProfile: "test_image.jpg",
	}

	// Сохраняем пользователя
	err = repo.SaveUserData(userData)
	require.NoError(t, err, "не удалось создать пользователя")

	// Проверяем, что пользователь был сохранен
	savedUserData, err := repo.GetUserData(userData.Login)
	require.NoError(t, err, "не удалось получить данные пользователя")
	assert.Equal(t, userData.Login, savedUserData.Login, "логин пользователя не совпадает")
}

func TestGetUserData(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}
	userData := &models.UserData{
		Login:        "testuser",
		Name:         "Test User",
		FirstName:    "Test",
		PhoneNumber:  "1234567890",
		ImageProfile: "test_image.jpg",
	}

	// Сначала создаем пользователя
	err = repo.SaveUserData(userData)
	require.NoError(t, err, "не удалось создать пользователя")

	// Получаем данные пользователя
	fetchedUserData, err := repo.GetUserData(userData.Login)
	require.NoError(t, err, "не удалось получить данные пользователя")
	assert.Equal(t, userData.Login, fetchedUserData.Login, "логин пользователя не совпадает")
	assert.Equal(t, userData.Name, fetchedUserData.Name, "имя пользователя не совпадает")
}
