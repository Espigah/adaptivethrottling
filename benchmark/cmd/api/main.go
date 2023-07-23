package main

import (
	"fmt"
	"math/rand"
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

type Config struct {
	FirstPointOfFailure int `json: "firstPointOfFailure"`
	Intermittency       int `json: "intermittency"`
	FailureStateEnabled bool
	FailureCount        int
	RequestCount        int
}

var configs map[string]*Config = make(map[string]*Config)

func main() {
	app := fiber.New()
	rand.Seed(time.Now().UnixNano())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/:testName", func(c *fiber.Ctx) error {
		testName := fmt.Sprintf("%s", c.Params("testName"))

		err := createConfigIfNotExist(testName, c)

		if err != nil {
			fmt.Println("error = ", err)
			return c.SendStatus(500)
		}
		config := configs[testName]
		config.RequestCount++

		if config.RequestCount < config.FirstPointOfFailure {
			return c.SendStatus(200)
		}

		config.FailureCount++

		if config.FailureCount%config.Intermittency == 0 {
			config.FailureStateEnabled = !config.FailureStateEnabled
		}

		if config.FailureStateEnabled {
			n := 1 + rand.Intn(10-1+1)
			if n >= 8 {
				return c.SendStatus(200)
			}

			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	})

	app.Listen(":3003")
}

func createConfigIfNotExist(testName string, c *fiber.Ctx) error {
	if _, ok := configs[testName]; !ok {
		config := new(Config)

		if err := c.BodyParser(config); err != nil {
			return err
		}

		configs[testName] = config
	}
	return nil
}
