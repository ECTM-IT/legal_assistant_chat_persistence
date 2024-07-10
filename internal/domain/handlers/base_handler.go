package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BaseHandler struct{}

func (h *BaseHandler) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (h *BaseHandler) RespondWithError(w http.ResponseWriter, code int, message string) {
	h.RespondWithJSON(w, code, APIError{Code: code, Message: message})
}

func (h *BaseHandler) ParseObjectID(r *http.Request, key string, fromHeader bool) (primitive.ObjectID, error) {
	if fromHeader {
		return primitive.ObjectIDFromHex(strings.TrimSpace(r.Header.Get("Authorization")))
	}
	return primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)[key]))
}

func (h *BaseHandler) DecodeJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
