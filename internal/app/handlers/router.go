package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	// to-do: implement status, healtcheck and monitoring
	// router.HandleFunc("/status", status).Methods("GET")
	return router
}
