package httphandler

import (
	"github.com/gorilla/mux"
	"github.com/romanzac/gorilla-feast/middleware"
	"net/http"
)

// InitRoutes for Gorilla Feast
func InitRoutes(r *mux.Router, apiv1 *APIv1) {

	// Test route
	r.HandleFunc("/ping", apiv1.PingPong)

	// WebSocket routes
	r.HandleFunc("/login-failures", apiv1.LoginFailures)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	// User routes
	v1.HandleFunc("/login", apiv1.Login).
		Methods("POST")

	v1.HandleFunc("/user", apiv1.SignupUser).
		Methods("POST")

	v1.Handle("/user",
		middleware.JWTHandler(http.HandlerFunc(apiv1.ListAllUsers))).
		Methods("GET")

	v1.Handle("/user/{fullname}",
		middleware.JWTHandler(http.HandlerFunc(apiv1.SearchUserbyFullname))).
		Methods("GET")

	v1.Handle("/user/{acct}/detail",
		middleware.JWTHandler(http.HandlerFunc(apiv1.GetUserDetail))).
		Methods("GET")

	v1.Handle("/user/{acct}",
		middleware.JWTHandler(http.HandlerFunc(apiv1.UpdateUser))).
		Methods("PATCH")

	v1.Handle("/user/{acct}",
		middleware.JWTHandler(http.HandlerFunc(apiv1.DeleteUser))).
		Methods("DELETE")
}
