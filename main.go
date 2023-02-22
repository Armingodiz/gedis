package main

import (
	"log"

	"gedis/config"

	"gedis/app"
)

func main() {

	app := app.NewApp()
	log.Fatalln(app.Start(":" + config.Configs.App.Port))
}
