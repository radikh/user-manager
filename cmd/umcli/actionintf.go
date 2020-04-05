// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/model"
	"github.com/lvl484/user-manager/storage"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// ActionChecker auxiliary interface for commands and testing
type ActionChecker interface {
	NewConfig() (*config.Config, error)
	DBConfig(context.Context) (*storage.DBConfig, error)
	ConnectToDB(*storage.DBConfig) (*sql.DB, error)
	ReturnRepo() (*model.UsersRepo, error)
	MessageCommandDone(msg string) error
}

// actionHandle structure that implements ActionChecker interface
type actionHandle struct {
	ccfg *config.Config
}

// NewConfig function that replace config.NewConfig()
func (ah *actionHandle) NewConfig() (*config.Config, error) {
	return config.NewConfig()
}

// DBConfig function that replace config.Config.DBConfig()
func (ah *actionHandle) DBConfig(ctx context.Context) (*storage.DBConfig, error) {
	return ah.ccfg.DBConfig(ctx)
}

// ConnectToDB function that replace storage.ConnectToDB
func (ah *actionHandle) ConnectToDB(dbconf *storage.DBConfig) (*sql.DB, error) {
	return storage.ConnectToDB(dbconf)
}

// returnRepo return the repo that holds database
func (ah *actionHandle) ReturnRepo() (*model.UsersRepo, error) {
	cfg, err := ah.NewConfig()
	if err != nil {
		return nil, errors.Wrap(err, msgErrorReadConfig)
	}
	ah.ccfg = cfg
	dbcfg, err := ah.DBConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, msgErrorDBConfig)
	}
	db, err := ah.ConnectToDB(dbcfg)
	if err != nil {
		return nil, errors.Wrap(err, msgErrorConnectDB)
	}
	repo := model.NewUsersRepo(db)
	return repo, nil
}

// splitParam split input argument into field name and value
func (ah *actionHandle) splitParam(param string) (pName string, pValue string, err error) {
	params := strings.Split(param, splitDelimiter)
	if len(params) < 2 {
		return "", "", errors.New(msgWrongInputArguments)
	}
	pName = params[0]
	pValue = params[1]
	return pName, pValue, nil
}

// appendParam assign value to User structure field by name
func (ah *actionHandle) appendParam(user *model.User, param string) error {

	pName, pValue, err := ah.splitParam(param)
	if err != nil {
		return err
	}
	pField := reflect.ValueOf(user).Elem().FieldByName(convertFieldName()[pName])

	if !pField.CanSet() {
		return errors.New(fmt.Sprintf("%s : %s\n", msgConnvertParam, pName))
	}
	reflect.ValueOf(user).Elem().FieldByName(convertFieldName()[pName]).SetString(pValue)

	return nil
}

// createUser draws draws out fields from arguments to User structure
func (ah *actionHandle) createUser(c *cli.Context) (*model.User, error) {
	var user model.User
	var err error
	params := c.Args().Slice()
	if len(params) == 0 {
		return nil, errors.New(msgEmptyInputArguments)
	}
	for _, param := range params {
		err = ah.appendParam(&user, param)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

// splitLogin draws out user name from input argument
func (ah *actionHandle) splitLogin(c *cli.Context) (string, error) {
	if c.Args().Len() == 0 {
		return "", errors.New(msgEmptyInputArguments)
	}
	param := c.Args().First()
	pName, pValue, err := ah.splitParam(param)
	if err != nil {
		return "", err
	}
	if pName != argLogin {
		return "", errors.New(msgWrongInputArguments)
	}
	return pValue, nil
}

// messageWorkDone messages that command is done
func (ah *actionHandle) MessageCommandDone(msg string) error {
	_, err := os.Stdout.WriteString(msg)
	return err
}
