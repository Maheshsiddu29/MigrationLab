package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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
		Short: "Validate PostgreSQL migration syntax",
		Long:  "Discover and parse SQL migrations at a file or directory path. Migration safety analysis is not yet implemented.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			output := cmd.OutOrStdout()
			fmt.Fprintln(output, "Analyzing migrations...")
			fmt.Fprintln(output)

			migrations, err := migration.Load(args[0])
			if err != nil {
				return err
			}

			totalStatements := 0
			for _, migration := range migrations {
				statementCount := len(migration.Statements)
				totalStatements += statementCount

				fmt.Fprintln(output, filepath.Base(migration.Path))
				fmt.Fprintf(output, "  parsed %d %s\n\n", statementCount, pluralize(statementCount, "statement"))
			}

			fmt.Fprintf(output, "%d %s\n", len(migrations), pluralize(len(migrations), "migration"))
			fmt.Fprintf(output, "%d %s\n", totalStatements, pluralize(totalStatements, "statement"))

			return nil
		},
	})

	return root
}

func pluralize(count int, singular string) string {
	if count == 1 {
		return singular
	}

	return singular + "s"
}
