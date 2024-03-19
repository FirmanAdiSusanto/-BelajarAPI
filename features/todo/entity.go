package todo

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TodoController interface {
	Add() echo.HandlerFunc
	// Update() echo.HandlerFunc
	// Delete() echo.HandlerFunc
	// ShowMyTodo() echo.HandlerFunc
}

type TodoModel interface {
	InsertTodo(pemilik string, kegiatanBaru Todo) (Todo, error)
	UpdateTodo(pemilik string, todoID uint, data Todo) (Todo, error)
	DeleteTodo(pemilik string, todoID uint) error
	GetTodoByOwner(pemilik string) ([]Todo, error)
}

type TodoService interface {
	AddTodo(pemilik *jwt.Token, kegiatanBaru Todo) (Todo, error)
	UpdateTodo(pemilik *jwt.Token, todoID string, data Todo) (Todo, error)
	DeleteTodo(pemilik *jwt.Token, todoID string) error
}

type Todo struct {
	Kegiatan string
}
