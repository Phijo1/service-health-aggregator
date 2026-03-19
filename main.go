package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Result struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	ResponseTimeMS int64  `json:"response_time_ms,omitempty"`
	Error          string `json:"error,omitempty"`
}

type AggregateResponse struct {
	Status    string   `json:"status"`
	Timestamp string   `json:"timestamp"`
	Services  []Result `json:"services"`
}

type Service struct {
	Name      string `yaml:"name"`
	Url       string `yaml:"url"`
	TimeoutMS int    `yaml:"timeout_ms"`
}

type Config struct {
	Services []Service `yaml:"services"`
}

func main() {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health/aggregate", func(w http.ResponseWriter, r *http.Request) {
		results := CheckServices(config.Services)
		status := Aggregate(results)

		response := AggregateResponse{
			Status:    status,
			Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			Services:  results,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Running on :8080, see http://127.0.0.1:8080/health/aggregate")
	http.ListenAndServe(":8080", nil)
}

func CheckServices(services []Service) []Result {
	var results []Result
	result := Result{
		Name:           "example",
		Status:         "foo",
		ResponseTimeMS: 100,
	}

	results = append(results, result)
	return results
}

func Aggregate(results []Result) string {
	return "healthy"
}
