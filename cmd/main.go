package main

import (
	server "api_sotr/app"
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	app := new(server.App)
	app.Run("8090")
}
