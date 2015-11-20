package cmddb

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var createLong = `Creates a new database with collections and indexes.

Example:

	db create -f ./scripts/database.json
`

// create contains the state for this command.
var create struct {
	file string
}

// addCreate handles the creation of users.
func addCreate() {
	cmd := &cobra.Command{
		Use:   "create [-f file]",
		Short: "Creates a new database from a script file",
		Long:  createLong,
		Run:   runCreate,
	}

	cmd.Flags().StringVarP(&create.file, "file", "f", "", "file path of script json file")

	dbCmd.AddCommand(cmd)
}

//==============================================================================

// DB is the container for all db objects.
type DB struct {
	Cols []Collection `json:"collections"`
}

// Collection is the container for a db collection definition.
type Collection struct {
	Name    string  `json:"name"`
	Indexes []Index `json:"indexes"`
}

// Index is the container for an index definition.
type Index struct {
	Name     string  `json:"name"`
	IsUnique bool    `json:"unique"`
	Fields   []Field `json:"fields"`
}

// Field is the container for a field definition.
type Field struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	OtherType string `json:"other"`
}

// runCreate is the code that implements the create command.
func runCreate(cmd *cobra.Command, args []string) {
	cols, err := getCollections(create.file)
	if err != nil {
		dbCmd.Printf("Error reading collections : %s : ERROR : %v\n", create.file, err)
		return
	}

	dbCmd.Println(cols)
}

// getCollections reads the specified file and returns the set of
// collection definitions that are defined.
func getCollections(fileName string) ([]Collection, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var db DB
	if err := json.NewDecoder(f).Decode(&db); err != nil {
		return nil, err
	}

	return db.Cols, nil
}
