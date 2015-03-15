package conf

import (
	"flag"
	"log"
)

func CliArgs() string {
	if !flag.Parsed() {
		flag.String(
			"marker",
			"",
			"Marker to resume from",
		)

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

	confFlag := flag.Lookup("config")
	return confFlag.Value.String()

}
