package cmdview

import "github.com/spf13/cobra"

// viewCmd represents the parent for all view cli commands.
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view provides a xenia CLI for managing view metadata.",
}

// GetCommands returns the view commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	return viewCmd
}
