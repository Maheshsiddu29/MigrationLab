package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Maheshsiddu29/MigrationLab/internal/migration"
	"github.com/spf13/cobra"
)

const version = "dev"

func main() {
	if err := newRootCommand(os.Stdout, os.Stderr).Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCommand(stdout, stderr io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:               "migrationlab",
		Short:             "Analyze PostgreSQL migrations before deployment",
		SilenceErrors:     true,
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the MigrationLab version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "migrationlab %s\n", version)
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "analyze <path>",
		Short: "Discover SQL migration files",
		Long:  "Discover SQL migration files at a file or directory path. SQL parsing and risk analysis are not yet implemented.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			migrations, err := migration.Discover(args[0])
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Discovered migrations:")
			for _, path := range migrations {
				fmt.Fprintln(cmd.OutOrStdout(), path)
			}

			return nil
		},
	})

	return root
}
