package main

import (
	"fmt"
	"testing"
)

// calls main.readConfigFile and checks if return Array of Job
func TestReadValidConfig(t *testing.T) {
	createConfigFile := "./test-data/config-success.yaml"
	jobs, err := readConfigFile(createConfigFile)
	if err != nil {
		t.Fatalf(`readConfigFile("./test-data/config-success.yaml") = _, %v, want _, nil`, err)
	}

	// print out the jobs
	for _, job := range jobs {
		fmt.Printf("Job: %v\n", job)
	}
}

func TestReadEmptyConfig(t *testing.T) {
}

func TestReadInvalidConfig(t *testing.T) {
}
