// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"context"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"

	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/mock"
	"github.com/lvl484/user-manager/model"
	"github.com/lvl484/user-manager/storage"
)

func TestReturnRepo(t *testing.T) {
	var err error
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := config.Config{
		PostgresUser: "postgres",
		PostgresPass: "1q2w3e4r",
		PostgresDB:   "um_db",
	}
	dbconf := &storage.DBConfig{
		Host:     "db",
		Port:     5432,
		User:     "user",
		Password: "password",
		DBName:   "dbname",
	}
	mockConf := mock.NewMockActionChecker(mockCtrl)
	mockConf.EXPECT().NewConfig().Return(&c, nil)
	conf, err := mockConf.NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, &c, conf)

	mockConf.EXPECT().DBConfig(context.Background()).Return(dbconf, nil)
	dconf, err := mockConf.DBConfig(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, dbconf, dconf)

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo1 := model.NewUsersRepo(db)
	mockConf.EXPECT().ConnectToDB(dbconf).Return(db, nil)
	ddb, err := mockConf.ConnectToDB(dbconf)
	assert.NoError(t, err)
	assert.NotNil(t, ddb)

	mockConf.EXPECT().ReturnRepo().Return(repo1, nil)
	repo, err := mockConf.ReturnRepo()
	assert.NoError(t, err)
	assert.NotNil(t, repo)

}

func TestSsplitParam(t *testing.T) {
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

func TestAppendParam(t *testing.T) {
	var ah actionHandle
	user := model.User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "password",
		Email:     "email1@company.com",
		FirstName: "Pedro",
		LastName:  "Petrenko",
		Phone:     "77777777777",
	}
	user1 := model.User{
		ID:        "3b60ac82-5e8f-4010-ac99-2344cfa72ce0",
		Username:  "user1",
		Password:  "password",
		Email:     "boss2@company.com",
		FirstName: "Pedro",
		LastName:  "Porotrenko",
		Phone:     "77777777777",
	}
	err := ah.appendParam(&user, "email=boss2@company.com")
	assert.NoError(t, err)
	err = ah.appendParam(&user, "lastname=Porotrenko")
	assert.NoError(t, err)
	assert.Equal(t, &user1, &user)
	err = ah.appendParam(&user, "lastnamePorotrenko")
	assert.Error(t, err)
}

func TestCreateUser(t *testing.T) {
	var ah actionHandle
	user := model.User{
		Username:  "user1",
		Password:  "password",
		Email:     "boss@company.com",
		FirstName: "Petro",
		LastName:  "Petrenko",
		Phone:     "7778777778877887",
	}
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1", "pwd=password", "email=boss@company.com", "phone=7778777778877887", "name=Petro", "lastname=Petrenko"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.createUser(context)
	assert.Equal(t, &user, user1)
	assert.NoError(t, err)
}

func TestSplitLogin(t *testing.T) {
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

func TestMessageCommandDone(t *testing.T) {
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
			err := ah.MessageCommandDone(tt.args.msg)
			assert.NoError(t, err)
		})
	}
}

func TestCreateAction(t *testing.T) {
	var ah actionHandle
	user := model.User{
		Username:  "user1",
		Password:  "password",
		Email:     "boss@company.com",
		FirstName: "Petro",
		LastName:  "Petrenko",
		Phone:     "7778777778877887",
	}
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1", "pwd=password", "email=boss@company.com", "phone=7778777778877887", "name=Petro", "lastname=Petrenko"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.createUser(context)
	assert.Equal(t, &user, user1)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConf := mock.NewMockActionChecker(mockCtrl)
	repo1 := &model.UsersRepo{}
	mockConf.EXPECT().ReturnRepo().Return(repo1, nil)
	repo, err := mockConf.ReturnRepo()

	assert.NoError(t, err)
	assert.Equal(t, repo1, repo)

	modelMock := mock.NewMockUsers(mockCtrl)
	modelMock.EXPECT().Add(user).Return(nil)
	err = modelMock.Add(user)
	assert.NoError(t, err)
}

