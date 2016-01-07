package cmdquery

import (
	"encoding/json"
	"strings"

	"github.com/coralproject/xenia/pkg/query"

	"github.com/ardanlabs/kit/db"
	"github.com/spf13/cobra"
)

var execLong = `Executes a query from the system by the query name.

Example:
	query exec -n "user_advice"

	query exec -n "my_query" -v "key:value,key:value"
`

// exec contains the state for this command.
var exec struct {
	name string
	vars string
}

// addExec handles the execution of queries.
func addExec() {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Executes a query by name.",
		Long:  execLong,
		Run:   runExec,
	}

	cmd.Flags().StringVarP(&exec.name, "name", "n", "", "Name of query.")
	cmd.Flags().StringVarP(&exec.vars, "vars", "v", "", "Variables required by query.")

	queryCmd.AddCommand(cmd)
}

// runExec is the code that implements the execute command.
func runExec(cmd *cobra.Command, args []string) {
	cmd.Printf("Exec Query : Name[%s] Vars[%v]\n", exec.name, exec.vars)

	if exec.name == "" {
		cmd.Help()
		return
	}

	db := db.NewMGO()
	defer db.CloseMGO()

	set, err := query.GetSetByName("", db, exec.name)
	if err != nil {
		cmd.Println("Exec Query : ", err)
		return
	}

	vars := make(map[string]string)
	if exec.vars != "" {
		vs := strings.Split(exec.vars, ",")
		for _, kvs := range vs {
			kv := strings.Split(kvs, ":")
			if len(kv) != 2 {
				continue
			}
			vars[kv[0]] = kv[1]
		}
	}

	result := query.ExecuteSet("", db, set, vars)

	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		cmd.Println("Exec Query : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
}
