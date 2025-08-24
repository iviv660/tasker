package main

import (
	"app/internal/config"
	"app/internal/database"
	_ "app/internal/docs"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/usecase"
	"log"
)

// @title           Task Manager API
// @version         1.0
// @description     CRUD задач с регистрацией, логином (JWT) и защищёнными эндпоинтами.
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description     Формат: "Bearer {token}"

// @contact.name    API Support
// @contact.email   volodya.mir05@mail.ru
// @license.name    MIT
func main() {
	DB, err := database.ConnectPostgres(config.C.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	UserDB := repository.NewUserRepo(DB)
	TaskDB := repository.NewTaskRepo(DB)

	UserUC := usecase.NewUserUseCase(UserDB)
	TaskUC := usecase.NewTaskUseCase(TaskDB)

	router, _ := handler.NewHandler(TaskUC, UserUC)
	if err = router.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
