package main

import (
	"context"
	"idp/Jobs/api"
	"idp/Jobs/config"
	"log"
)

const (
	confFile="config.json"
)

func main() {
	ctx := context.Background()
	conf, err := config.NewMicroserviceConfig(confFile).GetConfig()
	if err != nil {
		log.Println("Could not get server config: " + err.Error())
		panic(1)
	}

	if err := api.RunServer(ctx, conf); err != nil {
		log.Println("Could not start server: " + err.Error())
		panic(2)
	}
}
