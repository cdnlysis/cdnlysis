package conf

import (
	"flag"
	"log"
)

func cliArgs() string {
	var config = flag.String(
		"config",
		"",
		"Config [.yml format] file to load the configurations from",
	)

	flag.Parse()

	if *config == "" {
		log.Println("No config file supplied. Using defauls.")
	}

	return *config
}
