package test

import (
	"testing"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/repository"
)

func TestUpdatePassword(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}

	// Создание тестового пользователя
	user := &models.User{
		Login:        "testusers",
		HashPassword: "oldpasswordhash",
	}
	err = repo.SaveUser(user)
	if err != nil {
		t.Fatalf("Не удалось сохранить тестового пользователя: %v", err)
	}

	// Обновление пароля
	newHashPassword := "newpasswordhash"
	err = repo.UpdatePassword(user.Login, newHashPassword)
	if err != nil {
		t.Fatalf("Ошибка при обновлении пароля: %v", err)
	}

	// Проверка обновленного пароля
	updatedUser, err := repo.GetUser(user.Login)
	if err != nil {
		t.Fatalf("Ошибка при получении обновленного пользователя: %v", err)
	}
	if updatedUser.HashPassword != newHashPassword {
		t.Errorf("Ожидался пароль %s, но получен %s", newHashPassword, updatedUser.HashPassword)
	}

}

func TestSaveUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}

	user := &models.User{
		Login:        "testusers",
		HashPassword: "testpasswordhash",
	}

	// Сохранение пользователя
	err = repo.SaveUser(user)
	if err != nil {
		t.Fatalf("Ошибка при сохранении пользователя: %v", err)
	}

	// Проверка сохраненного пользователя
	savedUser, err := repo.GetUser(user.Login)
	if err != nil {
		t.Fatalf("Ошибка при получении сохраненного пользователя: %v", err)
	}
	if savedUser.Login != user.Login {
		t.Errorf("Ожидался логин %s, но получен %s", user.Login, savedUser.Login)
	}

}

func TestGetUser(t *testing.T) {

	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}

	user := &models.User{
		Login:        "testusers",
		HashPassword: "testpasswordhash",
	}

	// Сохранение пользователя
	err = repo.SaveUser(user)
	if err != nil {
		t.Fatalf("Ошибка при сохранении пользователя: %v", err)
	}

	// Получение пользователя
	retrievedUser, err := repo.GetUser(user.Login)
	if err != nil {
		t.Fatalf("Ошибка при получении пользователя: %v", err)
	}
	if retrievedUser.Login != user.Login {
		t.Errorf("Ожидался логин %s, но получен %s", user.Login, retrievedUser.Login)
	}

}
