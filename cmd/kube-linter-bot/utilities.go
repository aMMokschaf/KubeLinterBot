//Package main reads config files, and contains the hook-receiving server.
package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// type config struct {
// 	Repository struct {
// 		RepoName string `yaml:"reponame"`
// 		User     struct {
// 			Username    string `yaml:"username"`
// 			AccessToken string `yaml:"accessToken"`
// 		}
// 	}
// 	Bot struct {
// 		Port int `yaml:"port"`
// 	}
// 	Webhook struct {
// 		Secret string `yaml:"secret"`
// 	}
// }

type config struct {
	User struct {
		Username string `yaml:"username"`
	}
	Bot struct {
		Port int `yaml:"port"`
	}
	Repositories []struct {
		Name        string `yaml:"name"`
		AccessToken string `yaml:"accessToken"`
		Owner       string `yaml:"owner"`
		Webhook     struct {
			Secret  string   `yaml:"secret"`
			Events  []string `yaml:"events"`
			Address string   `yaml:"address"`
		} `yaml:"repositories"`
	}
}

var cfg config

//optionParser reads a config-file named "kube-linter-bot-configuration.yaml", that has
//to be located in the same folder as kube-linter-bot and parses its contents to a struct.
//A sample file is located in /samples/
func optionParser() config {
	dat, err := ioutil.ReadFile("kube-linter-bot-configuration.yaml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal([]byte(dat), &cfg)
	return cfg
}

//writeOptionsToFile saves changes to the configuration to kube-linter-bot-configuration.yaml.
func writeOptionsToFile() bool {
	status := false

	d, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile("./kube-linter-bot-configuration.yaml", d, 0666) //TODO: Check permissions
	if err != nil {
		panic(err)
	} else {
		status = true
	}
	return status
}
