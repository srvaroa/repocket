package repocket

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/srvaroa/repocket/pkg/pocket"
)

const RepocketConfigFile = ".repocket/config"

type Config struct {
	ConsumerKey string `yaml:"consumer_key"`
	AccessToken string `yaml:"access_token"`
	FavsDir     string `yaml:"favs_dir"`
	UnreadDir   string `yaml:"unread_dir"`
}

func (cfg *Config) Load() error {

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

func (cfg *Config) Save() error {

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

// Authenticate will ensure that a given Config object is autheticated,
// either by providing a ConsumerKey and AccessToken, or running the
// oauth auth process.  In the latter case, the token is persisted in
// the config file.
func (cfg *Config) Authenticate() {
	if len(cfg.AccessToken) > 0 {
		return
	}
	if len(cfg.ConsumerKey) <= 0 {
		log.Fatalf("Your config file seems empty.  It should contain " +
			"at least an entry with the consumer_key.  Please check the " +
			"README.md for details")
	}

	log.Printf("Loading access token..")

	var err error
	cfg.AccessToken, err = pocket.Authorize(cfg.ConsumerKey)
	if err != nil {
		log.Fatal("Failed to authorize against Pocket: %s", err)
	}

	err = cfg.Save()
	if err != nil {
		log.Printf("Failed to persist user token: %s", err)
	}
}
