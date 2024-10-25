package main

import (
	"log"

	env "github.com/the-witcher-knight/envconfig"
)

type Config struct {
	Name    string `env:"name,required"`
	Email   string `env:"email,required"`
	Gender  string `env:"gender,expectedValue=male female"`
	Enabled bool   `env:"enabled"`
}

func main() {
	// Setup env
	// export name=the-knight
	// export email=knight@witchertown.com
	// export gender=male

	var cfg Config
	if err := env.Lookup(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", cfg)
}
