// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"flag"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestCreateAction(t *testing.T) {
	var err error
	context := mockContext()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cliMock := NewMockActionChecker(mockCtrl)
	cliMock.EXPECT().ExecuteAction(context, actionCreate).Return(nil)
	err = cliMock.ExecuteAction(context, actionCreate)
	assert.NoError(t, err)
}

func TestInfoAction(t *testing.T) {
	var err error
	context := mockContext()
	cliMock := mockAction(t)
	cliMock.EXPECT().ExecuteAction(context, actionInfo).Return(nil)
	err = cliMock.ExecuteAction(context, actionInfo)
	assert.NoError(t, err)
}

func TestActivateAction(t *testing.T) {
	var err error
	context := mockContext()
	cliMock := mockAction(t)
	cliMock.EXPECT().ExecuteAction(context, actionActivate).Return(nil)
	err = cliMock.ExecuteAction(context, actionActivate)
	assert.NoError(t, err)
}

func TestDisableAction(t *testing.T) {
	var err error
	context := mockContext()
	cliMock := mockAction(t)
	cliMock.EXPECT().ExecuteAction(context, actionDisable).Return(nil)
	err = cliMock.ExecuteAction(context, actionDisable)
	assert.NoError(t, err)
}

func TestUpdateAction(t *testing.T) {
	var err error
	context := mockContext()
	cliMock := mockAction(t)
	cliMock.EXPECT().ExecuteAction(context, actionUpdate).Return(nil)
	err = cliMock.ExecuteAction(context, actionUpdate)
	assert.NoError(t, err)
}

func TestDeleteAction(t *testing.T) {
	var err error
	context := mockContext()
	cliMock := mockAction(t)
	cliMock.EXPECT().ExecuteAction(context, actionDelete).Return(nil)
	err = cliMock.ExecuteAction(context, actionDelete)
	assert.NoError(t, err)
}

func mockContext() *cli.Context {
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1", "pwd=password", "email=boss@company.com", "phone=7778777778877887", "name=Petro", "lastname=Porotrenko"})
	context := cli.NewContext(app, set, nil)
	return context
}

func mockAction(t *testing.T) *MockActionChecker {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	actionMock := NewMockActionChecker(mockCtrl)
	return actionMock
}
