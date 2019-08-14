package cmd

import (
	"github.com/spf13/cobra"
	"go.fraixed.es/storj/sntools/dbcleaner/internal"
)

type globalArgs struct {
	dbFile string
	debug  bool
}

// Execute executes the command line logic.
func Execute() error {
	return rootCmd().Execute()
}

// rootCmd returns the root command.
func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dbcleaner",
		Short: "Perform clean up DB operations",
		Long:  "Performs some clean up DB maintenance operations on a Storage Node",
	}

	// Map global flags
	args := &globalArgs{}
	cmd.PersistentFlags().StringVarP(
		&args.dbFile, "dbfile", "f", "", "SQLite database filepath (required)",
	)
	// TODO: investigate about this error and if needed, handle it
	_ = cmd.MarkPersistentFlagRequired("dbfile")

	cmd.PersistentFlags().BoolVarP(
		&args.debug, "debug", "d", false, "enable/disable debug messages (default: disabled)",
	)

	// Add subcommands
	cmd.AddCommand(unsentInvalidLimitsCmd(args))

	return cmd
}

func unsentInvalidLimitsCmd(gargs *globalArgs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unsent-invalid-limits",
		Short: "Delete unsent invalid limits",
		Long: "" +
			"Delete all the unsent orders which have an invalid order limits.\n" +
			"Currently it only deletes those order with an order limit which " +
			"cannot be unmarshal due to some Storj validations.",
		RunE: func(c *cobra.Command, args []string) error {
			return internal.UnsentInvalidLimits(gargs.dbFile, gargs.debug)
		},
	}

	return cmd
}
