package handlers

import "errors"

// Register
var ErrSameNicknameExists = errors.New("пользователь с таким ником уже существует")
var ErrSameEmailExists = errors.New("пользователь с таким email уже существует")

// Internal
var ErrInternalServer = errors.New("что-то пошло не так")

// Login
var ErrNoUserWithEmail = errors.New("пользователя с таким email не существует")
var ErrInvalidCredentials = errors.New("неверный пароль")
var ErrInvalidToken = errors.New("неверный токен подтверждения")

// Post
var ErrInvalidInputData = errors.New("неправильные входные данные")
