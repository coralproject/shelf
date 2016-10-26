package cmdquery

import (
	"strings"

	"github.com/coralproject/shelf/cmd/xenia/web"
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
		RunE:  runExec,
	}

	cmd.Flags().StringVarP(&exe.name, "name", "n", "", "Name of Set.")
	cmd.Flags().StringVarP(&exe.vars, "vars", "v", "", "Variables required by Set.")

	queryCmd.AddCommand(cmd)
}

// runExec is the code that implements the execute command.
func runExec(cmd *cobra.Command, args []string) error {
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

	return runExecWeb(cmd, vars)
}

// runExecWeb issues the command talking to the web service.
func runExecWeb(cmd *cobra.Command, vars map[string]string) error {
	verb := "GET"
	url := "/v1/exec/" + exe.name

	if len(vars) > 0 {
		var i int
		for k, v := range vars {
			if i == 0 {
				url += "?"
			} else {
				url += "&"
			}
			i++

			url += k + "=" + v
		}
	}

	resp, err := web.Request(cmd, verb, url, nil)
	if err != nil {
		return err
	}

	cmd.Printf("\n%s\n\n", resp)
	return nil
}