func TestInfoAction(t *testing.T) {
	userTest := model.User{
		Username:  "user1",
		Password:  "password",
		Email:     "boss@company.com",
		FirstName: "Petro",
		LastName:  "Petrenko",
		Phone:     "7778777778877887",
	}
	var ah actionHandle
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.splitLogin(context)
	assert.Equal(t, "user1", user1)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConf := mock.NewMockActionChecker(mockCtrl)
	repo1 := &model.UsersRepo{}
	mockConf.EXPECT().ReturnRepo().Return(repo1, nil)
	repo, err := mockConf.ReturnRepo()

	assert.NoError(t, err)
	assert.Equal(t, repo1, repo)

	modelMock := mock.NewMockUsers(mockCtrl)
	modelMock.EXPECT().GetInfo(user1).Return(&userTest, nil)
	users, err := modelMock.GetInfo(user1)
	assert.NoError(t, err)
	assert.NotNil(t, users)
}

func TestActivateAction(t *testing.T) {
	var ah actionHandle
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.splitLogin(context)
	assert.Equal(t, "user1", user1)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConf := mock.NewMockActionChecker(mockCtrl)
	repo1 := &model.UsersRepo{}
	mockConf.EXPECT().ReturnRepo().Return(repo1, nil)
	repo, err := mockConf.ReturnRepo()

	assert.NoError(t, err)
	assert.Equal(t, repo1, repo)

	modelMock := mock.NewMockUsers(mockCtrl)
	modelMock.EXPECT().Activate(user1).Return(nil)
	err = modelMock.Activate(user1)
	assert.NoError(t, err)
}

func TestDisableAction(t *testing.T) {
	var ah actionHandle
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.splitLogin(context)
	assert.Equal(t, "user1", user1)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConf := mock.NewMockActionChecker(mockCtrl)
	repo1 := &model.UsersRepo{}
	mockConf.EXPECT().ReturnRepo().Return(repo1, nil)
	repo, err := mockConf.ReturnRepo()

	assert.NoError(t, err)
	assert.Equal(t, repo1, repo)

	modelMock := mock.NewMockUsers(mockCtrl)
	modelMock.EXPECT().Disable(user1).Return(nil)
	err = modelMock.Disable(user1)
	assert.NoError(t, err)
}

func TestDeleteAction(t *testing.T) {
	var ah actionHandle
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.splitLogin(context)
	assert.Equal(t, "user1", user1)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConf := mock.NewMockActionChecker(mockCtrl)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo1 := model.NewUsersRepo(db)
	mockConf.EXPECT().ReturnRepo().Return(repo1, nil)
	repo, err := mockConf.ReturnRepo()

	assert.NoError(t, err)
	assert.Equal(t, repo1, repo)

	modelMock := mock.NewMockUsers(mockCtrl)
	modelMock.EXPECT().Delete(user1).Return(nil)
	err = modelMock.Delete(user1)
	assert.NoError(t, err)
}

func TestUpdateAction(t *testing.T) {
	userTest := model.User{
		Username:  "user1",
		Password:  "password",
		Email:     "boss@company.com",
		FirstName: "Petro",
		LastName:  "Petrenko",
		Phone:     "7778777778877887",
	}
	var ah actionHandle
	app := &cli.App{Writer: ioutil.Discard}
	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{"login=user1"})

	context := cli.NewContext(app, set, nil)

	user1, err := ah.splitLogin(context)
	assert.Equal(t, "user1", user1)
	assert.NoError(t, err)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConf := mock.NewMockActionChecker(mockCtrl)
	repo1 := &model.UsersRepo{}
	mockConf.EXPECT().ReturnRepo().Return(repo1, nil)
	repo, err := mockConf.ReturnRepo()

	assert.NoError(t, err)
	assert.Equal(t, repo1, repo)

	modelMock := mock.NewMockUsers(mockCtrl)
	modelMock.EXPECT().Update(&userTest).Return(nil)
	err = modelMock.Update(&userTest)
	assert.NoError(t, err)
}
