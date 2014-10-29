package main

import (
	"flag"
	"os"
)

func cliArgs(cfg *config) {
	args := map[string]*string{}

	configPtr := flag.String(
		"config",
		"",
		"Config [.ini format] file to Load the configurations from",
	)

	prefixPtr := flag.String(
		"prefix",
		"",
		"Directory prefix to process the logs for",
	)

	args["config"] = configPtr
	args["prefix"] = prefixPtr

	//Must be called after all flags are defined and
	//before flags are accessed by the program.

	flag.Parse()

	if *args["config"] == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg.MakeConfig(*args["config"])

	if *args["prefix"] != "" {
		cfg.S3.Prefix = *args["prefix"]
	}

}
