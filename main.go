package main

import (
	"github.com/armzerpa/ABC-api-test/cmd"
	"github.com/armzerpa/ABC-api-test/cmd/config"
)

func main() {
	config := config.GetConfig()

	app := &cmd.App{}
	app.Initialize(config)
}
