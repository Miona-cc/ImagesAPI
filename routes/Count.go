package routes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"xixo.cf/profileapi/database"
)

type result struct {
	Type  string `json:"type" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}

func Count(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := database.GetMongo().Database("ProfileAPI").Collection("images").Aggregate(ctx, bson.A{
		bson.M{
			"$group": bson.M{
				"_id":  "$type",
				"count": bson.M{"$sum": 1},
			},
		},
	})

	if err != nil {
		fmt.Println("Error while executing pipline: ", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong on our side.",
		})
	}

	defer cur.Close(ctx)

	var results []result

	for cur.Next(ctx) {
		var res result

		if err := cur.Decode(&res); err != nil {
			fmt.Println("Error while decoding response inside for: ", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Something went wrong on our side.",
			})
		}
		if strings.HasPrefix(res.Type, "match") {
			res.Count/=2
		}
		results = append(results, res)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error while cursor iterates: ", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Something went wrong on our side.",
		})
	}

	return c.JSON(results)
}
