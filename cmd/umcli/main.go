// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package pgclient

import (
	"os"

	"github.com/mitchellh/cli"

	"github.com/lvl484/user-manager/logger"
)

func main() {
	c := cli.NewCLI("umcli", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		//		"foo": fooCommandFactory,
		//		"bar": barCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		logger.LogUM.Error(err)
	}

	os.Exit(exitStatus)
}
