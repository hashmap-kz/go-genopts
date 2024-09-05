package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hashmap-kz/go-genopts/pkg/cfg"
	"github.com/hashmap-kz/go-genopts/pkg/out"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "", "config file location")
	flag.Parse()

	if configFile == "" {
		log.Fatal("expect config file")
	}

	opts := cfg.ReadInput(configFile)
	gen := out.GenOpts(opts)
	fmt.Println(gen)
}
