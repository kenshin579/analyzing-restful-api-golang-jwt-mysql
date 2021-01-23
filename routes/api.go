package routes

import (
	"github.com/gorilla/mux"
	"github.com/kenshin579/analyzing-restful-api-golang-jwt-mysql/app/controllers"
	"github.com/kenshin579/analyzing-restful-api-golang-jwt-mysql/app/controllers/auth"
)

func ApiRoutes(prefix string, r *mux.Router) {

	s := r.PathPrefix(prefix).Subrouter()

	s.HandleFunc("/login", auth.Login).Methods("POST")
	s.HandleFunc("/register", controllers.CreateUser).Methods("POST")

	s.HandleFunc("/users", auth.ValidateMiddleware(controllers.GetUsers)).Methods("GET")
	s.HandleFunc("/users/{id}", auth.ValidateMiddleware(controllers.GetUser)).Methods("GET")
	s.HandleFunc("/users", auth.ValidateMiddleware(controllers.CreateUser)).Methods("POST")
	s.HandleFunc("/users/{id:[0-9]+}", auth.ValidateMiddleware(controllers.GetUser)).Methods("GET")
	s.HandleFunc("/users/{id:[0-9]+}", auth.ValidateMiddleware(controllers.UpdateUser)).Methods("PUT")
	s.HandleFunc("/users/{id}", auth.ValidateMiddleware(controllers.DeleteUser)).Methods("DELETE")
}
