package main

import (
	"flag"

	"github.com/Beretta350/authentication/internal/app"
)

func main() {
	var env string
	flag.StringVar(&env, "env", "local", "Specify enviroment. Default is local")
	flag.Parse()

	app.Run(env)
}
