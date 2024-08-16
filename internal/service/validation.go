package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
)

const specialChars = `!@#$%^&*()_+-={}|\:;"'<>,.?/~`

// Валидация логина или пароля пользователя, в value указать что проверяется
func (s *UserService) ValidateUsername(login, value string) error {

	if login == "" {
		return fmt.Errorf("%s не может быть пустым", value)
	}

	if strings.ContainsAny(login, specialChars) {
		return fmt.Errorf("%s не может содержать спец. символы", value)
	}

	if len(login) > 20 {
		return fmt.Errorf("%s не может превышать длину в 20 символов", value)
	}
	return nil
}

// Общая валидация профиля пользователя
func (s *UserService) ValidateUserData(user *models.UserData) error {
	//Проверка на наличие загруженной аватарки
	if user.ImageProfile == "" {
		user.ImageProfile = "/images/stock/avatar-mane.png"
	}

	if len(user.ImageProfile) > 50 {
		return fmt.Errorf("service: validation.go/validateUserData: ссылка на фото не должна превышать 50 символов")
	}

	//Проверка номера телефона
	if err := s.ValidatePhoneNumber(user.PhoneNumber); err != nil {
		return err
	}
	//Проверка имени и фамилии
	if err := s.ValidateName(user.Name, user.FirstName); err != nil {
		return err
	}

	return nil
}

// Валидация мобильного номера пользователя
func (s *UserService) ValidatePhoneNumber(phone string) error {
	pattern := `^\+[1-9]\d{1,14}$`

	if matched, err := regexp.MatchString(pattern, phone); err != nil {
		return fmt.Errorf("service: validation.go/validatePhoneNumber: ошибка валидации номера телефона\nERROR: %v", err)

	} else if !matched {
		return fmt.Errorf("service: validation.go/validatePhoneNumber: некоректный номер телефона\nERROR: %v", err)
	}

	return nil
}

// Валидация ФИО пользователя
func (s *UserService) ValidateName(name, firstname string) error {
	//Проверка имени и фамилии
	if name == "" || firstname == "" {
		return fmt.Errorf("service: validation.go/validateUserData: поля ФИО не могут быть пустыми")
	} else if strings.ContainsAny(name, specialChars) || strings.ContainsAny(firstname, specialChars) {
		return fmt.Errorf("service: validation.go/validateUserData: поля ФИО не могут содержать спец. символы")
	} else if len(name) > 50 || len(firstname) > 50 {
		return fmt.Errorf("service: validation.go/validateUserData: поля ФИО не могут быть больше 50 символов")
	}
	return nil
}

// Общая валидация объявления пользователя
func (s *UserService) ValidateUserAds(user *models.UserAds) error {
	//Проверка имени и фамилии
	if user.Ads_Name == "" || user.Ads_Description == "" {
		return fmt.Errorf("service: validation.go/validateUserData: поля `название` и `описание` не могут быть пустыми")
	} else if len(user.Ads_Name) > 50 {
		return fmt.Errorf("service: validation.go/validateUserData: Название объявления не может быть больше 50 символов")
	} else if len(user.Ads_Description) > 200 {
		return fmt.Errorf("service: validation.go/validateUserData: Описание объявления не может быть больше 200 символов")
	}

	if err := s.ValidateAdsImage(user.Image_1, user.Image_2, user.Image_3); err != nil {
		return err
	}

	return nil
}

// Валидация фото объявления по длине ссылки на фото
func (s *UserService) ValidateAdsImage(img1, img2, img3 string) error {
	if len(img1) > 255 {
		return fmt.Errorf("service: validation.go/validateUserData: Фото1 объявления не может быть больше 255 символов")
	}

	if len(img2) > 255 {
		return fmt.Errorf("service: validation.go/validateUserData: Фото1 объявления не может быть больше 255 символов")
	}

	if len(img3) > 255 {
		return fmt.Errorf("service: validation.go/validateUserData: Фото1 объявления не может быть больше 255 символов")
	}

	return nil
}
