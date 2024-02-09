package handler

import (
	"net/http"

	http_errors "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/http"
	"github.com/gorilla/mux"
)

func (errors *http_errors) routes() *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(http_errors.NotFound())
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowed)

	router.Use(app.logAccess)
	router.Use(app.recoverPanic)

	router.HandleFunc("/status", status).Methods("GET")
	return router
}
