package interfaces

import "github.com/gorilla/mux"

func InitRoutes(r *mux.Router, userHandler *Handler) {
	r.HandleFunc("/api/users", userHandler.Create).Methods("POST")
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler.GetByID).Methods("GET")
	r.HandleFunc("/api/users", userHandler.GetList).Methods("GET")
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler.Delete).Methods("DELETE")
	r.HandleFunc("/api/users/{id:[0-9]+}", userHandler.Update).Methods("PUT")
}
