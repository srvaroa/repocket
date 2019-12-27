package repocket

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

const RepocketConfigFile = ".repocket"

type Config struct {
	ConsumerKey string `required:"true" split_words:"true"`
	AccessToken string
	OutputDir   string `split_words:"true"`
}

func LoadLocalConfig() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Println("Could not determine user home %s", err)
	}
	b, err := ioutil.ReadFile(usr.HomeDir + "/" + RepocketConfigFile)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func SaveLocalConfig(token string) error {
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

	_, err = io.WriteString(file, token)
	if err != nil {
		log.Printf("Failed to write config file %s: %s", file, err)
	}

	return err
}
