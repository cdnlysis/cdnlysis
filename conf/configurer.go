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

	var confData []byte

	if path != "" {
		confData = readYaml(path)
	} else {
		//Convert to a byte array to keep it consistent with the YAML
		//parser.
		confData = []byte(defaultYaml)
	}

	conf := Config{}

	err := yaml.Unmarshal(confData, &conf)
	if err != nil {
		panic(err)
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
