package main

import (
	"flag"
	"fmt"

	"github.com/laincloud/streamrouter/dispatcher"
	"github.com/mijia/sweb/log"
)

const version = 1.0

func main() {
	var showVersion, debug bool
	flag.BoolVar(&showVersion, "v", false, "Show watcher version")
	flag.BoolVar(&debug, "-debug", false, "Open debug log")
	if showVersion {
		fmt.Printf("%f\n", version)
	} else {
		if debug {
			log.EnableDebug()
		}
		dispatcher.Run()
	}

}
