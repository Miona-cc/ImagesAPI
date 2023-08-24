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

var tp = []string{"malephoto", "femalephoto", "malegif", "banner", "faceless", "femalegif", "cars", "nike", "nsfw", "aesthetic", "cartoon", "jewellry", "shoes", "guns", "drill", "money", "smoking", "animals", "soft", "hellokitty", "besties", "body", "food", "random", "anime"}

func Random(c *fiber.Ctx) error {
	t := c.Params("type", "")
	if t == "" || !slices.Contains(tp, t) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid type.",
		})
	}

	max, err := database.GetMongo().Database("ProfileAPI").Collection("images").CountDocuments(context.Background(), bson.M{
		"type": t,
	})

	if err != nil {
		fmt.Println("Erorr while counting items from db", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}

	var image types.Image

	err = database.GetMongo().Database("ProfileAPI").Collection("images").FindOne(context.Background(), bson.M{
		"type": t,
	}, options.FindOne().SetSkip(*RandomInt64InRange(max))).Decode(&image)

	if err != nil {
		fmt.Printf("Erorr while getting random item from db(%d): %s", max, err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong.",
		})
	}

	//If the image doesn't exists, call the function again
	if CheckImage(image.Src) {
		Random(c)
		return nil
	}

	return c.JSON(fiber.Map{
		"url": os.Getenv("PUBLIC_URL") + "/image/" + image.Id,
	})

}

func RandomInt64InRange(x int64) *int64 {
	if x <= 0 {
		var i int64 = 1
		return &i
	}
	// Create a new random number generator with a custom source seeded with the current time
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	// Generate a random integer between 0 and x (inclusive)
	randomInt := randomGen.Int63n(x)

	return &randomInt
}
