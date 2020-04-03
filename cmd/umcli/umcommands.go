// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/model"
	"github.com/lvl484/user-manager/storage"
)

const (
	splitDelimiter = "="

	msgConnvertParam       = "There is no such field to assigned value"
	msgConverteUser        = "Converting to structure User failed"
	msgErrorReadConfig     = "Reading configuration failed"
	msgErrorDBConfig       = "Creating database configuration failed"
	msgErrorConnectDB      = "Connecting to database failed"
	msgEmptyInputArguments = "No arguments for command line"
	msgUserActivated       = "User activated successfully"
	msgUserDeleted         = "User deleted successfully"
	msgUserDisable         = "User disabled successfully"
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
		ArgsUsage: "login=mylogin",
		Action:    DeleteAction,
		Description: `Delete account by user name (login)
		login=mylogin `,
	},
	{
		Name:      "disable",
		ArgsUsage: "login=mylogin",
		Action:    DisableAction,
		Description: `Disable account by user name (login), without deleting it
		login=mylogin `,
	},
	{
		Name:      "activate",
		ArgsUsage: "login=mylogin",
		Action:    ActivateAction,
		Description: `Activate account by user name (login), that was disabled
		login=mylogin `,
	},
	{
		Name:      "update",
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
		ArgsUsage: "login=mylogin",
		Action:    InfoAction,
		Description: `Show  account info by user name (login)
		login=mylogin `,
	},
}

// returnRepo return the repo that holds database
func returnRepo(c *cli.Context) (*model.UsersRepo, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, msgErrorReadConfig)
	}
	dbcfg, err := cfg.DBConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, msgErrorDBConfig)
	}
	db, err := storage.ConnectToDB(dbcfg)
	if err != nil {
		return nil, errors.Wrap(err, msgErrorConnectDB)
	}
	repo := model.NewUsersRepo(db)
	return repo, nil
}

// splitParam split input argument into field name and value
func splitParam(param string) (pName string, pValue string) {
	params := strings.Split(param, splitDelimiter)
	pName = params[0]
	pValue = params[1]
	return pName, pValue
}

// appendParam assign value to User structure field by name
func appendParam(user *model.User, param string) error {

	pName, pValue := splitParam(param)
	pField := reflect.ValueOf(user).Elem().FieldByName(convertFieldName()[pName])

	if !pField.CanSet() {
		return errors.New(msgConnvertParam)
	}
	reflect.ValueOf(user).Elem().FieldByName(convertFieldName()[pName]).SetString(pValue)

	return nil
}

// createUser draws draws out fields from arguments to User structure
func createUser(c *cli.Context) (*model.User, error) {
	var user model.User
	var err error
	params := c.Args().Slice()
	if len(params) == 0 {
		return nil, errors.New(msgEmptyInputArguments)
	}
	for _, param := range params {
		err = appendParam(&user, param)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

// splitLogin draws out user name from input argument
func splitLogin(c *cli.Context) (string, error) {
	if c.Args().Len() == 0 {
		return "", errors.New(msgEmptyInputArguments)
	}
	param := c.Args().Get(0)
	_, pValue := splitParam(param)
	return pValue, nil
}

// messageWorkDone messages that command is done
func messageCommandDone(msg string) error {
	_, err := os.Stdout.WriteString(msg)
	return err
}

// CreateAction create new user in database
func CreateAction(c *cli.Context) error {
	user, err := createUser(c)
	if err != nil {
		return errors.Wrap(err, msgConverteUser)
	}
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Add(user)
	return err
}

// InfoAction get info of user by its login name
func InfoAction(c *cli.Context) error {
	var user *model.User

	pValue, err := splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	user, err = repo.GetInfo(pValue)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(user)
}

// ActivateAction activate user that was disabled
func ActivateAction(c *cli.Context) error {
	pValue, err := splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Activate(pValue)
	if err != nil {
		return err
	}

	return messageCommandDone(msgUserActivated)
}

// DisableAction disable user by its login
func DisableAction(c *cli.Context) error {
	pValue, err := splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Disable(pValue)
	if err != nil {
		return err
	}

	return messageCommandDone(msgUserDisable)
}

// UpdateAction update information about user
func UpdateAction(c *cli.Context) error {
	user, err := createUser(c)
	if err != nil {
		return errors.Wrap(err, msgConverteUser)
	}

	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Update(user)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(user)
}

// DeleteAction delete user in database
func DeleteAction(c *cli.Context) error {
	pValue, err := splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Delete(pValue)
	if err != nil {
		return err
	}

	return messageCommandDone(msgUserDeleted)
}
