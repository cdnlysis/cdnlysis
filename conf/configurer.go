package conf

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var Settings Config

func GetConfig() Config {
	//Try to look for a module level cached value
	if Settings.initalized {
		return Settings
	}

	var path string = cliArgs()

	conf := Config{}
	yaml.Unmarshal([]byte(defaultYaml), &conf)

	if path != "" {
		confData := readYaml(path)
		yaml.Unmarshal(confData, &conf)
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
