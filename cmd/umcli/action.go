// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/lvl484/user-manager/model"
)

var actionHelper *actionHandle

// CreateAction create new user in database
func CreateAction(c *cli.Context) error {
	user, err := actionHelper.createUser(c)
	if err != nil {
		return errors.Wrap(err, msgConverteUser)
	}
	repo, err := actionHelper.ReturnRepo()
	if err != nil {
		return err
	}
	err = repo.Add(user)
	return err
}

// InfoAction get info of user by its login name
func InfoAction(c *cli.Context) error {
	var user *model.User

	pValue, err := actionHelper.splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := actionHelper.ReturnRepo()
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
	pValue, err := actionHelper.splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := actionHelper.ReturnRepo()
	if err != nil {
		return err
	}
	err = repo.Activate(pValue)
	if err != nil {
		return err
	}

	return actionHelper.MessageCommandDone(msgUserActivated)
}

// DisableAction disable user by its login
func DisableAction(c *cli.Context) error {
	pValue, err := actionHelper.splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := actionHelper.ReturnRepo()
	if err != nil {
		return err
	}
	err = repo.Disable(pValue)
	if err != nil {
		return err
	}

	return actionHelper.MessageCommandDone(msgUserDisable)
}

// UpdateAction update information about user
func UpdateAction(c *cli.Context) error {
	user, err := actionHelper.createUser(c)
	if err != nil {
		return errors.Wrap(err, msgConverteUser)
	}

	repo, err := actionHelper.ReturnRepo()
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
	pValue, err := actionHelper.splitLogin(c)
	if err != nil {
		return err
	}
	repo, err := actionHelper.ReturnRepo()
	if err != nil {
		return err
	}
	err = repo.Delete(pValue)
	if err != nil {
		return err
	}

	return actionHelper.MessageCommandDone(msgUserDeleted)
}
