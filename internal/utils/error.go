package utils

import "errors"

var (
	ErrNotFound     = errors.New("запись не найдена")
	ErrDeleteFail   = errors.New("ошибка удаления записи")
	ErrNoOwnerAds   = errors.New("пользователь не является владельцем объявления")
	ErrAlreadyExist = errors.New("запись уже существует")
)

const (
	ErrServer = "На сервере произошла ошибка, попробуйте позже."
)
