package conf

import (
	"flag"
	"os"
)

var Settings *Config

func CliArgs() {
	var cfg Config

	var config = flag.String(
		"config",
		"",
		"Config [.ini format] file to Load the configurations from",
	)

	var prefix = flag.String(
		"prefix",
		"",
		"Directory prefix to process the logs for",
	)

	var verbose = flag.Bool(
		"v",
		true,
		"Display activity progress. Errors are not supressed",
	)

	//Must be called after all flags are defined and
	//before flags are accessed by the program.

	flag.Parse()

	if *config == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg.MakeConfig(*config)

	if *prefix != "" {
		cfg.S3.Prefix = *prefix
	}

	cfg.Verbose = *verbose
	Settings = &cfg
}
