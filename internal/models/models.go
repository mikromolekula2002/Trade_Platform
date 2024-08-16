package models

//БД профиль
type UserData struct {
	Login        string `gorm:"primaryKey"`
	Name         string
	FirstName    string
	PhoneNumber  string
	ImageProfile string
}

//БД пользователь
type User struct {
	Login        string `gorm:"primaryKey"`
	HashPassword string
}

//БД объявления
type UserAds struct {
	Ads_Id          string `gorm:"primaryKey"`
	User_Id         string
	Image_1         string
	Image_2         string
	Image_3         string
	Ads_Name        string
	Ads_Description string
	Ads_Price       float64
	Ads_Likes       int `gorm:"default:0"`
}

// вывод данных о избранном
type Likes struct {
	User_Login string `gorm:"primaryKey"`
	Ads_Id     string
}

//структура для выводы данных профиля
type ProfileData struct {
	GuestData *UserData
	UserData  *UserData
	UserAds   []*UserAds
	Likes     []*UserAds
}

//структура для вывода данных на домашней странице
type HomeData struct {
	UserData *UserData
	UserAds  []*UserAds
}

//структура для вывода данных об объявлении
type OneAdsData struct {
	OwnerUserData *UserData
	UserData      *UserData
	Ads           *UserAds
}
