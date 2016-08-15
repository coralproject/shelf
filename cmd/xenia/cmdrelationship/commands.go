package cmdrelationship

import "github.com/spf13/cobra"

// relationshipCmd represents the parent for all relationship cli commands.
var relationshipCmd = &cobra.Command{
	Use:   "relationship",
	Short: "relationship provides a xenia CLI for managing relationship metadata.",
}

// GetCommands returns the relationship commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	return relationshipCmd
}
