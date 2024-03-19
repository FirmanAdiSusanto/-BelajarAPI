package services

import (
	"clean1/features/user"
	"clean1/helper"
	"clean1/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	model user.UserModel
	pm    helper.PasswordManager
	v     *validator.Validate
}

func NewService(m user.UserModel) user.UserService {
	return &service{
		model: m,
		pm:    helper.NewPasswordManager(),
		v:     validator.New(),
	}
}

func (s *service) Register(newData user.User) error {
	var registerValidate user.Register
	registerValidate.Hp = newData.Hp
	registerValidate.Nama = newData.Nama
	registerValidate.Password = newData.Password
	err := s.v.Struct(&registerValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	newPassword, err := s.pm.HashPassword(newData.Password)
	if err != nil {
		return errors.New(helper.ServiceGeneralError)
	}
	newData.Password = newPassword

	err = s.model.InsertUser(newData)
	if err != nil {
		return errors.New(helper.ServerGeneralError)
	}
	return nil
}
func (s *service) Login(loginData user.User) (user.User, string, error) {
	var loginValidate user.Login
	loginValidate.Hp = loginData.Hp
	loginValidate.Password = loginData.Password
	err := s.v.Struct(&loginValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return user.User{}, "", err
	}

	dbData, err := s.model.Login(loginValidate.Hp)
	if err != nil {
		return user.User{}, "", err
	}

	err = s.pm.ComparePassword(loginValidate.Password, dbData.Password)
	if err != nil {
		return user.User{}, "", errors.New(helper.UserCredentialError)
	}

	token, err := middlewares.GenerateJWT(dbData.Hp)
	if err != nil {
		return user.User{}, "", errors.New(helper.ServiceGeneralError)
	}

	return dbData, token, nil
}

func (s *service) Profile(token *jwt.Token) (user.User, error) {
	decodeHp := middlewares.DecodeToken(token)
	result, err := s.model.GetUserByHP(decodeHp)
	if err != nil {
		return user.User{}, err
	}

	return result, nil
}
