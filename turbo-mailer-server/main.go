//go:generate go run internal/models/generator/generator.go
//go:generate swag init --parseVendor --parseDependency --parseInternal
//go:generate swag fmt
package main

import (
	"turbo-mailer-server/cmd"
	_ "turbo-mailer-server/docs"
)

// @title		Turbo Mailer API
// @version	0.1.0
// @BasePath	/
func main() {
	cmd.Execute()
}
