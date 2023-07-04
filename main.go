package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

type Job struct {
	Cmd      string
	Schedule string
}

type Config struct {
	Jobs map[string]Job
}

func createConfigFile(configFilePath string) (*os.File, error) {
	configDir := path.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0770); err != nil {
		return nil, err
	}
	return os.Create(configFilePath)
}

func readConfigFile(configFilePath string) (map[string]Job, error) {
	// load config file
	configDir := path.Dir(configFilePath)
	fullName := path.Base(configFilePath)
	nameWithoutExt := fullName[0 : len(fullName)-len(filepath.Ext(fullName))]
	viper.SetConfigName(nameWithoutExt)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	errReading := viper.ReadInConfig()
	if errReading != nil { // Handle errors reading the config file
		if _, ok := errReading.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found; create new file at the path")
			_, errCreating := createConfigFile(configDir)
			if errCreating != nil {
				return nil, errCreating
			}
		} else {
			return nil, errReading
		}
	}

	var config Config

	errParsing := viper.Unmarshal(&config)
	if errParsing != nil {
		return nil, errParsing
	}

	return config.Jobs, nil
}

func runCronJobs(jobs map[string]Job) {
	// Create a scheduler instance
	s := gocron.NewScheduler(time.Local)

	// Iterate over the jobs and schedule them
	for name, job := range jobs {
		scheduleJob(s, name, job)
	}

	// Start the scheduler
	s.StartBlocking()
}

func scheduleJob(s *gocron.Scheduler, name string, job Job) {
	_, err := s.CronWithSeconds(job.Schedule).Do(func() {
		cmd := exec.Command("sh", "-c", job.Cmd)
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error running command: %s\n", err)
			return
		}
		fmt.Println(string(out))
	})
	if err != nil {
		fmt.Printf("Failed to schedule job %s: %s\n", name, err)
	}
}

func main() {
	configFilePath := os.Getenv("HOME") + "/.config/cronman/cronman.yaml"
	jobs, err := readConfigFile(configFilePath)
	if err != nil {
		panic(err)
	}

	runCronJobs(jobs)
}
