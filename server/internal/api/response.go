package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, statusCode int, errorRes []string, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	type errBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	res := errBody{
		Code:    errorRes[0],
		Message: errorRes[1],
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error: couldn't encode response body - %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{error: \"Something went wrong\"}"))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}

func successResponse(w http.ResponseWriter, statusCode int, res any) {
	w.Header().Set("Content-Type", "application/json")

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error: couldn't encode response body - %v\n", err)
		errorResponse(w, http.StatusInternalServerError, []string{"INTERNAL_SERVER_ERROR", "Something went wrong"}, err)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}
