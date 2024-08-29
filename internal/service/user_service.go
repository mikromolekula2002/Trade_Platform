package service

import (
	"fmt"
	"net/http"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/utils"
)

// Сервисная логика для домашней страницы, если пользователь не авторизован или у него ошибка с кукой
func (s *UserService) GetAdsData() (*models.HomeData, error) {
	usersAds, err := s.repo.GetAllAds()
	if err != nil {
		return nil, fmt.Errorf("userservice.go: GetAdsData - %v", err)
	}
	// Объединение данных в структуру ProfileData
	data := &models.HomeData{
		UserData: nil,
		UserAds:  usersAds,
	}
	return data, nil
}

// Сервисная логика, получение даты для домашней страницы если юзер авторизован и все норм
func (s *UserService) GetHomeData(cookie *http.Cookie) (*models.HomeData, error) {
	op := "userservice.go: GetHomeData:"

	claims, err := s.jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", op, err)
	}

	User, err := s.repo.GetUserData(claims.UserLogin)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", op, err)
	}

	userAds, err := s.repo.GetAllAds()
	if err != nil {
		return nil, fmt.Errorf("%s - %v", op, err)
	}

	// Объединение данных в структуру ProfileData
	data := &models.HomeData{
		UserData: User,
		UserAds:  userAds,
	}

	return data, nil

}

// Сервисная логика, получение данных аккаунта СКОРЕЕ ВСЕГО НАХУЙ НЕ НУЖНО, БЛЯТЬ
func (s *UserService) GetProfileData(login string) (*models.ProfileData, error) {
	op := "userservice.go: GetProfileData:"

	user, err := s.repo.GetUserData(login)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", op, err)
	}

	userAds, err := s.repo.GetUserAds(login)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", op, err)
	}

	likes, err := s.GetLikedAds(login)
	if err != nil {
		return nil, fmt.Errorf("%s - %v", op, err)
	}

	// Объединение данных в структуру ProfileData
	data := &models.ProfileData{
		UserData: user,
		UserAds:  userAds,
		Likes:    likes,
	}

	return data, nil
}

// Сервисная логика для логики профиля, получение данных аккаунта гостя переходящего на чужой профиль
func (s *UserService) GetGuestData(cookie *http.Cookie) (*models.UserData, error) {
	op := "userservice.go: GetGuestData:"

	claims, err := s.jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return nil, fmt.Errorf("%s \n%v", op, err)
	}

	data, err := s.repo.GetUserData(claims.UserLogin)
	if err != nil {
		return nil, fmt.Errorf("%s \n%v", op, err)
	}

	return data, nil
}

func (s *UserService) GetLikedAds(login string) ([]*models.UserAds, error) {
	op := "userservice.go: GetLikedAds:"

	var adsSlice []*models.UserAds

	likes, err := s.repo.GetAllLikes(login)
	if err != nil {
		return nil, fmt.Errorf("%s \n%v", op, err)
	}

	for i := range likes {
		ads, err := s.repo.GetOneAds(likes[i].Ads_Id)
		if err != nil {
			return nil, fmt.Errorf("%s \n%v", op, err)
		}

		adsSlice = append(adsSlice, ads)
	}
	return adsSlice, nil
}

// Добавление объявления в избранное
func (s *UserService) SaveLikedAd(adsID string, cookie *http.Cookie) error {
	op := "userservice.go: SaveLikedAd:"

	claims, err := s.jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return fmt.Errorf("%s - %v", op, err)
	}

	LikedAd := &models.Likes{
		Ads_Id:     adsID,
		User_Login: claims.UserLogin,
	}

	err = s.repo.SaveLikes(LikedAd)
	if err == utils.ErrAlreadyExist {

		err := s.repo.DelLikes(LikedAd.User_Login, LikedAd.Ads_Id)
		if err != nil {
			return fmt.Errorf("%s - %v", op, err)
		}

	} else if err != nil {
		return fmt.Errorf("%s - %v", op, err)
	}

	return nil
}

// Сервисная логика, получение данных аккаунта
func (s *UserService) GetUserData(cookie *http.Cookie) (*models.UserData, error) {
	op := "userservice.go: GetUserData: "

	claims, err := s.jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return nil, fmt.Errorf("%s \n%v", op, err)
	}

	User, err := s.repo.GetUserData(claims.UserLogin)
	if err != nil {
		return nil, fmt.Errorf("%s \n%v", op, err)
	}
	return User, nil
}

// Сервисная логика, получение данных для урл: /profile/:id (юзера, его объяв, а также лайкнутых)
func (s *UserService) OneAds(cookie *http.Cookie, cookieChecker bool, adsID string) (*models.OneAdsData, string, error) {
	op := "userservice.go: OneAds:"

	// Если owner == true это значит что юзер является владельцем объявления
	var owner bool
	// Заранее объявляем переменную с типом нашей модели
	data := &models.OneAdsData{}

	// Если cookieChecker == true это значит что пользователь авторизован
	switch cookieChecker {
	case true:
		// получаем обьявление
		AdsData, err := s.repo.GetOneAds(adsID)
		if err != nil {
			return nil, "", fmt.Errorf("%s ошибка получения объявления: \n%v", op, err)
		}
		// проверяем является ли юзер владельцем объявления(через куку и UserID в структуре объявления)
		//получаем логин - (user) с куки
		user, ownerChecker, err := s.VerifyAdsOwner(cookie, AdsData)
		if err != nil {
			return nil, "", fmt.Errorf("%s \n%v", op, err)
		}
		// присваиваем результат главной переменной owner
		owner = ownerChecker
		// получаем данные о юзере
		userProfile, err := s.repo.GetUserData(user)
		if err != nil {
			return nil, "", fmt.Errorf("%s \n%v", op, err)
		}
		// получаем данные о профиле владельца объявления
		ownerUserAds, err := s.repo.GetUserData(AdsData.User_Id)
		if err != nil {
			return nil, "", fmt.Errorf("%s \n%v", op, err)
		}
		// записываем все полученное в передаваему на html структуру
		data = &models.OneAdsData{
			OwnerUserData: ownerUserAds,
			UserData:      userProfile,
			Ads:           AdsData,
		}
		// Если cookieChecker == false это значит что пользователь не авторизован
	case false:
		// получаем обьявление
		AdsData, err := s.repo.GetOneAds(adsID)
		if err != nil {
			return nil, "", fmt.Errorf("%s ошибка получения объявления: \n%v", op, err)
		}

		// получаем данные о профиле владельца объявления
		ownerUserAds, err := s.repo.GetUserData(AdsData.User_Id)
		if err != nil {
			return nil, "", fmt.Errorf("%s \n%v", op, err)
		}

		// записываем все полученное в передаваему на html структуру
		data = &models.OneAdsData{
			OwnerUserData: ownerUserAds,
			UserData:      nil,
			Ads:           AdsData,
		}
		owner = false
	}

	// решил сделать так: тут идет проверка на наличие куки(cookieChecker), а также на владельца объявления(owner)
	// в ответах возвращаем название html под каждую ситуацию и структуру с данными
	if cookieChecker && owner {
		return data, "ads-owner.html", nil
	} else if cookieChecker && !owner {
		return data, "ads-no-owner.html", nil
	} else {
		return data, "ads-no-auth.html", nil
	}
}
