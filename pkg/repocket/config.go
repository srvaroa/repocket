package repocket

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

const RepocketConfigFile = ".repocket"

type Config struct {
	ConsumerKey string `yaml:"consumer_key"`
	AccessToken string `yaml:"access_token"`
	OutputDir   string `yaml:"output_dir"`
}

func (cfg *Config) LoadConfig() error {

	usr, err := user.Current()
	if err != nil {
		log.Println("Could not determine user home %s", err)
	}

	yamlFile, err := ioutil.ReadFile(usr.HomeDir + "/" + RepocketConfigFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}

	return err
}

func (cfg *Config) SaveConfig() error {

	usr, err := user.Current()
	if err != nil {
		log.Println("Could not determine user home %s", err)
	}

	file, err := os.Create(usr.HomeDir + "/" + RepocketConfigFile)
	if err != nil {
		log.Printf("Failed to create config file at %s: %s", file, err)
		return err
	}
	defer file.Close()

	outBytes, err := yaml.Marshal(cfg)
	_, err = io.WriteString(file, string(outBytes)) // TODO: use straight bytes
	if err != nil {
		log.Printf("Failed to write config file %s: %s", file, err)
	}

	return err
}
