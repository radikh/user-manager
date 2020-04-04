// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package main

import (
	"flag"
	"io/ioutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lvl484/user-manager/model"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestReturnRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ur := model.NewUsersRepo(db)

	assert.NotNil(t, ur)
}

func TestSsplitParam(t *testing.T) {
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
			gotPName, gotPValue, err := splitParam(tt.args.param)
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
	err := appendParam(&user, "email=boss2@company.com")
	assert.NoError(t, err)
	err = appendParam(&user, "lastname=Porotrenko")
	assert.NoError(t, err)
	assert.Equal(t, &user1, &user)
	err = appendParam(&user, "lastnamePorotrenko")
	assert.Error(t, err)
}

func TestCreateUser(t *testing.T) {
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

	user1, err := createUser(context)
	assert.NoError(t, err)
	assert.Equal(t, &user, user1)

}

func TestSplitLogin(t *testing.T) {
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
			gotPValue, err := splitLogin(context)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.wantPValue, gotPValue)
		})
	}
}

func Test_messageCommandDone(t *testing.T) {
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
			err := messageCommandDone(tt.args.msg)
			assert.NoError(t, err)
		})
	}
}

func TestCreateAction(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("CreateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInfoAction(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InfoAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InfoAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActivateAction(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ActivateAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ActivateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDisableAction(t *testing.T) {

}

func TestUpdateAction(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteAction(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_splitLogin(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitLogin(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("splitLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}
