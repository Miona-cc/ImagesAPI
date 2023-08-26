package routes

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

	pipeline := []bson.M{
		{
			"$match": bson.M{"type": t},
		},
		{
			"$sample": bson.M{"size": 1},
		},
	}

	var match types.Match

	cur, err := database.GetMongo().Database("ProfileAPI").Collection("matches").Aggregate(context.Background(), pipeline)
	if err != nil {
		fmt.Printf("Error while aggregating: %s", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}
	defer cur.Close(context.Background())

	if cur.Next(context.Background()) {
		err := cur.Decode(&match)
		if err != nil {
			fmt.Printf("Error while decoding: %s", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Something went wrong.",
			})
		}

		return c.JSON([]string{os.Getenv("PUBLIC_URL") + "/image/" + match.Srcs[0], os.Getenv("PUBLIC_URL") + "/image/" + match.Srcs[1]})
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "No match found.",
	})
}
