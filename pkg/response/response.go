// Package response provides a standard API response format.
// Framework-agnostic - only depends on net/http.
package response

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo contains error details.
type ErrorInfo struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Meta contains pagination info.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// --- Success Responses ---

func Success(w http.ResponseWriter, status int, data interface{}) {
	write(w, status, Response{Success: true, Data: data})
}

func OK(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusCreated, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Paginated(w http.ResponseWriter, data interface{}, page, perPage int, total int64) {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}
	write(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta:    &Meta{Page: page, PerPage: perPage, Total: total, TotalPages: totalPages},
	})
}

// --- Error Responses ---

func Error(w http.ResponseWriter, status int, code, message string) {
	write(w, status, Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	})
}

func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", msg)
}

func Unauthorized(w http.ResponseWriter, msg string) {
	Error(w, http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

func Forbidden(w http.ResponseWriter, msg string) {
	Error(w, http.StatusForbidden, "FORBIDDEN", msg)
}

func NotFound(w http.ResponseWriter, msg string) {
	Error(w, http.StatusNotFound, "NOT_FOUND", msg)
}

func Conflict(w http.ResponseWriter, msg string) {
	Error(w, http.StatusConflict, "CONFLICT", msg)
}

func ValidationError(w http.ResponseWriter, msg string, details map[string]string) {
	write(w, http.StatusUnprocessableEntity, Response{
		Success: false,
		Error:   &ErrorInfo{Code: "VALIDATION_ERROR", Message: msg, Details: details},
	})
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
}

// --- Internal ---

func write(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
