package handler

import (
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mikromolekula2002/Trade_Platform/internal/config"
	"github.com/mikromolekula2002/Trade_Platform/internal/models"
	"github.com/mikromolekula2002/Trade_Platform/internal/service"
	"github.com/mikromolekula2002/Trade_Platform/internal/utils"
	"github.com/sirupsen/logrus"
)

// Мейн структура через которую вы работаем, здесь и фреймворк, и БД, и Редис, и Логер
type Handler struct {
	Log     *logrus.Logger
	Echo    *echo.Echo
	UserSvc *service.UserService
}

// Инициализация хендлера с логером и фреймворком
func Init(logger *logrus.Logger, cfg *config.Config, userSvc *service.UserService) *Handler {
	echo := echo.New()
	echo.HideBanner = true

	echo.Server = &http.Server{
		Addr:         cfg.Server.Port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	return &Handler{
		Log:     logger,
		Echo:    echo,
		UserSvc: userSvc,
	}
}

// Определение маршрутов обработчика
func (h *Handler) Routes() {
	h.Echo.GET("/home", h.Home)
	h.Echo.GET("/profile/:id", h.Profile)      //через куку(И НЕТ, нужно через params)
	h.Echo.GET("/profile/form", h.ProfileForm) //редактирование профиля(форма для редакта)
	h.Echo.POST("/quit", h.QuitProfile)
	h.Echo.GET("/auth", h.Authorization)
	h.Echo.POST("/auth/login", h.Login)
	h.Echo.POST("/auth/register", h.Register)
	h.Echo.GET("/ads/:id", h.OneAds)        // определенная объява
	h.Echo.GET("/home/ads-form", h.AdsForm) //создание объявления(форма для загрузки объявы)

	//	h.Echo.POST("/profile/update", h.UpdateProfile)
	//	h.Echo.GET("/profile/likes:id", h.Likes)
	//	h.Echo.POST("/ads/post", h.PostAds)
	//Добавить маршруты /ads/liked и /ads/viewed

	// Приватные маршруты (добавить сюда редактирование профиля, добавление объявления)
	profileGroup := h.Echo.Group("/profile", h.JWTMiddleware)
	profileGroup.POST("/update", h.UpdateProfile)

	AdsGroup := h.Echo.Group("/ads", h.JWTMiddleware)
	AdsGroup.POST("/post", h.PostAds)
	AdsGroup.POST("/like/:id", h.PostLike)
	AdsGroup.POST("/delete/:id", h.DelAds)
	AdsGroup.GET("/form_update/:id", h.AdsUpdateForm)
	AdsGroup.POST("/update", h.UpdateUserAds)
	//Прописать мидлварь авторизации
	//А также прописать мидлварь подсчета уникальных просмотров объявления через REDIS
	// статик для стилей
	h.Echo.Static("/css", "web/css")
	// статик для изображений(аватарки, фото объяв и прочее)
	h.Echo.Static("/images", "./images")
	// Инициализация шаблонизатора
	t := &Template{
		templates: template.Must(template.ParseGlob("web/html/*.html")),
	}
	h.Echo.Renderer = t
}

// Обработчик GET /home, главная страница
func (h *Handler) Home(c echo.Context) error {
	// Попытаться получить куку с именем "session_id"
	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			data, errors := h.UserSvc.GetAdsData()
			if errors != nil {
				h.Log.Error(err)
				fmt.Println("HUESOSaaaa") //
				return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
			}
			// Кука не найдена
			err := c.Render(http.StatusOK, "nonhome.html", data)
			if err != nil {
				h.Log.Error("Ошибка чтения html файла: ", err)
				return c.String(http.StatusInternalServerError, utils.ErrServer)
			}
			return nil

		} else {
			// Если ошибка НЕ ЯВЛЯЕТСЯ "отсутствие куки"
			return c.String(http.StatusInternalServerError, "Произошла ошибка с куки-файлом, попробуйте перезайти в свой аккаунт.")
		}
	}
	// Если юзер авторизован и с кукой все норм
	data, err := h.UserSvc.GetHomeData(cookie)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
	}

	err = c.Render(http.StatusOK, "home.html", data)
	if err != nil {
		h.Log.Error("Ошибка чтения html файла: ", err)
		return c.String(http.StatusInternalServerError, utils.ErrServer)
	}
	return nil
}

