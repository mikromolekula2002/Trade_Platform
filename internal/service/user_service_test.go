package service

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mikromolekula2002/Trade_Platform/internal/jwt"
	"github.com/mikromolekula2002/Trade_Platform/internal/mocks"
	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/utils"
)

func TestUserService_GetAdsData(t *testing.T) {
	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для репозитория
	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := UserService{
		repo: mockRepo,
	}

	// Таблица тестов
	tests := []struct {
		name        string
		mockReturn  []*models.UserAds
		mockError   error
		expectedErr string
		expectedAds []*models.UserAds
	}{
		{
			name:        "Успешное получение данных",
			mockReturn:  []*models.UserAds{{Ads_Id: "1"}, {Ads_Id: "2"}},
			mockError:   nil,
			expectedErr: "",
			expectedAds: []*models.UserAds{{Ads_Id: "1"}, {Ads_Id: "2"}},
		},
		{
			name:        "Ошибка при получении данных",
			mockReturn:  nil,
			mockError:   fmt.Errorf("ошибка доступа к данным"),
			expectedErr: "userservice.go: GetAdsData - ошибка доступа к данным",
			expectedAds: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокаем метод GetAllAds
			mockRepo.EXPECT().GetAllAds().Return(tt.mockReturn, tt.mockError)

			// Выполняем тестируемый метод
			_, err := userService.GetAdsData()

			// Проверяем ошибку
			if err != nil {
				if err.Error() != tt.expectedErr {
					t.Fatalf("\nExpected - %v,\n Actual - %v", tt.expectedErr, err)
				}
			} else if tt.expectedErr != "" {
				t.Fatalf("\nExpected - %v,\n Actual - nil", tt.expectedErr)
			}
		})
	}
}

func TestUserService_GetHomeData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	userService := NewUserService(mockRepo, mockJWT, "mockKey")

	testCases := []struct {
		name           string
		cookie         *http.Cookie
		mockJWTToken   *jwt.Claims
		mockJWTError   error
		mockUserData   *models.UserData
		mockUserError  error
		mockAdsData    []*models.UserAds
		mockAdsError   error
		expectedError  string
		expectedUser   string
		expectedAdsLen int
	}{
		{
			name:           "Success",
			cookie:         &http.Cookie{Value: "mockToken"},
			mockJWTToken:   &jwt.Claims{UserLogin: "testuser"},
			mockJWTError:   nil,
			mockUserData:   &models.UserData{Login: "testuser"},
			mockUserError:  nil,
			mockAdsData:    []*models.UserAds{{Ads_Id: "1"}, {Ads_Id: "2"}},
			mockAdsError:   nil,
			expectedError:  "",
			expectedUser:   "testuser",
			expectedAdsLen: 2,
		},
		{
			name:          "ExtractTokenError",
			cookie:        &http.Cookie{Value: "invalidToken"},
			mockJWTToken:  nil,
			mockJWTError:  fmt.Errorf("jwt.go: ExtractToken:"),
			mockUserData:  nil,
			mockUserError: nil,
			mockAdsData:   nil,
			mockAdsError:  nil,
			expectedError: "userservice.go: GetHomeData: - jwt.go: ExtractToken:",
		},
		{
			name:          "GetUserDataError",
			cookie:        &http.Cookie{Value: "mockToken"},
			mockJWTToken:  &jwt.Claims{UserLogin: "testuser"},
			mockJWTError:  nil,
			mockUserData:  nil,
			mockUserError: fmt.Errorf("repository: GetUserData - ошибка получения данных аккаунта из БД."),
			mockAdsData:   nil,
			mockAdsError:  nil,
			expectedError: "userservice.go: GetHomeData: - repository: GetUserData - ошибка получения данных аккаунта из БД.",
		},
		{
			name:          "GetAllAdsError",
			cookie:        &http.Cookie{Value: "mockToken"},
			mockJWTToken:  &jwt.Claims{UserLogin: "testuser"},
			mockJWTError:  nil,
			mockUserData:  &models.UserData{Login: "testuser"},
			mockUserError: nil,
			mockAdsData:   nil,
			mockAdsError:  fmt.Errorf("repository: GetUserData - ошибка получения данных аккаунта из БД."),
			expectedError: "userservice.go: GetHomeData: - repository: GetUserData - ошибка получения данных аккаунта из БД.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Мокаем ExtractToken
			mockJWT.EXPECT().ExtractToken(tc.cookie.Value, "mockKey").Return(tc.mockJWTToken, tc.mockJWTError).Times(1)

			// Мокаем GetUserData и GetAllAds только если предыдущий шаг успешен
			if tc.mockJWTError == nil {
				mockRepo.EXPECT().GetUserData(tc.mockJWTToken.UserLogin).Return(tc.mockUserData, tc.mockUserError).Times(1)
				if tc.mockUserError == nil {
					mockRepo.EXPECT().GetAllAds().Return(tc.mockAdsData, tc.mockAdsError).Times(1)
				}
			}

			// Вызов GetHomeData
			data, err := userService.GetHomeData(tc.cookie)
			if tc.expectedError != "" {
				if err == nil || err.Error() != tc.expectedError {
					t.Fatalf("\nExpected %v \n Actual: %v", tc.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("\nExpected `nil` \n Actual: %v", err)
				}

				if data.UserData.Login != tc.expectedUser || len(data.UserAds) != tc.expectedAdsLen {
					t.Fatalf("\nWrong Data: %v", data)
				}
			}
		})
	}
}

