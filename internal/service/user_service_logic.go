package service

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/mikromolekula2002/Trade_Platform/internal/jwt"
	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/repository"
	"github.com/mikromolekula2002/Trade_Platform/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// Инит нашей структуры с базой данных
func NewUserService(repo repository.PostgreSQL, jwtkey string) *UserService {
	return &UserService{
		repo:   repo,
		jwtKey: jwtkey,
	}
}

// Сервисная логика, сохранение данных пользователя в БД
func (s *UserService) RegisterUser(login, password string) (string, error) {

	if err := s.ValidateUsername(login, "логин"); err != nil {
		return "", err
	}

	if err := s.ValidateUsername(password, "пароль"); err != nil {
		return "", err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("userservice.go: RegisterUser - ошибка хеширования пароля.\nERROR: %v", err)
	}

	user := models.User{
		Login:        login,
		HashPassword: string(hashPassword),
	}

	err = s.repo.SaveUser(&user)
	if err != nil {
		return "", fmt.Errorf("userservice.go: RegisterUser\n%v", err)
	}
	userData := models.UserData{
		Login:        login,
		Name:         "",
		FirstName:    "",
		PhoneNumber:  "",
		ImageProfile: "/images/stock/avatar-mane.png",
	}

	err = s.repo.SaveUserData(&userData)
	if err != nil {
		return "", fmt.Errorf("userservice.go: RegisterUser\n%v", err)
	}

	jwtToken, err := jwt.CreateToken(user.Login, s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("userservice.go: RegisterUser\n%v", err)
	}
	return jwtToken, nil
}

// Сервисная логика, аутентификация юзера(сверка паролей бд и введенного)
func (s *UserService) AuthUser(login, password string) (string, error) {
	user, err := s.repo.GetUser(login)
	if err != nil {
		return "", fmt.Errorf("userservice.go: AuthUser\n%v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("userservice.go: AuthUser - ошибка при проверке пароля.\nERROR: %v", err)
	}

	jwtToken, err := jwt.CreateToken(user.Login, s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("userservice.go: AuthUser\n%v", err)
	}
	return jwtToken, nil
}

// Сервисная логика, выход из аккаунта (удаление куки)
func (s *UserService) QuitFromAccount(cookie *http.Cookie) (string, error) {
	claims, err := jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("userservice.go: QuitFromAccount: \n%v", err)
	}

	jwtToken, err := jwt.DeleteToken(claims.UserLogin, s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("userservice.go: QuitFromAccount: \n%v", err)
	}

	return jwtToken, nil
}

// Сервисная логика, сохранение данных аккаунта
func (s *UserService) SaveUserData(data *models.UserData) error {
	err := s.ValidateUserData(data)
	if err != nil {
		return fmt.Errorf("userservice.go: SaveUserData\n%v", err)
	}
	err = s.repo.SaveUserData(data)
	if err != nil {
		return fmt.Errorf("userservice.go: SaveUserData\n%v", err)
	}
	return nil
}

// Метод сделан для проверки является ли юзер перешедший владельцем профиля и отдает true если это так
func (s *UserService) VerifyUserOwner(cookie *http.Cookie, loginParam string) (bool, error) {
	claims, err := jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return false, fmt.Errorf("userservice.go: VerifyUserOwner- ошибка с парсингом куки: \n%v", err)
	}

	switch {
	case loginParam == claims.UserLogin:
		return true, nil
	default:
		return false, nil
	}
}

// Сервисная логика, обновление данных аккаунта пока что нахуй не сдалось
func (s *UserService) UpdateUserData(cookie *http.Cookie, name, firstname, phonenumber string, images []*multipart.FileHeader, missingImage bool, oldImage string) (string, error) {
	var oldImages []string
	oldImages = append(oldImages, oldImage)

	claims, err := jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("userservice.go: VerifyUserOwner- ошибка с парсингом куки: \n%v", err)
	}
	data := &models.UserData{}
	switch missingImage {
	case true:
		data = &models.UserData{
			Login:       claims.UserLogin,
			Name:        name,
			FirstName:   firstname,
			PhoneNumber: phonenumber,
		}

	case false:
		pathsToImages, err := s.DownloadFiles(images, "avatars")
		if err != nil {
			return "", fmt.Errorf("userservice.go: SaveUserAds:\n%v", err)
		}
		data = &models.UserData{
			Login:        claims.UserLogin,
			Name:         name,
			FirstName:    firstname,
			PhoneNumber:  phonenumber,
			ImageProfile: pathsToImages[0],
		}
		err = s.DeleteImages(oldImages)
		if err != nil {
			return "", fmt.Errorf("userservice.go: SaveUserAds:\n%v", err)
		}
	}

	err = s.ValidateUserData(data)
	if err != nil {
		return "", fmt.Errorf("userservice.go: UpdateUserData:\n%v", err)
	}

	err = s.repo.UpdateUserData(data)
	if err != nil {
		return "", fmt.Errorf("userservice.go: UpdateUserData:\n%v", err)
	}
	return claims.UserLogin, nil
}

