package main

import (
	"testing"
)

func TestDefaultServiceBehaviour(t *testing.T) {
	config := []Service{
		{Name: "api-gateway", Url: "https://htt.pavonz.com/200", TimeoutMS: 5000},
		{Name: "user-service", Url: "https://htt.pavonz.com/500", TimeoutMS: 3000},
		{Name: "payment-service", Url: "https://htt.pavonz.com/503", TimeoutMS: 5000},
	}

	results := CheckServices(config)
	status := Aggregate(results)

	if status != "degraded" {
		t.Errorf("Expected degraded, got %s", status)
	}
}

func TestTimeoutHandling(t *testing.T) {
	config := []Service{
		{Name: "api-gateway", Url: "https://htt.pavonz.com/200", TimeoutMS: 10},
	}

	results := CheckServices(config)

	if results[0].Status != "down" {
		t.Errorf("Expected status down, got %s", results[0].Status)
	}

	if results[0].Error == "" {
		t.Errorf("Expected an error due to timeout, got none")
	}
}

func TestAggregate(t *testing.T) {
	tests := []struct {
		name     string
		results  []Result
		expected string
	}{
		{"all healthy", []Result{{Status: "healthy"}, {Status: "healthy"}}, "healthy"},
		{"all down", []Result{{Status: "down"}, {Status: "down"}}, "down"},
		{"mixed", []Result{{Status: "healthy"}, {Status: "down"}}, "degraded"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Aggregate(test.results)
			if got != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, got)
			}
		})
	}
}
