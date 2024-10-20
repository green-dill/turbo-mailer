package migrate

import (
	"context"
	"errors"
	"time"
	"turbo-mailer-server/internal/initialize"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("migrate database")
		initialize.Do(cmd.Context())
		start := time.Now()
		if err := query.Migrate(); err != nil {
			log.Fatal().Err(err).Msg("migrate failed")
		}
		log.Info().Dur("duration", time.Since(start)).Msg("migrate database success")

		if err := createDefaultUser(); err != nil {
			log.Fatal().Err(err).Msg("create default user failed")
		}
	},
}

func createDefaultUser() error {
	// Check if admin user exists (including soft-deleted)
	_, err := query.User.WithContext(context.Background()).Unscoped().Where(query.User.Username.Eq("admin")).First()
	if err == nil {
		log.Info().Msg("Admin user already exists")
		return nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create default admin user if not exists
	defaultPassword := "Tm@123qwe"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &models.User{
		Username: "admin",
		Password: string(hashedPassword),
	}

	err = query.User.WithContext(context.Background()).Create(newUser)
	if err != nil {
		return err
	}

	log.Info().Msg("Default admin user created successfully")
	return nil
}
