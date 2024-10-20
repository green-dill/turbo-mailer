package user

import (
	"context"
	"fmt"
	"os"
	"turbo-mailer-server/internal/initialize"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var Cmd = &cobra.Command{
	Use:   "user",
	Short: "User management",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updatePasswordCmd)
	Cmd.AddCommand(deleteUserCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Do(cmd.Context())
		users, err := query.User.WithContext(context.Background()).Unscoped().Find()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to list users")
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Username", "Created At", "Updated At", "Deleted"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)

		for _, user := range users {
			table.Append([]string{
				user.Username,
				user.CreatedAt.Format("2006-01-02 15:04:05"),
				user.UpdatedAt.Format("2006-01-02 15:04:05"),
				fmt.Sprintf("%v", user.DeletedAt.Valid),
			})
		}
		table.Render()
	},
}

var createCmd = &cobra.Command{
	Use:   "create <username> <password>",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Do(cmd.Context())
		username := args[0]
		password := args[1]

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to hash password")
		}

		newUser := &models.User{
			Username: username,
			Password: string(hashedPassword),
		}

		err = query.User.WithContext(context.Background()).Create(newUser)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create user")
		}

		fmt.Printf("User created successfully: ID: %d, Username: %s\n", newUser.ID, newUser.Username)
	},
}

var updatePasswordCmd = &cobra.Command{
	Use:   "update-password <username> <new-password>",
	Short: "Update a user's password",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Do(cmd.Context())
		username := args[0]
		newPassword := args[1]

		user, err := query.User.WithContext(context.Background()).Where(query.User.Username.Eq(username)).First()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to find user")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to hash password")
		}

		user.Password = string(hashedPassword)
		err = query.User.WithContext(context.Background()).Save(user)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update user password")
		}

		fmt.Printf("Password updated successfully for user: %s\n", username)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete <username>",
	Short: "Soft delete a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Do(cmd.Context())
		username := args[0]

		user, err := query.User.WithContext(context.Background()).Where(query.User.Username.Eq(username)).First()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to find user")
		}

		_, err = query.User.WithContext(context.Background()).Delete(user)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to delete user")
		}

		fmt.Printf("User %s has been soft deleted successfully\n", username)
	},
}