// Обработчик GET /auth
func (h *Handler) Authorization(c echo.Context) error {
	return c.Render(http.StatusOK, "Auth.html", nil)
}

// Обработчик POST /auth/login
func (h *Handler) Login(c echo.Context) error {
	login := c.FormValue("username")
	password := c.FormValue("password")

	jwtToken, err := h.UserSvc.AuthUser(login, password)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusUnauthorized, "Введены неверные данные, перепроверьте логин и пароль.")
	}
	// Отправка токена через куки
	h.UserSvc.SendToken(c.Response().Writer, jwtToken)

	//Наверное перенаправить человечка на домашнюю страницу, пока что не трогал
	err = c.Redirect(http.StatusSeeOther, "/home")
	if err != nil {
		h.Log.Error("Ошибка чтения html файла: ", err)
		return c.String(http.StatusInternalServerError, utils.ErrServer)
	}
	return nil
}

// Обработчик POST /auth/register
func (h *Handler) Register(c echo.Context) error {
	login := c.FormValue("username")
	password := c.FormValue("password")

	jwtToken, err := h.UserSvc.RegisterUser(login, password)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusUnauthorized, "Введены неверные данные.")
	}
	// Отправка токена через куки
	h.UserSvc.SendToken(c.Response().Writer, jwtToken)

	//Наверное перенаправить человечка на домашнюю страницу, пока что не трогал
	err = c.Redirect(http.StatusSeeOther, "/home")
	if err != nil {
		h.Log.Error("Ошибка чтения html файла: ", err)
		return c.String(http.StatusInternalServerError, utils.ErrServer)
	}

	return nil
}

// Обработчик GET /home/ads-form создание объявления
func (h *Handler) AdsForm(c echo.Context) error {
	return c.Render(http.StatusOK, "userAdsPost.html", nil)
}

// Обработчик GET /ads/:id ....
func (h *Handler) OneAds(c echo.Context) error {
	op := "handler.go - OneAds: "
	// заранее задаем параметр cookieChecker на то, что кука существует и юзер авторизован
	cookieChecker := true
	// получаем идентификатор объявления из URL ads/:id
	adsID := c.Param("id")
	if adsID == "" {
		h.Log.Error(op, "переход на URL: `/ads/` без указания идентификатора id невозможен")
		return c.String(http.StatusInternalServerError, "handler.go - OneAds: переход на URL: `/ads/` без указания идентификатора id невозможен")
	}

	// получаем куку
	cookie, err := c.Cookie("jwt")
	if err != nil {
		// если куки нет, значит пользователь не авторизован и мы задаем cookieChecker этот ответ
		if err == http.ErrNoCookie {
			cookieChecker = false
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла ошибка связанная с куки.")
		}
	}
	// получаем все данные для загрузки страницы(данные пользователя, объявления, название html файла)
	data, nameHTML, err := h.UserSvc.OneAds(cookie, cookieChecker, adsID)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "На сервере произошла ошибка с загрузкой страницы, попробуйте позже.")
	}
	//загружаем страницу и обрабатываем ошибку если появится
	err = c.Render(http.StatusOK, nameHTML, data)
	if err != nil {
		h.Log.Error("Ошибка чтения html файла: ", nameHTML, "\n", err)
		return c.String(http.StatusInternalServerError, utils.ErrServer)
	}

	return nil
}

