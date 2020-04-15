// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"github.com/urfave/cli/v2"
)

var actionHelper *actionHandle

// CreateAction create new user in database
func CreateAction(c *cli.Context) error {
	return actionHelper.ExecuteAction(c, actionCreate)
}

// InfoAction get info of user by its login name
func InfoAction(c *cli.Context) error {
	return actionHelper.ExecuteAction(c, actionInfo)
}

// ActivateAction activate user that was disabled
func ActivateAction(c *cli.Context) error {
	return actionHelper.ExecuteAction(c, actionActivate)
}

// DisableAction disable user by its login
func DisableAction(c *cli.Context) error {
	return actionHelper.ExecuteAction(c, actionDisable)
}

// UpdateAction update information about user
func UpdateAction(c *cli.Context) error {
	return actionHelper.ExecuteAction(c, actionUpdate)
}

// DeleteAction delete user in database
func DeleteAction(c *cli.Context) error {
	return actionHelper.ExecuteAction(c, actionDelete)
}
