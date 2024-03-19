package handler

import (
	"clean1/features/todo"
	"clean1/helper"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type controller struct {
	s todo.TodoService
}

func NewHandler(service todo.TodoService) todo.TodoController {
	return &controller{
		s: service,
	}
}

func (ct *controller) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input ToDoRequest
		err := c.Bind(&input)
		if err != nil {
			log.Println("error bind data:", err.Error())
			if strings.Contains(err.Error(), "unsupport") {
				return c.JSON(http.StatusUnsupportedMediaType,
					helper.ResponseFormat(http.StatusUnsupportedMediaType, helper.UserInputFormatError, nil))
			}
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest,
				helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		var inputProcess todo.Todo
		inputProcess.Kegiatan = input.Kegiatan
		result, err := ct.s.AddTodo(token, inputProcess)
		if err != nil {
			log.Println("error insert db:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "berhasil menambahkan kegiatan", result))
	}
}

func (ct *controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		var updatedTodo todo.Todo
		if err := c.Bind(&updatedTodo); err != nil {
			log.Println("error binding data:", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		// Mendapatkan ID todo yang akan di-update dari URL atau dari body request
		todoID := c.Param("id")

		// Memperbarui todo
		result, err := ct.s.UpdateTodo(token, todoID, updatedTodo)
		if err != nil {
			log.Println("error updating todo:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Todo updated successfully", result))
	}
}

func (ct *controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, helper.UserInputError, nil))
		}

		// Mendapatkan ID todo yang akan dihapus dari URL atau dari body request
		todoID := c.Param("id")

		// Menghapus todo
		if err := ct.s.DeleteTodo(token, todoID); err != nil {
			log.Println("error deleting todo:", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, helper.ServerGeneralError, nil))
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Todo deleted successfully", nil))
	}
}
