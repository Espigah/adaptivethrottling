package main

import (
	"fmt"

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
			return c.SendStatus(500)
		}

		return c.SendStatus(200)
	})

	app.Listen(":3000")
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
