// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/model"
	"github.com/lvl484/user-manager/storage"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func mockDB(t *testing.T) *sql.DB {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db
}

func mockUser() *model.User {
	user := model.User{
		Username:  "user1",
		Password:  "password",
		Email:     "boss@company.com",
		FirstName: "Petro",
		LastName:  "Porotrenko",
		Phone:     "7778777778877887",
	}
	return &user
}

func TestActionHandleConfig(t *testing.T) {
	c := &config.Config{
		PostgresUser: "postgres",
		PostgresPass: "1q2w3e4r",
		PostgresDB:   "um_db",
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cliMock := NewMockActionChecker(mockCtrl)
	cliMock.EXPECT().Config().Return(c, nil)
	conf, err := cliMock.Config()
	assert.Equal(t, c, conf)
	assert.NoError(t, err)

}

func TestActionHandleDBConfig(t *testing.T) {
	dbconf := &storage.DBConfig{
		Host:     "db",
		Port:     5432,
		User:     "user",
		Password: "password",
		DBName:   "dbname",
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cliMock := NewMockActionChecker(mockCtrl)
	cliMock.EXPECT().DBConfig(context.Background()).Return(dbconf, nil)
	dconf, err := cliMock.DBConfig(context.Background())
	assert.Equal(t, dbconf, dconf)
	assert.NoError(t, err)
}

func TestActionHandleConnectToDB(t *testing.T) {
	dbconf := &storage.DBConfig{
		Host:     "db",
		Port:     5432,
		User:     "user",
		Password: "password",
		DBName:   "dbname",
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cliMock := NewMockActionChecker(mockCtrl)
	db := mockDB(t)
	cliMock.EXPECT().ConnectToDB(dbconf).Return(db, nil)
	ddb, err := cliMock.ConnectToDB(dbconf)
	assert.NoError(t, err)
	assert.NotNil(t, ddb)

}

func TestActionHandleUsersRepo(t *testing.T) {
	var err error
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cliMock := NewMockActionChecker(mockCtrl)
	db := mockDB(t)
	repo1 := model.SetUsersRepo(db)
	cliMock.EXPECT().UsersRepo().Return(repo1, nil)
	repo, err := cliMock.UsersRepo()
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestActionHandleSplitParam(t *testing.T) {
	var ah actionHandle
	type args struct {
		param string
	}
	tests := []struct {
		name       string
		args       args
		wantPName  string
		wantPValue string
		wantErr    bool
	}{
		{name: "test1", args: args{"name=Petro"}, wantPName: "name", wantPValue: "Petro", wantErr: false},
		{name: "test1", args: args{"login=mylogin"}, wantPName: "login", wantPValue: "mylogin", wantErr: false},
		{name: "test1", args: args{"pwd=password"}, wantPName: "pwd", wantPValue: "password", wantErr: false},
		{name: "test1", args: args{"email=boss2@company.com"}, wantPName: "email", wantPValue: "boss2@company.com", wantErr: false},
		{name: "test1", args: args{"lastname=Petrenko"}, wantPName: "lastname", wantPValue: "Petrenko", wantErr: false},
		{name: "test1", args: args{"emalboss2@company.com"}, wantPName: "", wantPValue: "", wantErr: true},
		{name: "test1", args: args{"lastmePetrenko"}, wantPName: "", wantPValue: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPName, gotPValue, err := ah.splitParam(tt.args.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantPName, gotPName)
			assert.Equal(t, tt.wantPValue, gotPValue)
		})
	}
}

func TestActionHandleAppendParam(t *testing.T) {
	var ah actionHandle
	user := model.User{
		Username:  "user1",
		Password:  "password",
		Email:     "email1@company.com",
		FirstName: "Petro",
		LastName:  "Petrenko",
		Phone:     "7778777778877887",
	}
	user1 := mockUser()

	err := ah.appendParam(&user, "email=boss@company.com")
	assert.NoError(t, err)
	err = ah.appendParam(&user, "lastname=Porotrenko")
	assert.NoError(t, err)
	assert.Equal(t, user1, &user)
	err = ah.appendParam(&user, "lastnamePorotrenko")
	assert.Error(t, err)
}

func TestActionHandleCreateUser(t *testing.T) {
	var ah actionHandle
	user := mockUser()
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1", "pwd=password", "email=boss@company.com", "phone=7778777778877887", "name=Petro", "lastname=Porotrenko"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.createUser(context)
	assert.Equal(t, user, user1)
	assert.NoError(t, err)
}

func TestActionHandleSplitLogin(t *testing.T) {
	var ah actionHandle
	type args struct {
		msg []string
	}
	tests := []struct {
		name       string
		args       args
		wantPValue string
		wantErr    bool
	}{
		{name: "test1", args: args{[]string{"login=Petro"}}, wantPValue: "Petro", wantErr: false},
		{name: "test1", args: args{[]string{"login=mylogin"}}, wantPValue: "mylogin", wantErr: false},
		{name: "test1", args: args{[]string{"pwd=password"}}, wantPValue: "password", wantErr: true},
		{name: "test1", args: args{[]string{"emalboss2@company.com"}}, wantPValue: "", wantErr: true},
		{name: "test1", args: args{[]string{"lastmePetrenko"}}, wantPValue: "", wantErr: true},
	}

	for _, tt := range tests {
		app := &cli.App{Writer: ioutil.Discard}
		set := flag.NewFlagSet("test", 0)
		_ = set.Parse(tt.args.msg)

		context := cli.NewContext(app, set, nil)
		t.Run(tt.name, func(t *testing.T) {
			gotPValue, err := ah.splitLogin(context)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantPValue, gotPValue)
		})
	}
}

func TestActionHandleMessageCommandDone(t *testing.T) {
	var ah actionHandle
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"Some message1"}, false},
		{"test2", args{"Some message2"}, false},
		{"test3", args{"Some message3"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ah.MessageCommandDone(tt.args.msg, nil)
			assert.NoError(t, err)
		})
	}
}

func TestActionHandleSetConfiguration(t *testing.T) {
	c := &config.Config{
		PostgresUser: "postgres",
		PostgresPass: "1q2w3e4r",
		PostgresDB:   "um_db",
	}
	ah := &actionHandle{ccfg: c}
	var ahTest *actionHandle
	ahTest = ahTest.SetConfiguration(c)
	assert.Equal(t, ah, ahTest)

}

func TestActionHandleExecuteAction(t *testing.T) {
	var err error
	var ah actionHandle
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1", "pwd=password", "email=boss@company.com", "phone=7778777778877887", "name=Petro", "lastname=Porotrenko"})
	context := cli.NewContext(app, set, nil)
	err = ah.ExecuteAction(context, 7)
	assert.Error(t, err)
}

func TestGetCredentials(t *testing.T) {

	input := "Test\n"
	reader := bufio.NewReader(strings.NewReader(input))

	value, err := reader.ReadBytes(byte(10))
	assert.NoError(t, err)
	assert.Equal(t, input, string(value))
}
