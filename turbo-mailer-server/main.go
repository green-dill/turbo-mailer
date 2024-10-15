//go:generate go run internal/models/generator/generator.go
//go:generate swag init
//go:generate swag fmt
package main

import (
	"turbo-mailer-server/cmd"
	_ "turbo-mailer-server/docs"
)

// @title		Turbo Mailer API
// @version	0.1.0
// @BasePath	/api
func main() {
	cmd.Execute()
}
