package main

import (
	"clean1/config"
	td "clean1/features/todo/data"
	th "clean1/features/todo/handler"
	ts "clean1/features/todo/services"
	"clean1/features/user/data"
	"clean1/features/user/handler"
	"clean1/features/user/services"
	"clean1/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()            // inisiasi echo
	cfg := config.InitConfig() // baca seluruh system variable
	db := config.InitSQL(cfg)  // konek DB

	userData := data.New(db)
	userService := services.NewService(userData)
	userHandler := handler.NewUserHandler(userService)

	todoData := td.New(db)
	todoService := ts.NewTodoService(todoData)
	todoHandler := th.NewHandler(todoService)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) // ini aja cukup
	routes.InitRoute(e, userHandler, todoHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