// Обработчик GET /ads/form_update/:id форма для изменения объявления
func (h *Handler) AdsUpdateForm(c echo.Context) error {
	op := "handler.go - AdsUpdateForm: "
	// получаем данные и объявляем булевое значение
	adsID := c.Param("id")
	cookieChecker := true

	// получаем куку и делаем проверки
	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Log.Error(op, "отсутствует кука, направляем юзера на /home")
			return c.String(http.StatusInternalServerError, "Нам неизвестно как вы сюда попали, сначала авторизуйтесь.")
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}

	// вызываем уже существующий метод(из-за этого и пришлось создать cookieChecker)
	// получаем данные пользователя, объявления и обрабатываем ошибку
	data, ownerChecker, err := h.UserSvc.OneAds(cookie, cookieChecker, adsID)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "Произошла ошибка при получении данных, попробуйте позже.")
	}

	// с помощью полученной строки смотрим: является ли юзер владельцем объявления
	switch ownerChecker {
	//является владельцем
	case "ads-owner.html":
		err = c.Render(http.StatusOK, "adsUpdate.html", data)
		if err != nil {
			h.Log.Error("Ошибка чтения html файла: ", err)
			return c.String(http.StatusInternalServerError, utils.ErrServer)
		}
		// не является владельцем
	default:
		h.Log.Error("пользователь пытается изменить объявление, когда не является владельцем объявления")
		return c.String(http.StatusInternalServerError, "Вы не являетесь владельцем объявления.")
	}

	return nil
}

// Обработчик POST /ads/post
func (h *Handler) PostAds(c echo.Context) error {
	op := "handler.go - PostAds: "

	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Log.Error(op, "отсутствует кука, направляем юзера на /home")
			return c.String(http.StatusInternalServerError, "Нам неизвестно как вы сюда попали, сначала авторизуйтесь.")
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}

	adsName := c.FormValue("ads_name")
	adsDescription := c.FormValue("ads_description")
	price := c.FormValue("price")

	var priceInt float64

	if price == "" || price == "0" {
		priceInt = 0.00
	} else {
		priceInt, err = strconv.ParseFloat(price, 64)
		if err != nil {
			h.Log.Errorf(op, "ошибка конвертации в float64\nERROR: %v", err)
			return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
		}
	}

	images := []*multipart.FileHeader{}

	for i := 0; i <= 3; i++ {
		file, err := c.FormFile(fmt.Sprintf("image%d", i))
		if err != nil {
			// Skip if file is not uploaded
			continue
		}
		images = append(images, file)
	}

	id, err := h.UserSvc.SaveUserAds(cookie, images, adsName, adsDescription, priceInt)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
	}

	switch id {
	case "":
		c.Redirect(http.StatusSeeOther, "/home")
	default:
		urlName := "/ads/" + id
		c.Redirect(http.StatusSeeOther, urlName)
	}

	return nil
}

// Обработчик POST /ads/like
func (h *Handler) PostLike(c echo.Context) error {
	op := "handler.go - PostLike: "

	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Log.Error(op, "отсутствует кука, направляем юзера на /home")
			return c.String(http.StatusInternalServerError, "Нам неизвестно как вы сюда попали, сначала авторизуйтесь.")
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}

	adsID := c.Param("id")

	err = h.UserSvc.SaveLikedAd(adsID, cookie)
	if err != nil {
		h.Log.Error(op, "сбой при добавлении объявления в избранное: ", err)
	}

	urlPath := "/ads/" + adsID
	err = c.Redirect(http.StatusSeeOther, urlPath)
	if err != nil {
		h.Log.Error(op, "сбой при перенаправлении пользователя: ", err)
	}
	return nil
}

// Обработчик POST /ads/update
func (h *Handler) UpdateUserAds(c echo.Context) error {
	op := "handler.go - UpdateUserAds: "

	var oldImages []string
	// получение куки и проверка есть ли она вообще
	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Log.Error(op, "отсутствует кука, направляем юзера на /home")
			return c.String(http.StatusInternalServerError, "Нам неизвестно как вы сюда попали, сначала авторизуйтесь.")
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}

	// получение данных из формы (user_id и adsID пользователь не может менять, а даже если как-то поменяет у нас есть проверка)
	userID := c.FormValue("userid")
	adsID := c.FormValue("id")
	adsName := c.FormValue("ads_name")
	adsDescription := c.FormValue("ads_description")
	price := c.FormValue("price")
	oldImage1 := c.FormValue("old_image_1") // для удаления старой фотки
	oldImage2 := c.FormValue("old_image_2") // для удаления старой фотки
	oldImage3 := c.FormValue("old_image_3") // для удаления старой фотки

	oldImages = append(oldImages, oldImage1, oldImage2, oldImage3)
	// конвертим строку цены в значение float64
	var priceInt float64
	if price == "" || price == "0" {
		priceInt = 0.00
	} else {
		priceInt, err = strconv.ParseFloat(price, 64)
		if err != nil {
			h.Log.Errorf(op, "ошибка конвертации в float64\nERROR: %v", err)
			return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
		}
	}

	// создаем слайс фоток для удобной работы в методе
	images := []*multipart.FileHeader{}
	formImages, err := c.MultipartForm()
	if err != nil {
		h.Log.Errorf(op, "ошибка получения обновленного фото\nERROR: %v", err)
	} else {
		images = formImages.File["images"]
	}

	// метод, который валидирует полученные данные, делает проверку на пользователя и апдейтит наше объявление
	err = h.UserSvc.UpdateUserAds(cookie, userID, adsID, adsName, adsDescription, priceInt, images, oldImages)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
	}

	// в случае успешного обновления, перенаправляем на это же объявление
	UrlPath := "/ads/" + adsID
	err = c.Redirect(http.StatusSeeOther, UrlPath)
	if err != nil {
		h.Log.Error(op, "Ошибка перенаправления пользователя: ", err)
		return err
	}

	return nil
}

