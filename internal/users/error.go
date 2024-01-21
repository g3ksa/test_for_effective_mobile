package users

import "errors"

var ErrInvalidAccessToken = errors.New("некорректный токен авторизаци")
var ErrUserDoesNotExist = errors.New("такого польхователя не существует")
var ErrUserAlreadyExists = errors.New("пользователь с такими данными уже существует")
var ErrIncorrectPassword = errors.New("неверный пароль")
