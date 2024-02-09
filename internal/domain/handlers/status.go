package handler

import (
	"net/http"

	response "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/http"
)

func status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		response.RespondWithSlugError(err, w, r)
	}
}
