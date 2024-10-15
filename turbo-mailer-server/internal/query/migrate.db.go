package query

import "turbo-mailer-server/internal/models"

func Migrate() error {
	return DB.AutoMigrate(models.ALL...)
}
