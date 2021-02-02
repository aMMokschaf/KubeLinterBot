//Package config reads and writes config-files.
package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//Config is the representation of kube-linter-bot-configuration.yaml as a struct.
type Config struct {
	User struct {
		Username    string `yaml:"username"`
		AccessToken string `yaml:"accessToken"`
	}
	Bot struct {
		Port int `yaml:"port"`
	}
	Repositories []struct {
		Name    string `yaml:"name"`
		Owner   string `yaml:"owner"`
		Webhook struct {
			Secret  string   `yaml:"secret"`
			Events  []string `yaml:"events"`
			Address string   `yaml:"address"`
		} `yaml:"webhook"`
	} `yaml:"repositories"`
}

//OptionParser reads a config-file named "kube-linter-bot-configuration.yaml", that has
//to be located in the same folder as kube-linter-bot and parses its contents to a struct.
//A sample file is located in /samples/
func OptionParser() (*Config, error) {
	var cfg Config
	dat, err := ioutil.ReadFile("kube-linter-bot-configuration.yaml")
	if err != nil {
		//panic(err)
		return nil, err
	}
	err = yaml.Unmarshal([]byte(dat), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

//WriteOptionsToFile saves changes to the configuration to kube-linter-bot-configuration.yaml.
func WriteOptionsToFile(cfg Config) error {
	d, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile("./kube-linter-bot-configuration.yaml", d, 0666) //TODO: Check permissions
	if err != nil {
		panic(err)
	}
	return err
}
