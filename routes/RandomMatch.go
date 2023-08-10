package routes

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
	"xixo.cf/profileapi/database"
	"xixo.cf/profileapi/types"
)

var choices = []string{"matchphoto", "matchgif", "matchbanner"}

func RandomMatch(c *fiber.Ctx) error {
	t := c.Params("type", "")
	if t == "" || !slices.Contains(choices, t) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid type.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	max, err := database.GetMongo().Database("ProfileAPI").Collection("matches").CountDocuments(ctx, bson.M{
		"type": t,
	})

	if err != nil {
		fmt.Println("Erorr while counting items from db", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}

	var match types.Match

	err = database.GetMongo().Database("ProfileAPI").Collection("matches").FindOne(ctx, bson.M{
		"type": t,
	}, options.FindOne().SetSkip(*RandomInt64InRange(max))).Decode(&match)

	if err != nil {
		fmt.Printf("Erorr while getting random item from db(%d): %s", max, err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}

	return c.JSON([]string{os.Getenv("PUBLIC_URL") + "/image/" + match.Srcs[0], os.Getenv("PUBLIC_URL") + "/image/" + match.Srcs[1]})
}
