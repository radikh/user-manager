// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/logger"
	"github.com/lvl484/user-manager/model"
	"github.com/lvl484/user-manager/storage"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

// ActionChecker auxiliary interface for commands and testing
type ActionChecker interface {
	Config() (*config.Config, error)
	DBConfig(context.Context) (*storage.DBConfig, error)
	ConnectToDB(*storage.DBConfig) (*sql.DB, error)
	UsersRepo() (*model.UsersRepo, error)
	MessageCommandDone(msg string, err error) error
	ExecuteAction(c *cli.Context, action int) error
}

// actionHandle structure that implements ActionChecker interface
type actionHandle struct {
	ccfg *config.Config
}

// NewConfig function that replace config.NewConfig()
func (ah *actionHandle) Config() (*config.Config, error) {
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
func (ah *actionHandle) UsersRepo() (*model.UsersRepo, error) {
	cfg, err := ah.Config()
	if err != nil {
		return nil, errors.Wrap(err, msgErrorReadConfig)
	}
	ah = ah.SetConfiguration(cfg)
	dbcfg, err := ah.DBConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, msgErrorDBConfig)
	}
	loggerConfig, err := cfg.LoggerConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, msgErrorDBConfig)
	}
	err = logger.SetLogger(loggerConfig)
	if err != nil {
		return nil, errors.Wrap(err, msgErrorDBConfig)
	}
	db, err := ah.ConnectToDB(dbcfg)
	if err != nil {
		return nil, errors.Wrap(err, msgErrorConnectDB)
	}
	repo := model.SetUsersRepo(db)
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
func (ah *actionHandle) MessageCommandDone(msg string, err error) error {
	if err != nil {
		_, err = os.Stdout.WriteString(fmt.Sprintf("%s exit with an error: %s\n", msg, err.Error()))
	} else {
		_, err = os.Stdout.WriteString(msg)
	}
	return err
}

// logAction log messages about command execution
func (ah *actionHandle) logAction(msg string, err error) {
	if err != nil {
		msg += fmt.Sprintf(" caused an error: %s\n", err.Error())
		logger.LogUM.Error(msg)
	} else {
		msg += fmt.Sprintf("%s\n", " was successful")
		logger.LogUM.Info(msg)
	}
}

// SetConfiguration set configuration to handle that manipulate commands
func (ah *actionHandle) SetConfiguration(cfg *config.Config) *actionHandle {
	return &actionHandle{ccfg: cfg}
}

// ExecuteAction execute prepare and     command
func (ah *actionHandle) ExecuteAction(c *cli.Context, action int) error {
	var argumentValue interface{}
	returnMessage := msgErrorActionInput
	logMessage := msgErrorActionInput
	err := ah.checkRole()
	if err != nil {
		return err
	}
	repo, err := ah.UsersRepo()
	if err != nil {
		return err
	}
	if err != nil {
		return errors.Wrap(err, msgErrorConnectDB)
	}

	if action > 1 {
		argumentValue, err = ah.splitLogin(c)
	} else {
		argumentValue, err = ah.createUser(c)
	}
	if err != nil {
		return errors.Wrap(err, msgConverteUser)
	}
	switch action {
	case actionCreate:
		logMessage = msgCreate
		user := argumentValue.(*model.User)
		err = repo.Add(user)
		returnMessage = fmt.Sprintf("%s was created", user.Username)
	case actionUpdate:
		logMessage = msgUpdate
		user := argumentValue.(*model.User)
		err = repo.Update(user)
		returnMessage = fmt.Sprintf("%s was updated", user.Username)
	case actionInfo:
		var user *model.User
		logMessage = msgGetInfo
		user, err = repo.GetInfo(argumentValue.(string))
		user.Password = "*******"
		returnMessage = fmt.Sprintf("%+v", user)
	case actionDelete:
		logMessage = msgDelete
		returnMessage = msgUserDeleted
		err = repo.Delete(argumentValue.(string))
	case actionActivate:
		logMessage = msgActivate
		returnMessage = msgUserActivated
		err = repo.Activate(argumentValue.(string))
	case actionDisable:
		logMessage = msgDisable
		returnMessage = msgUserDisable
		err = repo.Disable(argumentValue.(string))
	default:
		err = errors.New(msgErrorActionInput)
	}
	ah.logAction(fmt.Sprintf(msgFormat, logMessage, returnMessage), err)
	return actionHelper.MessageCommandDone(returnMessage, err)
}

// ExecuteAction execute prepare and     command
func (ah *actionHandle) checkRole() error {
	username, password, err := ah.getCredentials()
	if err != nil {
		return errors.Wrap(err, msgErrorCheckCredentials)
	}
	repo, err := ah.UsersRepo()
	if err != nil {
		return errors.Wrap(err, msgErrorConnectDB)
	}
	status, err := repo.CheckAdminRole(username, password)
	if err != nil {
		return err
	}
	if !status {
		return errors.Wrap(err, msgErrorCheckCredentials)
	}
	return nil
}

// ExecuteAction execute prepare and     command
func (ah *actionHandle) getCredentials() (login string, pwd string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your's credentials")
	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	password := string(bytePassword)

	return username, password, err
}
