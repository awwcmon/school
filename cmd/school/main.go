// Package main is the http server of the application.
package main

import (
	"github.com/go-dev-frame/sponge/pkg/app"

	"school/cmd/school/initial"
)

func main() {
	initial.InitApp()
	services := initial.CreateServices()
	closes := initial.Close(services)
	a := app.New(services, closes)
	a.Run()
}
