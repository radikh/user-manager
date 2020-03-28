// Command umcli provides admin command line tool to manipulate accounts with admin rights.
package pgclient

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/lvl484/user-manager/model"
)

func TestNewUsersRepo(t *testing.T) {
	type args struct {
		data *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *usersRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUsersRepo(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUsersRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_usersRepo_Add(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			if err := ur.Add(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_usersRepo_Update(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			if err := ur.Update(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_usersRepo_Delete(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			if err := ur.Delete(tt.args.login); (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_usersRepo_Disable(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			if err := ur.Disable(tt.args.login); (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.Disable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_usersRepo_Activate(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			if err := ur.Activate(tt.args.login); (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.Activate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_usersRepo_GetInfo(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			got, err := ur.GetInfo(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.GetInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("usersRepo.GetInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_usersRepo_getUserDiactived(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			got, err := ur.getUserDiactived(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.getUserDiactived() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("usersRepo.getUserDiactived() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_usersRepo_CheckLoginExist(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := &usersRepo{
				db: tt.fields.db,
			}
			got, err := ur.CheckLoginExist(tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("usersRepo.CheckLoginExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("usersRepo.CheckLoginExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
