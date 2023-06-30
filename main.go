package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
)

func hello(name string) {
	message := fmt.Sprintf("Hello, %s!", name)
	fmt.Println(message)
}

func runCronJobs() {
	s := gocron.NewScheduler(time.Local)

	s.Every(4).Second().Do(func() {
		hello("Le Van Dat")
	})

	s.StartBlocking()
}

var (
	configDir  = os.Getenv("HOME") + "/.config/cronman"
	configFile = configDir + "/cronman.yaml"
)

func createConfigFile() (*os.File, error) {
	if err := os.MkdirAll(configDir, 0770); err != nil {
		return nil, err
	}
	return os.Create(configFile)
}

func main() {
	// load config file
	viper.SetConfigName("cronman") // name of config file (without extension)
	viper.SetConfigType(
		"yaml",
	) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configDir) // call multiple times to add many search paths
	err := viper.ReadInConfig()    // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; create new file at the path
			fmt.Println("Config file not found; create new file at the path")
			_, err2 := createConfigFile()
			if err2 != nil {
				panic(fmt.Errorf("fatal error config file: %w", err2))
			}
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	// TODO: parse config file

	// runCronJobs()
}
