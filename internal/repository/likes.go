package repository

import (
	"errors"
	"fmt"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/utils"
	"gorm.io/gorm"
)

// метод создания объявы
func (p *PostgreSQL) SaveLikes(likes *models.Likes) error {
	// Проверка на существующую запись
	var existingLike models.Likes
	err := p.DB.Where("user_login = ? AND ads_id = ?", likes.User_Login, likes.Ads_Id).First(&existingLike).Error
	if err == nil {
		// Если запись уже существует, возвращаем ошибку
		return utils.ErrAlreadyExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Если произошла ошибка, кроме "запись не найдена", возвращаем её
		return fmt.Errorf("repository: SaveLikes - ошибка проверки существующего лайка.\n ERROR: %v", err)
	}

	// Если запись не найдена, добавляем новый лайк
	result := p.DB.Create(&likes)
	if result.Error != nil {
		return fmt.Errorf("repository: SaveLikes - ошибка сохранения лайка.\n ERROR: %v", result.Error)
	}

	// Обновление количества лайков в объявлении
	resultLike := p.DB.Model(&models.UserAds{}).Where("ads_id = ?", likes.Ads_Id).UpdateColumn("ads_likes", gorm.Expr("ads_likes + ?", 1))
	if resultLike.Error != nil {
		return fmt.Errorf("repository: SaveLikes - ошибка обновления кол-ва лайков в объявлении.\n ERROR: %v", resultLike.Error)
	}

	return nil
}

// метод удаления объявы
func (p *PostgreSQL) DelLikes(userlogin string, adsId string) error {
	result := p.DB.Where("user_login = ? AND ads_id = ?", userlogin, adsId).
		Delete(&models.Likes{})

	if result.Error != nil {
		return fmt.Errorf("repository: DelLikes - ошибка удаления лайка.\n ERROR: %v", result.Error)
	}

	// Обновление количества лайков в объявлении
	resultLike := p.DB.Model(&models.UserAds{}).Where("ads_id = ?", adsId).UpdateColumn("ads_likes", gorm.Expr("ads_likes - ?", 1))
	if resultLike.Error != nil {
		return fmt.Errorf("repository: SaveLikes - ошибка обновления кол-ва лайков в объявлении.\n ERROR: %v", resultLike.Error)
	}

	return nil
}

// метод получения инфо обо всех лайкнутых объявах
func (p *PostgreSQL) GetAllLikes(userlogin string) ([]*models.Likes, error) {
	var usersLikes []*models.Likes
	result := p.DB.Where("user_login = ?", userlogin).
		Find(&usersLikes)

	if result.Error != nil {
		return nil, fmt.Errorf("repository: GetAllLikes - ошибка получения всех лайков.\n ERROR: %v", result.Error)
	}
	return usersLikes, nil
}
