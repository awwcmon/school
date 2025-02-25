// Package main is the http server of the application.
package main

import (
	"school/cmd/school/initial"

	"github.com/go-dev-frame/sponge/pkg/app"
)

func main() {
	initial.InitApp()
	services := initial.CreateServices()
	closes := initial.Close(services)
	initial.InitSnailJob()
	a := app.New(services, closes)
	a.Run()
}

// define Person and Persons
type Person struct {
	ID   int
	Name string
}

type Persons []Person

// implement sort interface
func (p Persons) Len() int {
	return len(p)
}

func (p Persons) Less(i, j int) bool {
	return p[i].ID < p[j].ID
}

func (p Persons) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