func TestUserService_GetProfileData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)
	userService := UserService{
		repo: mockRepo,
		jwt:  mockJWT,
	}

	mockUser := &models.UserData{Login: "testuser"}
	mockAds := []*models.UserAds{{Ads_Id: "1"}, {Ads_Id: "2"}}
	mockLikes := []*models.UserAds{{Ads_Id: "3"}, {Ads_Id: "4"}}

	// Мокаем методы репозитория
	mockRepo.EXPECT().GetUserData("testuser").Return(mockUser, nil).Times(1)
	mockRepo.EXPECT().GetUserAds("testuser").Return(mockAds, nil).Times(1)
	mockRepo.EXPECT().GetAllLikes("testuser").Return([]*models.Likes{{Ads_Id: "3", User_Login: "testuser"}, {Ads_Id: "4", User_Login: "testuser"}}, nil).Times(1)
	mockRepo.EXPECT().GetOneAds("3").Return(mockLikes[0], nil).Times(1)
	mockRepo.EXPECT().GetOneAds("4").Return(mockLikes[1], nil).Times(1)

	// Выполняем тестируемый метод
	data, err := userService.GetProfileData("testuser")
	if err != nil {
		t.Fatalf("\n Expected - nil,\n Actual - %v", err)
	}

	// Проверяем результат
	if data.UserData.Login != "testuser" || len(data.UserAds) != 2 || len(data.Likes) != 2 {
		t.Fatalf("\n Wrong Data %v", data)
	}
}

func TestUserService_GetGuestData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	userService := UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "test_jwt_key",
	}

	// Создаем тестовую таблицу
	tests := []struct {
		name          string
		cookie        *http.Cookie
		mockJWTSetup  func()
		mockRepoSetup func()
		expectedData  *models.UserData
		expectedError error
	}{
		{
			name:   "Success getting guest data",
			cookie: &http.Cookie{Value: "valid_token"},
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("valid_token", "test_jwt_key").Return(&jwt.Claims{UserLogin: "guestuser"}, nil).Times(1)
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().GetUserData("guestuser").Return(&models.UserData{Login: "guestuser"}, nil).Times(1)
			},
			expectedData:  &models.UserData{Login: "guestuser"},
			expectedError: nil,
		},
		{
			name:   "Error extracting token",
			cookie: &http.Cookie{Value: "invalid_token"},
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("invalid_token", "test_jwt_key").Return(nil, fmt.Errorf("invalid token")).Times(1)
			},
			mockRepoSetup: func() {},
			expectedData:  nil,
			expectedError: fmt.Errorf("userservice.go: GetGuestData: \ninvalid token"),
		},
		{
			name:   "Error getting guest user data",
			cookie: &http.Cookie{Value: "valid_token"},
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("valid_token", "test_jwt_key").Return(&jwt.Claims{UserLogin: "guestuser"}, nil).Times(1)
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().GetUserData("guestuser").Return(nil, fmt.Errorf("database error")).Times(1)
			},
			expectedData:  nil,
			expectedError: fmt.Errorf("userservice.go: GetGuestData: \ndatabase error"),
		},
	}

	// Проходим по тестам
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockJWTSetup()
			tt.mockRepoSetup()

			data, err := userService.GetGuestData(tt.cookie)

			// Проверка ошибок
			if err == nil && tt.expectedError != nil {
				t.Fatalf("\n Expected - %v,\n Actual - nil", tt.expectedError)
			}

			if err != nil && tt.expectedError == nil {
				t.Fatalf("\n Expected - %v", err)
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Fatalf("\n Expected - %v,\n Actual - %v", tt.expectedError, err)
			}

			// Проверка данных
			if !reflect.DeepEqual(data, tt.expectedData) {
				t.Fatalf("\n Expected - %v,\n Actual - %v", tt.expectedData, data)
			}
		})
	}
}

