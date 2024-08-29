package repository

import (
	"fmt"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
)

// метод обновления пароля аккаунта
func (p *PostgreSQL) UpdatePassword(login, hashPassword string) error {
	result := p.DB.Model(&models.User{}).
		Where("login = ?", login).
		Update("hash_password", hashPassword)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// метод создания аккаунта(логин и пароль)
func (p *PostgreSQL) SaveUser(user *models.User) error {
	result := p.DB.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("repository: SaveUser - ошибка сохранения данных.\n ERROR: %v", result.Error)
	}
	return nil
}

// метод удаления аккаунта(логин и пароль)
func (p *PostgreSQL) DelUser(login string) error {
	result := p.DB.Delete(&models.User{}, login)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// метод получения инфо об аккаунте (логин, пароль)
func (p *PostgreSQL) GetUser(login string) (*models.User, error) {
	var user models.User
	result := p.DB.Where("login = ?", login).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("repository: GetUser - ошибка получения данных юзера из БД.\n ERROR: %v", result.Error)
	}
	return &user, nil
}
