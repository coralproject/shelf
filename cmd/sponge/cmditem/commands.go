package cmditem

import "github.com/spf13/cobra"

// itemCmd represents the parent for all item cli commands.
var itemCmd = &cobra.Command{
	Use:   "item",
	Short: "item provides a sponge CLI for managing items.",
}

// GetCommands returns the item commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	return itemCmd
}
