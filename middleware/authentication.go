// Package middleware provides different implementations of HTTP middleware.
// Each implementation should be of type MiddlewareFunc
// https://pkg.go.dev/github.com/gorilla/mux?tab=doc#MiddlewareFunc
//
// Example
//
// A very basic middleware which logs the URI of the request being handled could be written as
//
//		func simpleMw(next http.Handler) http.Handler {
//			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//				// Do stuff here
//				log.Println(r.RequestURI)
//				// Call the next handler, which can be another middleware in the chain, or the final handler.
//				next.ServeHTTP(w, r)
//			})
//		}
//
// Middleware can be added to a router using `Router.Use()`
//
//		r := mux.NewRouter()
//		r.HandleFunc("/", handler)
//		r.Use(simpleMw)
//
// Source: https://pkg.go.dev/github.com/gorilla/mux?tab=doc#pkg-overview
package middleware
