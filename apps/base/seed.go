package main

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"

	"dls/dls-core/adapters/pb"
	"dls/dls-core/seed"
)

// newSeedCommand fills a development database with a fake channel. The data
// is deterministic, so everyone running it sees the same numbers and can
// talk about the same episode.
func newSeedCommand(app *pocketbase.PocketBase) *cobra.Command {
	var reset bool

	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Fill the database with the development dataset",
		RunE: func(_ *cobra.Command, _ []string) error {
			counts, err := pb.Seed(app, seed.Build(), reset)
			if err != nil {
				return err
			}
			fmt.Println("seeded:", counts)
			return nil
		},
	}
	cmd.Flags().BoolVar(&reset, "reset", false, "delete the existing rows first")

	return cmd
}
