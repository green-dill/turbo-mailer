package main

import (
	"turbo-mailer-server/internal/models"

	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "internal/query",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.ApplyBasic(models.ALL...)

	g.WithDbNameOpts(func(d *gorm.DB) string {
		return d.Name()
	})

	g.Execute()
}
