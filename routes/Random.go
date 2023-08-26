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

var tp = []string{"malephoto", "femalephoto", "malegif", "banner", "faceless", "femalegif", "cars", "nike", "nsfw", "aesthetic", "cartoon", "jewellry", "shoes", "guns", "drill", "money", "smoking", "animals", "soft", "hellokitty", "besties", "body", "food", "random", "anime"}

func Random(c *fiber.Ctx) error {
	t := c.Params("type", "")
	if t == "" || !slices.Contains(tp, t) {
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

	var image types.Image

	cur, err := database.GetMongo().Database("ProfileAPI").Collection("images").Aggregate(context.Background(), pipeline)
	if err != nil {
		fmt.Printf("Error while aggregating: %s", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}
	defer cur.Close(context.Background())

	if cur.Next(context.Background()) {
		err := cur.Decode(&image)
		if err != nil {
			fmt.Printf("Error while decoding: %s", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Something went wrong.",
			})
		}

		// If the image doesn't exist, call the function again
		if !CheckImage(image.Src) {
			return Random(c)
		}

		return c.JSON(fiber.Map{
			"url": os.Getenv("PUBLIC_URL") + "/image/" + image.Id,
		})
	}

	return c.Status(404).JSON(fiber.Map{
		"error": "No image found.",
	})
}
