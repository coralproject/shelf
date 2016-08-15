package cmdmask

import "github.com/spf13/cobra"

// maskCmd represents the parent for all mask cli commands.
var maskCmd = &cobra.Command{
	Use:   "mask",
	Short: "mask provides a xenia CLI for managing and executing masks.",
}

// GetCommands returns the mask commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	return maskCmd
}
