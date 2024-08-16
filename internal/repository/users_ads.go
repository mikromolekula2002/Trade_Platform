package repository

import (
	"fmt"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/utils"
)

// метод обновления инфо об объяве
func (p *PostgreSQL) UpdateUserAds(userAds *models.UserAds) error {
	result := p.DB.Model(&models.UserAds{}).
		Where("ads_id = ?", userAds.Ads_Id).
		Updates(userAds)

	if result.Error != nil {
		return fmt.Errorf("repository: UpdateUserAds - ошибка обновления данных объявления.\n ERROR: %v", result.Error)
	}
	return nil
}

// метод создания объявы
func (p *PostgreSQL) SaveUserAds(userAds *models.UserAds) (string, error) {
	result := p.DB.Create(&userAds)
	if result.Error != nil {
		return "", fmt.Errorf("repository: SaveUserAds - ошибка сохранения данных объявления.\n ERROR: %v", result.Error)
	}
	id := string(userAds.Ads_Id)
	return id, nil
}

// метод удаления всех объявлений пользователя
func (p *PostgreSQL) DelAllUserAds(login string) error {
	result := p.DB.Delete(&models.UserAds{}, login)
	if result.Error != nil {
		return fmt.Errorf("repository: DelUserAds - ошибка удаления данных объявления.\n ERROR: %v", result.Error)
	}
	return nil
}

// метод удаления конкретной объявы пользователя
func (p *PostgreSQL) DelUserAds(userID string, adsID string) error {
	result := p.DB.Where("user_id = ? AND ads_id = ?", userID, adsID).Delete(&models.UserAds{})
	if result.Error != nil {
		return fmt.Errorf("repository: DelUserAds - ошибка удаления данных объявления.\n ERROR: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNoOwnerAds
	}
	return nil
}

// метод получения инфо об конкретной объяве
func (p *PostgreSQL) GetOneAds(adsID string) (*models.UserAds, error) {
	var userAds models.UserAds
	result := p.DB.Where("ads_id = ?", adsID).First(&userAds)
	if result.Error != nil {
		return nil, fmt.Errorf("repository: GetOneAds - ошибка получения данных объявления из БД.\n ERROR: %v", result.Error)
	}
	return &userAds, nil
}

// метод получения всех объявлений юзера
func (p *PostgreSQL) GetUserAds(login string) ([]*models.UserAds, error) {
	var data []*models.UserAds
	result := p.DB.Where("user_id = ?", login).Find(&data)
	if result.Error != nil {
		return nil, fmt.Errorf("repository: GetUserAds - ошибка получения объявлений юзера из БД.\nERROR: %v", result.Error)
	}
	return data, nil
}

// метод получения инфо обо всех объявах
func (p *PostgreSQL) GetAllAds() ([]*models.UserAds, error) {
	var usersAds []*models.UserAds
	result := p.DB.Find(&usersAds)
	if result.Error != nil {
		return nil, fmt.Errorf("repository: GetAllAds - ошибка получения данных объявлений из БД.\n ERROR: %v", result.Error)
	}
	return usersAds, nil
}
