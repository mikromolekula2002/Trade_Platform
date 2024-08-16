package repository

import (
	"fmt"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
)

// метод обновления инфо в профиле
func (p *PostgreSQL) UpdateUserData(userData *models.UserData) error {
	result := p.DB.Model(&models.UserData{}).
		Where("login = ?", userData.Login).
		Updates(userData)

	if result.Error != nil {
		return fmt.Errorf("repository: UpdateUserData - ошибка обновления данных аккаунта.\n ERROR: %v", result.Error)
	}
	return nil
}

// метод создания профиля
func (p *PostgreSQL) SaveUserData(userdata *models.UserData) error {
	result := p.DB.Create(&userdata)
	if result.Error != nil {
		return fmt.Errorf("repository: SaveUserData - ошибка сохранения данных аккаунта.\n ERROR: %v", result.Error)
	}
	return nil
}

// метод удаления профиля
func (p *PostgreSQL) DelUserData(login string) error {
	result := p.DB.Delete(&models.UserData{}, login)
	if result.Error != nil {
		return fmt.Errorf("repository: DelUserData - ошибка удаления данных аккаунта.\n ERROR: %v", result.Error)
	}
	return nil
}

// метод получения инфо о профиле
func (p *PostgreSQL) GetUserData(login string) (*models.UserData, error) {
	var userData models.UserData
	result := p.DB.Where("login = ?", login).First(&userData)
	if result.Error != nil {
		return nil, fmt.Errorf("repository: GetUserData - ошибка получения данных аккаунта из БД.\n ERROR: %v", result.Error)
	}
	return &userData, nil
}
