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

	args["config"] = configPtr

	//Must be called after all flags are defined and
	//before flags are accessed by the program.

	flag.Parse()

	if *args["config"] == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg.MakeConfig(*args["config"])
}
