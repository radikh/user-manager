// Package http presents HTTP server implementation.
// It provides a REST API to perform a set of CRUD
// to manage users and an endpoint to authenticate.
package server

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lvl484/user-manager/config"
	"github.com/lvl484/user-manager/model"
)

func TestNewHTTP(t *testing.T) {
	type args struct {
		cfg *config.Config
		ur  *model.UsersRepo
	}
	tests := []struct {
		name string
		args args
		want *HTTP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTP(tt.args.cfg, tt.args.ur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTP_newRouter(t *testing.T) {
	type fields struct {
		srv *http.Server
		acc *account
		ur  *model.UsersRepo
	}
	tests := []struct {
		name   string
		fields fields
		want   *mux.Router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTP{
				srv: tt.fields.srv,
				acc: tt.fields.acc,
				ur:  tt.fields.ur,
			}
			if got := h.newRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTP.newRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTP_Start(t *testing.T) {
	type fields struct {
		srv *http.Server
		acc *account
		ur  *model.UsersRepo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTP{
				srv: tt.fields.srv,
				acc: tt.fields.acc,
				ur:  tt.fields.ur,
			}
			if err := h.Start(); (err != nil) != tt.wantErr {
				t.Errorf("HTTP.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHTTP_Stop(t *testing.T) {
	type fields struct {
		srv *http.Server
		acc *account
		ur  *model.UsersRepo
	}
	type args struct {
		ctx context.Context
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
			h := &HTTP{
				srv: tt.fields.srv,
				acc: tt.fields.acc,
				ur:  tt.fields.ur,
			}
			if err := h.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("HTTP.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
