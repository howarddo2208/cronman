package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/howarddo2208/cronman/models"
)

func InitConfig(c *models.Config) error {
	configFilePath := os.Getenv("HOME") + "/.config/cronman/cronman.yaml"
	setViperConfig(configFilePath)

	err := readConfigFile(configFilePath, c)
	return err
}

func setViperConfig(configFilePath string) {
	configDir := path.Dir(configFilePath)
	fullName := path.Base(configFilePath)
	nameWithoutExt := fullName[0 : len(fullName)-len(filepath.Ext(fullName))]
	viper.SetConfigName(nameWithoutExt)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
}

func createConfigFile(configFilePath string) (*os.File, error) {
	configDir := path.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0770); err != nil {
		return nil, err
	}
	return os.Create(configFilePath)
}

func readConfigFile(configFilePath string, c *models.Config) error {
	configDir := path.Dir(configFilePath)
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
