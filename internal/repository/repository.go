package repository

import (
	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	DB *gorm.DB
}

func PostgreInit(db *gorm.DB) *PostgreSQL {
	datab := PostgreSQL{
		DB: db,
	}
	return &datab
}

// UserRepository описывает методы работы с данными пользователей.
type UserRepository interface {
	// методы для работы с данными пользователей
	UpdateUserData(userData *models.UserData) error
	SaveUserData(userData *models.UserData) error
	DelUserData(login string) error
	GetUserData(login string) (*models.UserData, error)

	// методы для работы с учетными записями пользователей
	UpdatePassword(login, hashPassword string) error
	SaveUser(user *models.User) error
	DelUser(login string) error
	GetUser(login string) (*models.User, error)

	// методы для работы с объявлениями пользователей
	UpdateUserAds(userAds *models.UserAds) error
	SaveUserAds(userAds *models.UserAds) (string, error)
	DelUserAds(userID string, adsID string) error
	GetOneAds(adsID string) (*models.UserAds, error)
	GetUserAds(login string) ([]*models.UserAds, error)
	GetAllAds() ([]*models.UserAds, error)

	// методы для работы с лайками
	SaveLikes(likes *models.Likes) error
	DelLikes(userlogin string, adsId string) error
	GetAllLikes(userLogin string) ([]*models.Likes, error)
}
