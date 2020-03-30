// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/model"
	"github.com/lvl484/user-manager/storage"
)

var ErrUserConnvertParam = errors.New("There is no such field to assigned value")

// umcliCommands commands of CUI interface
var umcliCommands = cli.Commands{
	{
		Name:      "create",
		ArgsUsage: "login=mylogin pwd=ttrtyrghgfh email=boss@company.com phone=7778777778877887 name=Borys lastname=Petrenko ",
		Action:    CreateAction,
	},
	{
		Name:      "delete",
		ArgsUsage: "login=mylogin",
		Action:    DeleteAction,
	},
	{
		Name:      "disable",
		ArgsUsage: "login=mylogin",
		Action:    DisableAction,
	},
	{
		Name:      "activate",
		ArgsUsage: "login=mylogin",
		Action:    ActivateAction,
	},
	{
		Name:      "update",
		ArgsUsage: "login=mylogin pwd=ttrtyrghgfh email=boss@company.com phone=7778777778877887 name=Borys lastname=Petrenko ",
		Action:    UpdateAction,
	},
	{
		Name:      "info",
		ArgsUsage: "login=mylogin",
		Action:    InfoAction,
	},
}

// GetCommands return commands of CUI
func GetCommands() cli.Commands {
	return umcliCommands
}

// returnRepo return the repo that holds database
func returnRepo(c *cli.Context) (*usersRepo, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	dbcfg, err := cfg.DBConfig(context.Background())
	if err != nil {
		return nil, err
	}
	db, err := storage.ConnectToDB(dbcfg)
	if err != nil {
		return nil, err
	}
	repo := NewUsersRepo(db)
	return repo, nil
}

// splitParam split input argument into field name and value
func splitParam(param string) (pName string, pValue string) {
	params := strings.Split(param, "=")
	pName = params[0]
	pValue = params[1]
	return pName, pValue
}

// appendParam assign value to User structure field by name
func appendParam(user *model.User, param string) error {
	pName, pValue := splitParam(param)
	switch pName {
	case "login":
		user.Username = pValue
	case "pwd":
		user.Password = pValue
	case "email":
		user.Email = pValue
	case "phone":
		user.Phone = pValue
	case "name":
		user.FirstName = pValue
	case "lastname":
		user.LastName = pValue
	default:
		return ErrUserConnvertParam
	}
	return nil
}

// CreateAction create new user in database
func CreateAction(c *cli.Context) error {
	var user model.User
	var err error
	params := c.Args().Slice()
	for _, param := range params {
		err = appendParam(&user, param)
		if err != nil {
			return err
		}
	}
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Add(&user)
	return err
}

// InfoAction get info of user by its login name
func InfoAction(c *cli.Context) error {
	var user *model.User
	param := c.Args().Get(0)
	_, pValue := splitParam(param)
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	user, err = repo.GetInfo(pValue)
	if err != nil {
		return err
	}
	err = json.NewEncoder(os.Stdout).Encode(user)
	if err != nil {
		return err
	}
	return nil
}

// ActivateAction activate user that was disabled
func ActivateAction(c *cli.Context) error {
	param := c.Args().Get(0)
	_, pValue := splitParam(param)
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Activate(pValue)
	if err != nil {
		return err
	}
	_, err = os.Stdout.WriteString("User activated successfully")
	if err != nil {
		return err
	}
	return nil
}

// DisableAction disable user by its login
func DisableAction(c *cli.Context) error {
	param := c.Args().Get(0)
	_, pValue := splitParam(param)
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Disable(pValue)
	if err != nil {
		return err
	}
	_, err = os.Stdout.WriteString("User disabled successfully")
	if err != nil {
		return err
	}
	return nil
}

// UpdateAction update information about user
func UpdateAction(c *cli.Context) error {
	var user model.User
	var err error
	params := c.Args().Slice()
	for _, param := range params {
		err = appendParam(&user, param)
		if err != nil {
			return err
		}
	}
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Update(&user)
	if err != nil {
		return err
	}
	_, err = os.Stdout.WriteString("User update successfully")
	if err != nil {
		return err
	}
	err = json.NewEncoder(os.Stdout).Encode(user)
	if err != nil {
		return err
	}
	return err
}

// DeleteAction delete user in database
func DeleteAction(c *cli.Context) error {
	param := c.Args().Get(0)
	_, pValue := splitParam(param)
	repo, err := returnRepo(c)
	if err != nil {
		return err
	}
	err = repo.Delete(pValue)
	if err != nil {
		return err
	}
	_, err = os.Stdout.WriteString("User deleted successfully")
	if err != nil {
		return err
	}
	return nil
}
