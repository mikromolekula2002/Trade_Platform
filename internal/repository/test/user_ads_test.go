package test

import (
	"testing"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserAds(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}
	userAds := &models.UserAds{
		User_Id:         "testuser",
		Ads_Id:          "1",
		Image_1:         "image1.jpg",
		Image_2:         "image2.jpg",
		Image_3:         "image3.jpg",
		Ads_Name:        "Test Ad",
		Ads_Description: "Test Ad Description",
		Ads_Price:       100.00,
	}

	// Сначала создаем объявление
	_, err = repo.SaveUserAds(userAds)
	require.NoError(t, err, "не удалось создать объявление")

	// Обновляем объявление
	userAds.Ads_Name = "Updated Ad"
	err = repo.UpdateUserAds(userAds)
	require.NoError(t, err, "не удалось обновить объявление")

	// Проверяем обновленное объявление
	updatedAds, err := repo.GetOneAds("1")
	require.NoError(t, err, "не удалось получить объявление")
	assert.Equal(t, "Updated Ad", updatedAds.Ads_Name, "название объявления не обновлено")
}

func TestSaveUserAds(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}
	userAds := &models.UserAds{
		User_Id:         "testuser",
		Ads_Id:          "1",
		Image_1:         "image1.jpg",
		Image_2:         "image2.jpg",
		Image_3:         "image3.jpg",
		Ads_Name:        "Test Ad",
		Ads_Description: "Test Ad Description",
		Ads_Price:       100.00,
	}

	// Сохраняем объявление
	id, err := repo.SaveUserAds(userAds)
	require.NoError(t, err, "не удалось создать объявление")
	assert.Equal(t, "1", id, "id объявления не совпадает")

	// Проверяем, что объявление было сохранено
	savedAds, err := repo.GetOneAds("1")
	require.NoError(t, err, "не удалось получить объявление")
	assert.Equal(t, userAds.Ads_Name, savedAds.Ads_Name, "название объявления не совпадает")
}

func TestGetUserAds(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}
	userAds := &models.UserAds{
		User_Id:         "testuser",
		Ads_Id:          "1",
		Image_1:         "image1.jpg",
		Image_2:         "image2.jpg",
		Image_3:         "image3.jpg",
		Ads_Name:        "Test Ad",
		Ads_Description: "Test Ad Description",
		Ads_Price:       100.00,
	}

	// Сначала создаем объявление
	_, err = repo.SaveUserAds(userAds)
	require.NoError(t, err, "не удалось создать объявление")

	// Получаем все объявления пользователя
	adsList, err := repo.GetUserAds("testuser")
	require.NoError(t, err, "не удалось получить объявления пользователя")
	assert.Len(t, adsList, 1, "объявления пользователя не найдены")
	assert.Equal(t, "Test Ad", adsList[0].Ads_Name, "название объявления не совпадает")
}

func TestGetAllAds(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Не удалось создать тестовую базу данных: %v", err)
	}

	defer truncateTables(db)

	repo := repository.PostgreSQL{DB: db}
	userAds1 := &models.UserAds{
		User_Id:         "user1",
		Ads_Id:          "1",
		Image_1:         "image1.jpg",
		Image_2:         "image2.jpg",
		Image_3:         "image3.jpg",
		Ads_Name:        "Ad 1",
		Ads_Description: "Description 1",
		Ads_Price:       100.00,
	}
	userAds2 := &models.UserAds{
		User_Id:         "user2",
		Ads_Id:          "2",
		Image_1:         "image1.jpg",
		Image_2:         "image2.jpg",
		Image_3:         "image3.jpg",
		Ads_Name:        "Ad 2",
		Ads_Description: "Description 2",
		Ads_Price:       200.00,
	}

	// Сначала создаем объявления
	_, err = repo.SaveUserAds(userAds1)
	require.NoError(t, err, "не удалось создать объявление 1")
	_, err = repo.SaveUserAds(userAds2)
	require.NoError(t, err, "не удалось создать объявление 2")

	// Получаем все объявления
	adsList, err := repo.GetAllAds()
	require.NoError(t, err, "не удалось получить все объявления")
	assert.Len(t, adsList, 2, "объявления не найдены")
	assert.ElementsMatch(t, []*models.UserAds{userAds1, userAds2}, adsList, "список объявлений не совпадает")
}