// Обработчик POST /ads/delete
func (h *Handler) DelAds(c echo.Context) error {
	op := "handler.go - DelAds: "

	adsID := c.Param("id")
	if adsID == "" {
		h.Log.Error(op, "Отсутствует параметр adsID")
		return c.String(http.StatusInternalServerError, "Мы не можем понять какое объявление вы хотите удалить.")
	}

	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Log.Error(op, "отсутствует кука, направляем юзера на /home")
			return c.String(http.StatusInternalServerError, "Нам неизвестно как вы сюда попали, сначала авторизуйтесь.")
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}

	owner, err := h.UserSvc.DeleteAds(cookie, adsID)
	if err != nil {
		switch owner {
		case true:
			h.Log.Error(err)
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		case false:
			h.Log.Error(err)
			return c.String(http.StatusInternalServerError, "У вас нет доступа к удалению чужих объявлений.")
		}
	}

	err = c.Redirect(http.StatusSeeOther, "/home")
	if err != nil {
		h.Log.Error(op, "ошибка перенаправления пользователя на домашнюю страницу:\n", err)
		return err
	}

	return nil
}

// Обработчик GET /Profile?id=...
func (h *Handler) Profile(c echo.Context) error {
	op := "handler.go - Profile: "

	// проверка на авторизацию
	cook := true
	// получение идентификатора профиля
	login := c.Param("id")
	// проверка на владельца профиля
	owner := false
	var guestData *models.UserData
	// получение куки
	cookie, err := c.Cookie("jwt")
	if err != nil {
		// если куки нет, выставляем значение cook на false
		if err == http.ErrNoCookie {
			cook = false
			// все равно выставляем значение cook на false
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			cook = false
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}

	// проверяем куку, если она есть то проверяем на владельца профиля
	if cook {
		owner, err = h.UserSvc.VerifyUserOwner(cookie, login)
		if err != nil {
			h.Log.Error(err)
			owner = false
		}
	}

	if cook && !owner {
		guestData, err = h.UserSvc.GetGuestData(cookie)
		if err != nil {
			h.Log.Error(err)
			cook = false
		}
	}

	// получаем данные о профиле из БД
	data, err := h.UserSvc.GetProfileData(login)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
	}

	// в этом случае направляем на страницу для владельца профиля
	if owner && cook {
		err := c.Render(http.StatusOK, "profile.html", data)
		if err != nil {
			h.Log.Error("ошибка загрузки html файла:\n", err)
			return err
		}
		return nil
		// в этом случае направляем на страницу для гостей
	} else if !owner && cook {
		// прописываем в структуру данные гостя
		data = &models.ProfileData{
			GuestData: guestData,
			UserData:  data.UserData,
			UserAds:   data.UserAds,
			Likes:     data.Likes,
		}

		err := c.Render(http.StatusOK, "NoOwnerProfile.html", data)
		if err != nil {
			h.Log.Error("ошибка загрузки html файла:\n", err)
			return err
		}
		return nil
		// в этом случае направляем на страницу для неавторизованных гостей
	} else {
		err := c.Render(http.StatusOK, "NoAuthProfile.html", data)
		if err != nil {
			h.Log.Error("ошибка загрузки html файла:\n", err)
			return err
		}
		return nil
	}
}

