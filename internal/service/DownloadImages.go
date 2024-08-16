package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (s *UserService) DownloadFiles(files []*multipart.FileHeader, path string) ([]string, error) {

	completeImages := []string{}
	for _, c := range files {
		uniqueFileName := uuid.New().String()
		uniqueFileName = uniqueFileName[:20]
		// Открытие загруженного файла
		src, err := c.Open()
		if err != nil {
			return nil, fmt.Errorf("DownloadImages.go - DownloadFiles: ошибка при открытии фото")
		}
		defer src.Close()

		// Создание файла назначения на сервере
		fileExtension := filepath.Ext(c.Filename)
		fullName := fmt.Sprintf("%s%s", uniqueFileName, fileExtension)
		dstPath := fmt.Sprintf("./images/%s/%s", path, fullName)

		dst, err := os.Create(dstPath)
		if err != nil {
			return nil, fmt.Errorf("DownloadImages.go - DownloadFiles: ошибка при создании файла")
		}
		defer dst.Close()

		// Копирование загруженного файла в файл назначения
		if _, err = io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("DownloadImages.go - DownloadFiles: ошибка при записи фото в наш файл")
		}
		// Обрезаем первый символ
		dstPath = dstPath[1:]

		completeImages = append(completeImages, dstPath)
	}

	return completeImages, nil
}

func (s *UserService) DeleteImages(pathImages []string) error {
	for _, c := range pathImages {
		path := "." + c
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("ошибка удаления файла: %v", err)
		}
	}

	return nil
}
