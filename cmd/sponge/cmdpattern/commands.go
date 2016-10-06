package cmdpattern

import "github.com/spf13/cobra"

// patternCmd represents the parent for all pattern cli commands.
var patternCmd = &cobra.Command{
	Use:   "pattern",
	Short: "pattern provides a xenia CLI for managing pattern metadata.",
}

// GetCommands returns the pattern commands.
func GetCommands() *cobra.Command {
	addUpsert()
	addGet()
	addDel()
	return patternCmd
}
