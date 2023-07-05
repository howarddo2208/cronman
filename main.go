package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
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

func (c *Config) init() error {
	configFilePath := os.Getenv("HOME") + "/.config/cronman/cronman.yaml"
	err := readConfigFile(configFilePath, c)
	return err
}

func createConfigFile(configFilePath string) (*os.File, error) {
	configDir := path.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0770); err != nil {
		return nil, err
	}
	return os.Create(configFilePath)
}

func readConfigFile(configFilePath string, c *Config) error {
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
				return errCreating
			}
		} else {
			return errReading
		}
	}

	errParsing := viper.Unmarshal(c)
	if errParsing != nil {
		return errParsing
	}
	return nil
}

func scheduleJob(s *gocron.Scheduler, name string, job Job, stdout io.Writer) error {
	log.SetOutput(stdout)

	_, err := s.CronWithSeconds(job.Schedule).Do(func() {
		cmd := exec.Command("sh", "-c", job.Cmd)
		out, errCmd := cmd.Output()
		if errCmd != nil {
			fmt.Fprintf(os.Stderr, "%s\n", errCmd)
		} else {
			log.Printf("output of job %v: %s\n", name, out)
		}
	})
	if err != nil { // error scheduling
		return err
	}
	return nil
}

func run(ctx context.Context, c *Config, stdout io.Writer) error {
	if err := c.init(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			s := gocron.NewScheduler(time.Local)

			// Iterate over the jobs and schedule them
			for name, job := range c.Jobs {
				if errSchedule := scheduleJob(s, name, job, stdout); errSchedule != nil {
					return errSchedule
				}
			}

			// Start the scheduler
			s.StartBlocking()
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	c := &Config{}

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					c.init()
				case os.Interrupt:
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

	if err := run(ctx, c, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
