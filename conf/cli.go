package conf

import (
	"flag"
	"log"
)

func CliArgs() string {
	confFlag := flag.Lookup("config")

	if confFlag == nil || len(confFlag.Value.String()) == 0 {
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

	return confFlag.Value.String()

}
