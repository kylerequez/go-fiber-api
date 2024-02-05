package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/kylerequez/go-fiber-crud/api/todo"
	"github.com/kylerequez/go-fiber-crud/common"
)

func run() error {
	err := common.LoadEnvVariables()
	if err != nil {
		return err
	}

	err = common.ConnectDB()
	if err != nil {
		return err
	}

	defer common.CloseDB()

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	todo.InitRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("pingpong")
	})

	var port string = os.Getenv("PORT")
	app.Listen(":" + port)
	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
