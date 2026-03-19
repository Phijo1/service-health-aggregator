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