func TestUserService_SaveLikedAd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	userService := UserService{
		repo:   mockRepo,
		jwt:    mockJWT,
		jwtKey: "mockKey",
	}

	// Создаем тестовую таблицу
	tests := []struct {
		name          string
		cookie        *http.Cookie
		adsID         string
		mockJWTSetup  func()
		mockRepoSetup func()
		expectedError error
	}{
		{
			name:   "successfully saves liked ad",
			cookie: &http.Cookie{Value: "validToken"},
			adsID:  "123",
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("validToken", "mockKey").Return(&jwt.Claims{UserLogin: "testuser"}, nil).Times(1)
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().SaveLikes(gomock.Any()).Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name:   "ad already liked, deleting like",
			cookie: &http.Cookie{Value: "validToken"},
			adsID:  "123",
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("validToken", "mockKey").Return(&jwt.Claims{UserLogin: "testuser"}, nil).Times(1)
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().SaveLikes(gomock.Any()).Return(utils.ErrAlreadyExist).Times(1)
				mockRepo.EXPECT().DelLikes("testuser", "123").Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name:   "error extracting token",
			cookie: &http.Cookie{Value: "invalidToken"},
			adsID:  "123",
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("invalidToken", "mockKey").Return(nil, fmt.Errorf("invalid token")).Times(1)
			},
			mockRepoSetup: func() {
				// Не ожидаем вызова репозитория, так как ошибка на этапе токена
			},
			expectedError: fmt.Errorf("userservice.go: SaveLikedAd: - invalid token"),
		},
		{
			name:   "error saving like",
			cookie: &http.Cookie{Value: "validToken"},
			adsID:  "123",
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("validToken", "mockKey").Return(&jwt.Claims{UserLogin: "testuser"}, nil).Times(1)
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().SaveLikes(gomock.Any()).Return(fmt.Errorf("save error")).Times(1)
			},
			expectedError: fmt.Errorf("userservice.go: SaveLikedAd: - save error"),
		},
		{
			name:   "error deleting like",
			cookie: &http.Cookie{Value: "validToken"},
			adsID:  "123",
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("validToken", "mockKey").Return(&jwt.Claims{UserLogin: "testuser"}, nil).Times(1)
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().SaveLikes(gomock.Any()).Return(utils.ErrAlreadyExist).Times(1)
				mockRepo.EXPECT().DelLikes("testuser", "123").Return(fmt.Errorf("delete error")).Times(1)
			},
			expectedError: fmt.Errorf("userservice.go: SaveLikedAd: - delete error"),
		},
	}

	// Проходим по тестам
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockJWTSetup()
			tt.mockRepoSetup()

			err := userService.SaveLikedAd(tt.adsID, tt.cookie)

			if err == nil && tt.expectedError != nil {
				t.Fatalf("\n Expected - %v,\n Actual - nil", tt.expectedError)
			}

			if err != nil && tt.expectedError == nil {
				t.Fatalf("\n Expected - nil,\n Actual - %v", err)
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Fatalf("\n Expected - %v,\n Actual - %v", tt.expectedError, err)
			}
		})
	}
}

