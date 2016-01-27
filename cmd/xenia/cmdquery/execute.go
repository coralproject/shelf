package cmdquery

import (
	"encoding/json"
	"strings"

	"github.com/coralproject/xenia/pkg/exec"
	"github.com/coralproject/xenia/pkg/query"

	"github.com/spf13/cobra"
)

var execLong = `Executes a Set from the system by the sets name.

Example:
	query exec -n "user_advice"

	query exec -n "my_set" -v "key:value,key:value"
`

// exe contains the state for this command.
var exe struct {
	name string
	vars string
}

// addExec handles the execution of queries.
func addExec() {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Executes a Set by name.",
		Long:  execLong,
		Run:   runExec,
	}

	cmd.Flags().StringVarP(&exe.name, "name", "n", "", "Name of Set.")
	cmd.Flags().StringVarP(&exe.vars, "vars", "v", "", "Variables required by Set.")

	queryCmd.AddCommand(cmd)
}

// runExec is the code that implements the execute command.
func runExec(cmd *cobra.Command, args []string) {
	cmd.Printf("Exec Set : Name[%s] Vars[%v]\n", exe.name, exe.vars)

	if exe.name == "" {
		cmd.Help()
		return
	}

	set, err := query.GetByName("", conn, exe.name)
	if err != nil {
		cmd.Println("Exec Set : ", err)
		return
	}

	vars := make(map[string]string)
	if exe.vars != "" {
		vs := strings.Split(exe.vars, ",")
		for _, kvs := range vs {
			kv := strings.Split(kvs, ":")
			if len(kv) != 2 {
				continue
			}
			vars[kv[0]] = kv[1]
		}
	}

	result := exec.Exec("", conn, set, vars)

	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		cmd.Println("Exec Set : ", err)
		return
	}

	cmd.Printf("\n%s\n\n", string(data))
}
