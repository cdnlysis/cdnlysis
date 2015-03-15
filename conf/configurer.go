package conf

import (
	"flag"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

var Settings Config

func GetConfig() Config {
	//Try to look for a module level cached value
	if Settings.initalized {
		return Settings
	}

	var path string = CliArgs()

	conf := Config{}
	err := yaml.Unmarshal([]byte(defaultYaml), &conf)
	if err != nil {
		log.Println(err)
	}

	if path != "" {
		confData := readYaml(path)
		yaml.Unmarshal(confData, &conf)
	}

	prefix := flag.Lookup("prefix").Value.String()
	if prefix != "" {
		conf.S3.Prefix = prefix
	}

	//Set the module level cached value.
	Settings = conf
	return Settings
}

func readYaml(path string) []byte {
	//Load YAML file from the path provided.
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return yamlFile
}
