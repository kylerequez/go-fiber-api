package todo

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kylerequez/go-fiber-crud/common"
	"github.com/kylerequez/go-fiber-crud/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionName string = "todos"

func InitRoutes(app *fiber.App) {
	todosRoute := app.Group("/api/v1/todos")
	todosRoute.Get("/", GetAllTodos)
	todosRoute.Get("/:id", GetTodo)
	todosRoute.Post("/", CreateTodo)
	todosRoute.Put("/:id", PutUpdateTodo)
	todosRoute.Patch("/:id", PatchUpdateTodo)
	todosRoute.Delete("/:id", DeleteTodo)
}

func GetAllTodos(c *fiber.Ctx) error {
	todos := make([]types.Todo, 0)

	coll := common.GetCollection(collectionName)

	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for cursor.Next(c.Context()) {
		todo := types.Todo{}
		err := cursor.Decode(&todo)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		todos = append(todos, todo)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"todos":   todos,
	})
}

func GetTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Params must not be empty",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	coll := common.GetCollection(collectionName)

	var todo types.Todo
	filters := bson.M{"_id": objectId}
	err = coll.FindOne(context.TODO(), filters).Decode(&todo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"todo":    todo,
	})
}

func CreateTodo(c *fiber.Ctx) error {
	type RequestBody struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	req := new(RequestBody)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	coll := common.GetCollection(collectionName)
	newTodo := types.Todo{Title: req.Title, Body: req.Body}

	result, resErr := coll.InsertOne(context.TODO(), newTodo)
	if resErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": resErr.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"todo":    result,
	})
}

func PutUpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Params must not be empty",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	type RequestBody struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	req := new(RequestBody)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	coll := common.GetCollection(collectionName)
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"title": req.Title, "body": req.Body}}
	opts := options.Update().SetUpsert(true)

	result, resErr := coll.UpdateOne(context.TODO(), filter, update, opts)
	if resErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": resErr.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"todo":    result,
	})
}

func PatchUpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Params must not be empty",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	type RequestBody struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	req := new(RequestBody)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	coll := common.GetCollection(collectionName)

	var todo types.Todo
	filters := bson.M{"_id": objectId}
	err = coll.FindOne(context.TODO(), filters).Decode(&todo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"title": req.Title, "body": req.Body}}
	opts := options.Update().SetUpsert(false)

	result, resErr := coll.UpdateOne(context.TODO(), filter, update, opts)
	if resErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": resErr.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"todo":    result,
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Params must not be empty",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	coll := common.GetCollection(collectionName)
	filter := bson.M{"_id": objectId}

	result, resErr := coll.DeleteOne(context.TODO(), filter)
	if resErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": resErr.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"todo":    result,
	})
}
