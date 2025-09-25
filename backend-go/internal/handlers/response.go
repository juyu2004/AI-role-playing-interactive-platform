package handlers

import (
	"encoding/json"
	"net/http"
)

type baseResponse[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func writeJSON[T any](w http.ResponseWriter, status int, data T, message string, success bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(baseResponse[T]{
		Code:    status,
		Data:    data,
		Message: message,
		Success: success,
	})
}

func ok[T any](w http.ResponseWriter, data T) {
	writeJSON(w, http.StatusOK, data, "ok", true)
}

func created[T any](w http.ResponseWriter, data T) {
	writeJSON(w, http.StatusCreated, data, "created", true)
}

func badRequest(w http.ResponseWriter, message string) {
	writeJSON[map[string]any](w, http.StatusBadRequest, map[string]any{}, message, false)
}

func unauthorized(w http.ResponseWriter, message string) {
	writeJSON[map[string]any](w, http.StatusUnauthorized, map[string]any{}, message, false)
}

func notFound(w http.ResponseWriter, message string) {
	writeJSON[map[string]any](w, http.StatusNotFound, map[string]any{}, message, false)
}

func serverError(w http.ResponseWriter, message string) {
	writeJSON[map[string]any](w, http.StatusInternalServerError, map[string]any{}, message, false)
}
