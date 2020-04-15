// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"github.com/urfave/cli/v2"
)

const (
	splitDelimiter = "="
	argLogin       = "login"

	msgConnvertParam       = "There is no such field to assigned value"
	msgConverteUser        = "Converting to structure User failed"
	msgErrorReadConfig     = "Reading configuration failed"
	msgErrorDBConfig       = "Creating database configuration failed"
	msgErrorConnectDB      = "Connecting to database failed"
	msgEmptyInputArguments = "No arguments for command line"
	msgWrongInputArguments = "Arguments for command line are wrong"
	msgUserActivated       = "User activated successfully"
	msgUserDeleted         = "User deleted successfully"
	msgUserDisable         = "User disabled successfully"
	msgErrorActionInput    = "No such command for execute"

	msgCreate   = "Creating user "
	msgUpdate   = "Updating user "
	msgGetInfo  = "GetInfo user "
	msgDelete   = "Deleting user "
	msgDisable  = "Disabling user "
	msgActivate = "Activating user "

	msgFormat = "%s <<%v>>"
)

const (
	actionCreate = iota
	actionUpdate
	actionInfo
	actionDelete
	actionActivate
	actionDisable
)

var convertFieldName = func() map[string]string {
	return map[string]string{
		`login`:    `Username`,
		`pwd`:      `Password`,
		`email`:    `Email`,
		`phone`:    `Phone`,
		`name`:     `FirstName`,
		`lastname`: `LastName`,
	}
}

// umcliCommands commands of CUI interface
var umcliCommands = cli.Commands{
	{
		Name:      "create",
		Usage:     "Create new account in database",
		ArgsUsage: "login=mylogin pwd=ttrtyrghgfh email=boss@company.com phone=7778777778877887 name=Borys lastname=Petrenko ",
		Action:    CreateAction,
		Description: `Create new account, all fields are obligatory
		login=mylogin 
		pwd=ttrtyrghgfh 
		email=boss@company.com 
		phone=7778777778877887 
		name=Borys 
		lastname=Petrenko `,
	},
	{
		Name:      "delete",
		Usage:     "Delete account in database",
		ArgsUsage: "login=mylogin",
		Action:    DeleteAction,
		Description: `Delete account by user name (login)
		login=mylogin `,
	},
	{
		Name:      "disable",
		Usage:     "Disable account in database",
		ArgsUsage: "login=mylogin",
		Action:    DisableAction,
		Description: `Disable account by user name (login), without deleting it
		login=mylogin `,
	},
	{
		Name:      "activate",
		Usage:     "Activate previously disabled account in database",
		ArgsUsage: "login=mylogin",
		Action:    ActivateAction,
		Description: `Activate account by user name (login), that was disabled
		login=mylogin `,
	},
	{
		Name:      "update",
		Usage:     "Update account in database",
		ArgsUsage: "login=mylogin pwd=ttrtyrghgfh email=boss@company.com phone=7778777778877887 name=Borys lastname=Petrenko ",
		Action:    UpdateAction,
		Description: `Update account, obligatory field is login
		login=mylogin 
		pwd=ttrtyrghgfh 
		email=boss@company.com 
		phone=7778777778877887 
		name=Borys 
		lastname=Petrenko`,
	},
	{
		Name:      "info",
		Usage:     "Show information about user stored in database",
		ArgsUsage: "login=mylogin",
		Action:    InfoAction,
		Description: `Show  account info by user name (login)
		login=mylogin `,
	},
}
