package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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
	var wg sync.WaitGroup
	ch := make(chan Result)

	for _, service := range services {
		wg.Add(1)

		go func(svc Service) {
			defer wg.Done()
			ch <- check(service)
		}(service)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var results []Result
	for result := range ch {
		results = append(results, result)
	}

	return results
}

func check(service Service) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(service.TimeoutMS)*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", service.Url, nil)
	if err != nil {
		return Result{Name: service.Name, Status: "down", Error: err.Error()}
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return Result{Name: service.Name, Status: "down", Error: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return Result{Name: service.Name, Status: "healthy", ResponseTimeMS: duration}
	}

	return Result{Name: service.Name, Status: "down", Error: fmt.Sprintf("HTTP %d", resp.StatusCode)}
}

func Aggregate(results []Result) string {
	return "healthy"
}
