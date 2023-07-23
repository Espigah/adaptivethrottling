package testers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const api = "http://localhost:3003/"

type Config struct {
	FirstPointOfFailure int `json: "firstPointOfFailure"`
	Intermittency       int `json: "intermittency"`
}

func get(testName string, body *Config) (interface{}, error) {
	url := api + testName

	// marshall data to json (like json_encode)
	marshalled, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("impossible to marshall teacher: %s", err)
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(marshalled))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode > 200 {
		return nil, errors.New("error")
	}

	return res, err
}
