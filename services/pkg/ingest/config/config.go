package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	AccessKey      string `yaml:"access_key" env:"ACCESS_KEY"`
	AccessSecret   string `yaml:"access_secret" env:"ACCESS_SECRET"`
	ConsumerKey    string `yaml:"consumer_key" env:"CONSUMER_KEY"`
	ConsumerSecret string `yaml:"consumer_secret" env:"CONSUMER_SECRET"`
}

// Validate implements the validation interface.
func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AccessKey, validation.Required),
		validation.Field(&c.AccessSecret, validation.Required),
		validation.Field(&c.ConsumerKey, validation.Required),
		validation.Field(&c.ConsumerSecret, validation.Required),
	)
}

// ReadConfig reads config file and returns new config instance
func ReadConfig(configPath string) (*Config, error) {
	var config Config
	if configPath != "" {
		log.Infof("read config from %s", configPath)
		configFile, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Errorf("error loading app config: %v", err)
			return nil, err
		}
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			log.Errorf("error parsing app config : %v", err)
			return nil, err
		}
	} else {
		log.Infof("build config from environment variables")
		err := configor.Load(&config)
		if err != nil {
			log.Errorf("error loading app config: %v", err)
			return nil, err
		}
	}

	if err := config.Validate(); err != nil {
		log.Errorf("error validating app config: %v", err)
		return nil, err
	}

	return &config, nil
}
