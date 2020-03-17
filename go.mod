module github.com/lvl484/user-manager

go 1.13
//We don't need this depdendency
//But docker conteiner didn't start without any dependency in  go.mod
require github.com/gorilla/mux v1.6.2
