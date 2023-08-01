package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"xixo.cf/profileapi/database"
	"xixo.cf/profileapi/types"
)

// main route

func Image(c *fiber.Ctx) error {
	id := c.Params("id", "")

	if id == "" {
		return c.SendStatus(400)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var image types.Image

	err := database.GetMongo().Database("ProfileAPI").Collection("images").FindOne(ctx, bson.M{
		"id": id,
	}).Decode(&image)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.SendStatus(505)
	}

	res, err := http.Get(image.Src)

	if err != nil {
		return c.SendStatus(500)
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=604800")
	c.Set(fiber.HeaderLastModified, time.Now().UTC().Format(time.RFC1123))

	c.Set(fiber.HeaderContentType, res.Header.Get("Content-Type"))

	return c.SendStream(res.Body, int(res.ContentLength))
}
