package main

import (
	"flag"
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/laincloud/streamrouter/watcher"
)

const version = 1.0

func main() {
	var showVersion, debug bool
	flag.BoolVar(&showVersion, "v", false, "Show watcher version")
	flag.BoolVar(&debug, "-debug", false, "Open debug log")
	if showVersion {
		fmt.Println("%f\n", version)
	} else {
		if debug {
			log.EnableDebug()
		}
		log.Fatal(watcher.Run().Error())
	}

}