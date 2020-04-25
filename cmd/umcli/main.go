// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

const (
	nameAPP  = "umcli"
	usageAPP = "Command line tool to manipulate accounts with admin rights"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = nameAPP
	app.Usage = usageAPP

	app.Commands = umcliCommands

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
