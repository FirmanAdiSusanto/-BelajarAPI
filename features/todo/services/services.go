package services

import (
	"clean1/features/todo"
	"clean1/helper"
	"clean1/middlewares"
	"errors"
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	m todo.TodoModel
	v *validator.Validate
}

func NewTodoService(model todo.TodoModel) todo.TodoService {
	return &service{
		m: model,
		v: validator.New(),
	}
}

func (s *service) AddTodo(pemilik *jwt.Token, kegiatanBaru todo.Todo) (todo.Todo, error) {
	hp := middlewares.DecodeToken(pemilik)
	if hp == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return todo.Todo{}, errors.New("data tidak valid")
	}

	err := s.v.Struct(&kegiatanBaru)
	if err != nil {
		log.Println("error validasi", err.Error())
		return todo.Todo{}, err
	}

	result, err := s.m.InsertTodo(hp, kegiatanBaru)
	if err != nil {
		return todo.Todo{}, errors.New(helper.ServerGeneralError)
	}

	return result, nil
}

func (s *service) UpdateTodo(pemilik *jwt.Token, todoID string, data todo.Todo) (todo.Todo, error) {
	// Mendapatkan informasi pengguna dari token
	hp := pemilik.Claims.(jwt.MapClaims)["username"].(string)

	// Mengonversi todoID dari string ke uint
	id, err := strconv.ParseUint(todoID, 10, 64)
	if err != nil {
		return todo.Todo{}, errors.New("ID todo tidak valid")
	}

	// Memperbarui todo
	result, err := s.m.UpdateTodo(hp, uint(id), data)
	if err != nil {
		return todo.Todo{}, errors.New("gagal memperbarui kegiatan")
	}

	return result, nil
}

func (s *service) DeleteTodo(pemilik *jwt.Token, todoID string) error {
	hp := pemilik.Claims.(jwt.MapClaims)["username"].(string)
	if hp == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return errors.New("data tidak valid")
	}

	// Konversi todoID dari string ke uint atau tipe lain yang sesuai
	id, err := strconv.ParseUint(todoID, 10, 64)
	if err != nil {
		log.Println("error parsing todo ID:", err.Error())
		return errors.New("ID todo tidak valid")
	}

	// Melakukan operasi penghapusan
	err = s.m.DeleteTodo(hp, uint(id))
	if err != nil {
		log.Println("error deleting todo:", err.Error())
		return errors.New("gagal menghapus kegiatan")
	}

	return nil
}
