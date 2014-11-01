package main

import (
	"flag"
	"os"
)

func cliArgs(cfg *config) {
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

	args["config"] = configPtr
	args["prefix"] = prefixPtr

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

}