// Сервисная логика, сохранение данных объявления в БД
func (s *UserService) VerifyAdsOwner(cookie *http.Cookie, UserAds *models.UserAds) (string, bool, error) {
	// получаем токен и его данные
	claims, err := jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return "", false, fmt.Errorf("userservice.go: VerifyAdsOwner - ошибка с парсингом куки: \n%v", err)
	}

	// сверяем полученный идентификатор из токена с идентификатором владельца объявления
	if UserAds.User_Id == claims.UserLogin {
		return claims.UserLogin, true, nil
	} else {
		return claims.UserLogin, false, nil
	}
}

// Сервисная логика, сохранение данных объявления в БД
func (s *UserService) SaveUserAds(cookie *http.Cookie, images []*multipart.FileHeader, adsName, adsDescription string, price float64) (string, error) {
	claims, err := jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return "", fmt.Errorf("userservice.go: SaveUserAds- ошибка с парсингом куки: \n%v", err)
	}

	//Создание уникального идентификатора объявления (длиной на 36 символjd, в целях не повторения)
	SaveAdsID := uuid.New().String()

	pathsToImages, err := s.DownloadFiles(images, "ads")
	if err != nil {
		return "", fmt.Errorf("userservice.go: SaveUserAds:\n%v", err)
	}

	data := &models.UserAds{
		Ads_Id:          SaveAdsID,
		User_Id:         claims.UserLogin,
		Image_1:         pathsToImages[0],
		Image_2:         pathsToImages[1],
		Image_3:         pathsToImages[2],
		Ads_Name:        adsName,
		Ads_Description: adsDescription,
		Ads_Price:       price,
	}

	err = s.ValidateUserAds(data)
	if err != nil {
		return "", fmt.Errorf("userservice.go: SaveUserAds:\n%v", err)
	}

	adsID, err := s.repo.SaveUserAds(data)
	if err != nil {
		return "", fmt.Errorf("userservice.go: SaveUserAds:\n%v", err)
	}
	return adsID, nil
}

// Сервисная логика, обновление данных объявления в БД
func (s *UserService) UpdateUserAds(cookie *http.Cookie, userID, adsID, adsName, adsDescription string, price float64, images []*multipart.FileHeader, oldImages []string) error {
	// получаем данные из токена(идентификатор пользователя)
	claims, err := jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return fmt.Errorf("userservice.go: UpdateUserAds - ошибка с парсингом куки: \n%v", err)
	}
	// сверяем идентификатор из токена с идентификатором владельца объявления
	if claims.UserLogin != userID {
		return fmt.Errorf("userservice.go: UpdateUserAds - ошибка, юзер пытается изменить объявление, не являясь его владельцем")
	}

	data := &models.UserAds{
		Ads_Id:          adsID,
		User_Id:         claims.UserLogin,
		Ads_Name:        adsName,
		Ads_Description: adsDescription,
		Ads_Price:       price,
	}

	pathsToImages, err := s.DownloadFiles(images, "ads")
	if err != nil {
		return fmt.Errorf("userservice.go: UpdateUserAds:\n%v", err)
	}

	var delImages []string

	for i := range pathsToImages {
		if oldImages[i] != pathsToImages[i] && pathsToImages[i] != "" {
			delImages = append(delImages, oldImages[i])
		}
	}

	s.DeleteImages(delImages)

	// Проверяем какие фото были обновлены, а какие изменены и вписываем в структуру
	if len(pathsToImages) > 0 {
		if pathsToImages[0] == "" {
			data.Image_1 = oldImages[0]
		}

		data.Image_1 = pathsToImages[0]
	}

	if len(pathsToImages) > 1 {
		if pathsToImages[1] == "" {
			data.Image_2 = oldImages[1]
		}

		data.Image_2 = pathsToImages[1]
	}

	if len(pathsToImages) > 2 {
		if pathsToImages[2] == "" {
			data.Image_3 = oldImages[2]
		}

		data.Image_3 = pathsToImages[2]
	}

	// проводим валидацию данных
	if err := s.ValidateUserAds(data); err != nil {
		return fmt.Errorf("userservice.go: UpdateUserAds - ошибка валидации: \n%v", err)
	}

	//апдейтим данные в БД
	if err := s.repo.UpdateUserAds(data); err != nil {
		return fmt.Errorf("userservice.go: UpdateUserAds - ошибка обновления БД: \n%v", err)
	}

	return nil
}

// Сервисная логика, удаление существующего объявления
func (s *UserService) DeleteAds(cookie *http.Cookie, adsID string) (bool, error) {
	// получаем данные из токена(идентификатор)
	claims, err := jwt.ExtractToken(cookie.Value, s.jwtKey)
	if err != nil {
		return true, fmt.Errorf("userservice.go: DeleteAds: \n%v", err)
	}

	// удаляем запись по айди объявления, а также идентификатора пользователя
	// если они не совпадают с записью, возвращаем ошибку
	if err := s.repo.DelUserAds(claims.UserLogin, adsID); err != nil {
		if err == utils.ErrNoOwnerAds {
			return false, fmt.Errorf("userservice.go: DeleteAds: \n%v", err)
		} else {
			return true, fmt.Errorf("userservice.go: DeleteAds: \n%v", err)
		}
	}
	return true, nil
}
