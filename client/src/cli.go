package client

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

var getAgeCmd = &cobra.Command{
	Use:           "get-age <name>...",
	Short:         "Gets the age of a person by name",
	ValidArgs:     []string{"name"},
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          GetAge,
}

var createUserCmd = &cobra.Command{
	Use:           "cr-user <name> <age>",
	Short:         "Creates user in server",
	Args:          cobra.ExactArgs(2),
	ValidArgs:     []string{"name", "age"},
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          CreateUser,
}

var deleteUserCmd = &cobra.Command{
	Use:           "del-user <name>",
	Short:         "Deletes user from server",
	Args:          cobra.ExactArgs(1),
	ValidArgs:     []string{"name"},
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          DeleteUser,
}

var getAppCmd = &cobra.Command{
	Use:           "get-app <name>",
	Short:         "Gets the app info",
	Args:          cobra.ExactArgs(1),
	ValidArgs:     []string{"name"},
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          GetApp,
}

var createAppCmd = &cobra.Command{
	Use:           "cr-app <name> <price>",
	Short:         "Creates app in server",
	Args:          cobra.ExactArgs(2),
	ValidArgs:     []string{"name", "price"},
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          CreateApp,
}

var deleteAppCmd = &cobra.Command{
	Use:           "del-app <name>",
	Short:         "Deletes app from server",
	Args:          cobra.ExactArgs(1),
	ValidArgs:     []string{"name"},
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          DeleteApp,
}

func Execute() {
	rootCmd.AddCommand(getAgeCmd)
	rootCmd.AddCommand(createUserCmd)
	rootCmd.AddCommand(deleteUserCmd)
	rootCmd.AddCommand(getAppCmd)
	rootCmd.AddCommand(createAppCmd)
	rootCmd.AddCommand(deleteAppCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("There was an error while executing your CLI: %v", err)
	}
}
