package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karanjar/cargobackend_fibre_framework.git/config"
	"github.com/karanjar/cargobackend_fibre_framework.git/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	options2 "go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Mu sync.Mutex

func getId(ctx context.Context) (int, error) {
	col := config.Client.Database("car_inventory").Collection("cars")

	var counter struct {
		ID  string `bson:"_id"`
		Seq int    `bson:"seq"`
	}
	filter := bson.M{"_id": "car_id"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	options := options2.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options2.After)

	if err := col.FindOneAndUpdate(ctx, filter, update, options).Decode(&counter); err != nil {
		return 0, fmt.Errorf("errer finding car id from context %v", err)
	}
	return counter.Seq, nil
}

// Createcar  godoc
// @Summary Create a new car
// @Description Add a new car to the database
// @Tags cars
// @Accept  json
// @Produce  json
// @Param car body models.Car true "Car data"
// @Success 200 {object} models.Car
// @Failure 400 {object}  models.Error
// @Router /cars [post]
func Createcar(c *fiber.Ctx) error {
	Mu.Lock()
	defer Mu.Unlock()

	car := &models.Car{}

	if err := c.BodyParser(car); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.Error{
			Message: "incorrect input body",
			Details: err.Error(),
		})
	}

	//if err := car.Insert(); err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"error": "incorrect input body",
	//	})
	//}

	//car insert

	id, err := getId(c.Context())
	if err != nil {
		fmt.Printf("error getting id from context: %v\n", err)
	}
	car.Id = id

	coll := config.Client.Database("car_inventory").Collection("cars")
	reault, err := coll.InsertOne(c.Context(), car)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.Error{
			Message: "unable to add a new car",
			Details: err.Error(),
		})
	}

	fmt.Println(reault.InsertedID)

	fmt.Println("Car created with the id:", car.Id)
	return c.Status(fiber.StatusCreated).JSON(car)
}

// Getcar godoc
// @Summary Get  a new car
// @Description Get a car from the inventory
// @Tags cars
// @Accept  json
// @Produce  json
// @Param id path string true "Car id"
// @Success 200 {object} models.Car
// @Failure 400 {object}  models.Error
// @Failure 404 {object} models.Error
// @Router /cars/{id} [get]
func Getcar(c *fiber.Ctx) error {
	Mu.Lock()
	defer Mu.Unlock()

	//check  if the car is already presnt in the cache
	//if not then only goto the postgres

	car := &models.Car{}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.Error{
			Message: "invalid car id",
			Details: err.Error(),
		})
	}
	key := strconv.FormatInt(int64(id), 10)
	val, err := config.Cache.Get(c.Context(), key).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(val), car); err != nil {
			fmt.Printf("unable to unmarshal cached car: %v\n", err)
		} else {
			fmt.Printf("cache hit for car ID:%v", id)
			return c.Status(fiber.StatusOK).JSON(car)
		}
	} else if !errors.Is(err, redis.Nil) {
		fmt.Printf("redis error:%v\n", id)
	} else {
		fmt.Printf("cache miss for id:%v\n", id)
	}

	//car.Id = id

	if err := car.Get(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "car with the given id is not found",
			"id":    car.Id,
		})
	}

	//fmt.Println("Car found with the id:", id)
	b, _ := json.Marshal(car)
	err = config.Cache.Set(c.Context(), key, b, 60*time.Minute).Err()
	if err != nil {
		fmt.Printf("unable to add key to the redis:  %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(car)

}

func Deletecar(c *fiber.Ctx) error {
	Mu.Lock()
	defer Mu.Unlock()
	car := &models.Car{}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid delete car id",
		})
	}

	//car.Id = id
	if err := car.Delete(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "car with the given id does not found",
		})
	}

	fmt.Println("Car deleted with the id:", id)
	return c.SendStatus(fiber.StatusNoContent)
}
func Updatecar(c *fiber.Ctx) error {
	Mu.Lock()
	defer Mu.Unlock()
	car := &models.Car{}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid update car id",
		})
	}

	if err := c.BodyParser(car); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "incorrect request body",
		})
	}

	//car.Id = id

	if err := car.Update(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "car with the given id is not found",
		})
	}
	fmt.Println("Car Updated with the id:", id)
	return c.Status(fiber.StatusCreated).JSON(car)

}
