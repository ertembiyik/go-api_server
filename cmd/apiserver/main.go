package main

import (
	"flag"
	"log"
	apiserver "webserver/internal/app/apiservser"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	s := apiserver.New(config)

	if err != nil {
		log.Fatal(err)
	}

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}







