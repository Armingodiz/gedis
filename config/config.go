package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		App      App
		Database Database
		Cache    Cache
	}

	App struct {
		Port      string `envconfig:"APP_PORT" default:"3000"`
		JwtSecret string `envconfig:"JWT_SECRET"`
	}

	Database struct {
		Host     string `envconfig:"DATABASE_HOST"`
		Port     int    `envconfig:"DATABASE_PORT"`
		User     string `envconfig:"DATABASE_USER"`
		Password string `envconfig:"DATABASE_PASSWORD"`
		DbName   string `envconfig:"DATABASE_DBNAME"`
		Extras   string `envconfig:"DATABASE_EXTRAS"`
		Driver   string `envconfig:"DATABASE_DRIVER" default:"postgres"`
	}

	Cache struct {
		CacheAddr string `envconfig:"CACHE_ADDR"`
	}
)

var Configs Config

func init() {
	err := envconfig.Process("", &Configs)
	fmt.Println(Configs)
	if err != nil {
		log.Fatal(err.Error())
	}
}
