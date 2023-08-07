package routes

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
	"xixo.cf/profileapi/database"
	"xixo.cf/profileapi/types"
)

var tp = []string{"malephoto", "femalephoto", "malegif", "banner", "faceless", "anime", "femalegif", "cars", "nike", "nsfw", "aesthetic", "cartoon", "jewellry", "shoes", "guns", "drill", "money", "smoking", "animals", "soft", "hellokitty", "besties", "match"}

func Random(c *fiber.Ctx) error {
	t := c.Params("type", "")
	if t == "" || !slices.Contains(tp, t) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid type.",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	max, err := database.GetMongo().Database("ProfileAPI").Collection("images").CountDocuments(ctx, bson.M{
		"type": t,
	})

	if err != nil {
		fmt.Println("Erorr while counting items from db", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}

	var image types.Image

	err = database.GetMongo().Database("ProfileAPI").Collection("images").FindOne(ctx, bson.M{
		"type": t,
	}, options.FindOne().SetSkip(*randomInt64InRange(max))).Decode(&image)

	if err != nil {
		fmt.Printf("Erorr while getting random item from db(%d): %s", max, err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}

	return c.JSON(fiber.Map{
		"url": os.Getenv("PUBLIC_URL") + "/image/" + image.Id,
	})

}

func randomInt64InRange(x int64) *int64 {
	// Create a new random number generator with a custom source seeded with the current time
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	// Generate a random integer between 0 and x (inclusive)
	randomInt := randomGen.Int63n(x + 1)

	return &randomInt
}