func TestUserService_GetLikedAds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userService := UserService{
		repo: mockRepo,
	}

	// Создаем тестовую таблицу
	tests := []struct {
		name          string
		login         string
		mockRepoSetup func()
		expectedAds   []*models.UserAds
		expectedError error
	}{
		{
			name:  "Success getting liked ads",
			login: "user1",
			mockRepoSetup: func() {
				likes := []*models.Likes{
					{Ads_Id: "ads1"},
					{Ads_Id: "ads2"},
				}
				mockRepo.EXPECT().GetAllLikes("user1").Return(likes, nil).Times(1)
				mockRepo.EXPECT().GetOneAds("ads1").Return(&models.UserAds{Ads_Id: "ads1"}, nil).Times(1)
				mockRepo.EXPECT().GetOneAds("ads2").Return(&models.UserAds{Ads_Id: "ads2"}, nil).Times(1)
			},
			expectedAds: []*models.UserAds{
				{Ads_Id: "ads1"},
				{Ads_Id: "ads2"},
			},
			expectedError: nil,
		},
		{
			name:  "Error getting likes",
			login: "user1",
			mockRepoSetup: func() {
				mockRepo.EXPECT().GetAllLikes("user1").Return(nil, fmt.Errorf("database error")).Times(1)
			},
			expectedAds:   nil,
			expectedError: fmt.Errorf("userservice.go: GetLikedAds: \ndatabase error"),
		},
		{
			name:  "Error getting one ad",
			login: "user1",
			mockRepoSetup: func() {
				likes := []*models.Likes{
					{Ads_Id: "ads1"},
					{Ads_Id: "ads2"},
				}
				mockRepo.EXPECT().GetAllLikes("user1").Return(likes, nil).Times(1)
				mockRepo.EXPECT().GetOneAds("ads1").Return(&models.UserAds{Ads_Id: "ads1"}, nil).Times(1)
				mockRepo.EXPECT().GetOneAds("ads2").Return(nil, fmt.Errorf("ad not found")).Times(1)
			},
			expectedAds:   nil,
			expectedError: fmt.Errorf("userservice.go: GetLikedAds: \nad not found"),
		},
	}

	// Проходим по тестам
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()

			ads, err := userService.GetLikedAds(tt.login)

			// Проверяем ошибки
			if err == nil && tt.expectedError != nil {
				t.Fatalf("\nExpected - %v, \n Actual - nil", tt.expectedError)
			}

			if err != nil && tt.expectedError == nil {
				t.Fatalf("\nExpected - nil, \n Actual - %v", err)
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Fatalf("\nExpected - %v, \n Actual - %v", tt.expectedError, err)
			}

			// Проверяем полученные данные
			if !reflect.DeepEqual(ads, tt.expectedAds) {
				t.Fatalf("\nExpected - %v, \n Actual - %v", tt.expectedAds, ads)
			}
		})
	}
}

