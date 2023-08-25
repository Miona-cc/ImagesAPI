package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"xixo.cf/profileapi/database"
	"xixo.cf/profileapi/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Couldn't load .env file")
	}

	//Connect databases
	database.GetMongo()
	database.GetRedis()

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/image") {
			return c.Next()
		}

		split := strings.Split(c.Get("Authorization"), " ")
		if len(split) < 2 || split[1] == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := database.GetRedis().Get(ctx, split[1]).Err(); err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		return c.Next()
	})
	if os.Getenv("ENV") == "production" {
		app.Use(logger.New())
	}

	app.Use(limiter.New(
		limiter.Config{
			Max:        100000,
			Expiration: time.Hour,
			Next: func(c *fiber.Ctx) bool {
				return strings.HasPrefix(c.Path(), "/image")
			},
			KeyGenerator: func(c *fiber.Ctx) string {
				return strings.Split(c.Get("Authorization"), " ")[1]
			},
		},
	))
	app.Get("/count", routes.Count)
	app.Get("/image/:id", routes.Image)
	app.Get("/random/:type", routes.Random)
	app.Get("/match/:type", routes.RandomMatch)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	env := os.Getenv("ENV")
	if env != "production" {
		app.Listen("127.0.0.1:" + port)
	} else {
		app.Listen(":" + port)
	}
}
