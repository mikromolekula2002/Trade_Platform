package service

import (
	"testing"

	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateUsername(t *testing.T) {
	us := &UserService{}

	tests := []struct {
		login  string
		errMsg string
	}{
		{"", "логин не может быть пустым"},
		{"test@user", "логин не может содержать спец. символы"},
		{"verylongusername12345", "логин не может превышать длину в 20 символов"},
		{"validuser", ""},
	}

	for _, tt := range tests {
		err := us.ValidateUsername(tt.login, "логин")
		if tt.errMsg == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tt.errMsg)
		}
	}
}

func TestValidateUserData(t *testing.T) {
	us := &UserService{}
	data := &models.UserData{
		PhoneNumber:  "invalidphone",
		Name:         "ValidName",
		FirstName:    "ValidFirstName",
		ImageProfile: "",
	}

	err := us.ValidateUserData(data)
	assert.Error(t, err)
}

func TestValidatePhoneNumber(t *testing.T) {
	us := &UserService{}

	tests := []struct {
		phone  string
		errMsg string
	}{
		{"", "service: validation.go/validatePhoneNumber: некоректный номер телефона\nERROR: <nil>"},
		{"+1234567890", ""},
		{"+1234567895378525250", "service: validation.go/validatePhoneNumber: некоректный номер телефона\nERROR: <nil>"},
		{"12345678529", "service: validation.go/validatePhoneNumber: некоректный номер телефона\nERROR: <nil>"},
	}

	for _, tt := range tests {
		err := us.ValidatePhoneNumber(tt.phone)
		if tt.errMsg == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tt.errMsg)
		}
	}
}

func TestValidateName(t *testing.T) {
	us := UserService{}
	name := "ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"

	tests := []struct {
		name      string
		firstname string
		errMsg    string
	}{
		{"", "", "service: validation.go/validateUserData: поля ФИО не могут быть пустыми"},
		{"ValidName", "ValidFirstName", ""},
		{name, "asdg", "service: validation.go/validateUserData: поля ФИО не могут быть больше 50 символов"},
		{".!fsdfr", "asfr", "service: validation.go/validateUserData: поля ФИО не могут содержать спец. символы"},
	}

	for _, tt := range tests {
		err := us.ValidateName(tt.name, tt.firstname)
		if tt.errMsg == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tt.errMsg)
		}
	}
}

func TestValidateUserAds(t *testing.T) {
	us := UserService{}

	testTable := []struct {
		data   *models.UserAds
		errMsg string
	}{
		{
			data: &models.UserAds{
				User_Id:         "user1",
				Image_1:         "Valid_Image1.jpg",
				Image_2:         "Valid_Image2.jpg",
				Image_3:         "Valid_Image3.jpg",
				Ads_Name:        "",
				Ads_Description: "",
				Ads_Price:       100.0,
			},
			errMsg: "service: validation.go/validateUserData: поля `название` и `описание` не могут быть пустыми",
		},
		{
			data: &models.UserAds{
				User_Id:         "user2",
				Image_1:         "too-long-image-url-that-exceeds-the-limitssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss.jpg",
				Image_2:         "Valid_Image2.jpg",
				Image_3:         "Valid_Image3.jpg",
				Ads_Name:        "ValidName",
				Ads_Description: "Valid Description",
				Ads_Price:       200.0,
			},
			errMsg: "service: validation.go/validateUserData: Фото1 объявления не может быть больше 255 символов",
		},
		{
			data: &models.UserAds{
				User_Id:         "user1",
				Image_1:         "Valid_Image1.jpg",
				Image_2:         "Valid_Image2.jpg",
				Image_3:         "Valid_Image3.jpg",
				Ads_Name:        "too long Ads Name that exceeds the limitsssssssssssssssssssssssssssssssssssssssssssssssssssss",
				Ads_Description: "Valid Description",
				Ads_Price:       100.0,
			},
			errMsg: "service: validation.go/validateUserData: Название объявления не может быть больше 50 символов",
		},
		{
			data: &models.UserAds{
				User_Id:         "user1",
				Image_1:         "Valid_Image1.jpg",
				Image_2:         "Valid_Image2.jpg",
				Image_3:         "Valid_Image3.jpg",
				Ads_Name:        "ValidName",
				Ads_Description: "Description is more than 200 simbols aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Ads_Price:       100.0,
			},
			errMsg: "service: validation.go/validateUserData: Описание объявления не может быть больше 200 символов",
		},
		{
			data: &models.UserAds{
				User_Id:         "user1",
				Image_1:         "Valid_Image1.jpg",
				Image_2:         "Valid_Image2.jpg",
				Image_3:         "Valid_Image3.jpg",
				Ads_Name:        "ValidName",
				Ads_Description: "Valid Description",
				Ads_Price:       100.0,
			},
			errMsg: "",
		},
	}

	for _, tt := range testTable {
		err := us.ValidateUserAds(tt.data)
		if tt.errMsg == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tt.errMsg)
		}
	}
}