// Обработчик GET /profile/form
func (h *Handler) ProfileForm(c echo.Context) error {
	op := "handler.go - ProfileFrom: "

	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Log.Error(op, "отсутствует кука, направляем юзера на /home")
			return c.String(http.StatusInternalServerError, "Нам неизвестно как вы сюда попали, сначала авторизуйтесь.")
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}
	data, err := h.UserSvc.GetUserData(cookie)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "Произошла ошибка при получении данных, попробуйте позже.")
	}
	err = c.Render(http.StatusOK, "profileChange.html", data)
	if err != nil {
		h.Log.Error("Ошибка чтения html файла: ", err)
		return c.String(http.StatusInternalServerError, utils.ErrServer)
	}
	return nil
}

// Обработчик POST /profile/update
func (h *Handler) UpdateProfile(c echo.Context) error {
	op := "handler.go - UpdateProfile: "

	missingImage := false
	cookie, err := c.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Log.Error(op, "отсутствует кука, направляем юзера на /home")
			return c.String(http.StatusInternalServerError, "Нам неизвестно как вы сюда попали, сначала авторизуйтесь.")
		} else {
			h.Log.Error(op, "сбой при получении куки.")
			return c.String(http.StatusInternalServerError, "На сервере произошла непредвиденная ошибка.")
		}
	}
	name := c.FormValue("name")
	firstname := c.FormValue("first_name")
	phonenumber := c.FormValue("phone_number")
	oldImage := c.FormValue("old_image_profile") // для удаления старой фотки
	images := []*multipart.FileHeader{}
	file, err := c.FormFile("image_profile")
	if err != nil {
		// Пользователь не загрузил файл, просто пропускаем обработку
		missingImage = true
		// Обрабатываем другие ошибки
		h.Log.Error(op, "сбой при получении аватарки либо ее не загружали.")
	}

	images = append(images, file)

	login, err := h.UserSvc.UpdateUserData(cookie, name, firstname, phonenumber, images, missingImage, oldImage)
	if err != nil {
		h.Log.Error(err)
		return c.String(http.StatusInternalServerError, "Ошибка загрузки страницы.")
	}

	UrlPath := "/profile/" + login
	err = c.Redirect(http.StatusSeeOther, UrlPath)
	if err != nil {
		h.Log.Error("Ошибка чтения html файла: ", err)
		return c.String(http.StatusInternalServerError, utils.ErrServer)
	}

	return nil
}

// Обработчик POST /profile/quit
func (h *Handler) QuitProfile(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		h.Log.Error(err)
		err = c.Redirect(http.StatusSeeOther, "/home")
		if err != nil {
			h.Log.Error(err)
		}
	}

	jwtToken, err := h.UserSvc.QuitFromAccount(cookie)
	if err != nil {
		h.Log.Error(err)
	}

	// Отправка токена через куки
	h.UserSvc.SendExpiredToken(c.Response().Writer, jwtToken)

	//Наверное перенаправить человечка на домашнюю страницу, пока что не трогал
	err = c.Redirect(http.StatusSeeOther, "/home")
	if err != nil {
		h.Log.Error("Ошибка чтения html файла: ", err)
		return c.String(http.StatusInternalServerError, utils.ErrServer)
	}

	return nil
}

// JWTMiddleware проверяет подлинность JWT токена и извлекает логин
func (h *Handler) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			h.Log.Error(err)
			return c.Redirect(http.StatusInternalServerError, "/home")
		}

		if err = h.UserSvc.CheckToken(cookie.Value); err != nil {
			h.Log.Error(err)
			return c.Redirect(http.StatusInternalServerError, "/home")
		}

		// Передаем управление следующему обработчику
		return next(c)
	}
}

// JWTMiddleware проверяет подлинность JWT токена и извлекает логин
//func (h *Handler) ViewsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := context.Background()
//		adsID := c.Param(":id")
//
//		if err := h.UserSvc.IncrViews(ctx, adsID); err != nil {
//			h.Log.Error(err)
//		}
//		// Передаем управление следующему обработчику
//		return next(c)
//	}
//}