func TestUserService_OneAds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockJWT := mocks.NewMockJWTService(ctrl)

	userService := UserService{
		repo: mockRepo,
		jwt:  mockJWT,
	}

	// Создаем тестовую таблицу
	tests := []struct {
		name          string
		cookie        *http.Cookie
		cookieChecker bool
		adsID         string
		mockRepoSetup func()
		mockJWTSetup  func()
		expectedData  *models.OneAdsData
		expectedHTML  string
		expectedError error
	}{
		{
			name:          "Authorized user, owner of the ad",
			cookie:        &http.Cookie{Value: "validToken"},
			cookieChecker: true,
			adsID:         "ads1",
			mockRepoSetup: func() {
				adsData := &models.UserAds{Ads_Id: "ads1", User_Id: "user1"}
				mockRepo.EXPECT().GetOneAds("ads1").Return(adsData, nil).Times(1)
				mockRepo.EXPECT().GetUserData("user1").Return(&models.UserData{Login: "user1"}, nil).Times(1)
				mockRepo.EXPECT().GetUserData("user1").Return(&models.UserData{Login: "user1"}, nil).Times(1)
			},
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("validToken", gomock.Any()).Return(&jwt.Claims{UserLogin: "user1"}, nil).Times(1)
			},
			expectedData: &models.OneAdsData{
				OwnerUserData: &models.UserData{Login: "user1"},
				UserData:      &models.UserData{Login: "user1"},
				Ads:           &models.UserAds{Ads_Id: "ads1", User_Id: "user1"},
			},
			expectedHTML:  "ads-owner.html",
			expectedError: nil,
		},
		{
			name:          "Authorized user, not the owner of the ad",
			cookie:        &http.Cookie{Value: "validToken"},
			cookieChecker: true,
			adsID:         "ads1",
			mockRepoSetup: func() {
				adsData := &models.UserAds{Ads_Id: "ads1", User_Id: "user1"}
				mockRepo.EXPECT().GetOneAds("ads1").Return(adsData, nil).Times(1)
				mockRepo.EXPECT().GetUserData("user1").Return(&models.UserData{Login: "user1"}, nil).Times(1)
				mockRepo.EXPECT().GetUserData("testuser").Return(&models.UserData{Login: "testuser"}, nil).Times(1)
			},
			mockJWTSetup: func() {
				mockJWT.EXPECT().ExtractToken("validToken", gomock.Any()).Return(&jwt.Claims{UserLogin: "testuser"}, nil).Times(1)
			},
			expectedData: &models.OneAdsData{
				OwnerUserData: &models.UserData{Login: "user1"},
				UserData:      &models.UserData{Login: "testuser"},
				Ads:           &models.UserAds{Ads_Id: "ads1", User_Id: "user1"},
			},
			expectedHTML:  "ads-no-owner.html",
			expectedError: nil,
		},
		{
			name:          "Unauthorized user",
			cookie:        nil,
			cookieChecker: false,
			adsID:         "ads1",
			mockRepoSetup: func() {
				adsData := &models.UserAds{Ads_Id: "ads1", User_Id: "user1"}
				mockRepo.EXPECT().GetOneAds("ads1").Return(adsData, nil).Times(1)
				mockRepo.EXPECT().GetUserData("user1").Return(&models.UserData{Login: "user1"}, nil).Times(1)
			},
			mockJWTSetup: func() {},
			expectedData: &models.OneAdsData{
				OwnerUserData: &models.UserData{Login: "user1"},
				UserData:      nil,
				Ads:           &models.UserAds{Ads_Id: "ads1", User_Id: "user1"},
			},
			expectedHTML:  "ads-no-auth.html",
			expectedError: nil,
		},
		{
			name:          "Error getting ad data",
			cookie:        nil,
			cookieChecker: false,
			adsID:         "ads1",
			mockRepoSetup: func() {
				mockRepo.EXPECT().GetOneAds("ads1").Return(nil, fmt.Errorf("ad not found")).Times(1)
			},
			mockJWTSetup:  func() {},
			expectedData:  nil,
			expectedHTML:  "",
			expectedError: fmt.Errorf("userservice.go: OneAds: ошибка получения объявления: \nad not found"),
		},
	}

	// Проходим по тестам
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockJWTSetup()
			tt.mockRepoSetup()

			data, html, err := userService.OneAds(tt.cookie, tt.cookieChecker, tt.adsID)

			// Проверяем ошибки
			if err == nil && tt.expectedError != nil {
				t.Fatalf("\nExpected - %v, \n Actual - nil", tt.expectedError)
			}

			if err != nil && tt.expectedError == nil {
				t.Fatalf("\nExpected - `nil`, \n Actual - %v", err)
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Fatalf("\nExpected - %v, \n Actual - %v", tt.expectedError, err)
			}

			// Проверяем данные и шаблон
			if !reflect.DeepEqual(data, tt.expectedData) {
				t.Fatalf("\nExpected - %v, \n Actual - %v", tt.expectedData, data)
			}

			if html != tt.expectedHTML {
				t.Fatalf("\nExpected - %v, \n Actual - %v", tt.expectedHTML, html)
			}
		})
	}
}
