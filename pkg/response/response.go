// Package response provides a standard API response format.
// It is framework-agnostic and only depends on net/http.
package response

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response.
// All API responses follow this consistent format.
type Response struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty" swaggertype:"object"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo contains error details for failed responses.
type ErrorInfo struct {
	Code    string            `json:"code" example:"BAD_REQUEST"`
	Message string            `json:"message" example:"Invalid request body"`
	Details map[string]string `json:"details,omitempty"` // Field-level validation errors
}

// Meta contains pagination information.
type Meta struct {
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"20"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"5"`
}

// -----------------------------------------------------------------------------
// Success Responses
// -----------------------------------------------------------------------------

// Success sends a success response with data.
func Success(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, Response{
		Success: true,
		Data:    data,
	})
}

// OK sends a 200 OK response with data.
func OK(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusOK, data)
}

// Created sends a 201 Created response with data.
func Created(w http.ResponseWriter, data interface{}) {
	Success(w, http.StatusCreated, data)
}

// NoContent sends a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// SuccessWithMeta sends a success response with data and pagination metadata.
func SuccessWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	writeJSON(w, status, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Paginated sends a 200 OK response with paginated data.
func Paginated(w http.ResponseWriter, data interface{}, page, perPage int, total int64) {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	SuccessWithMeta(w, http.StatusOK, data, &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

// -----------------------------------------------------------------------------
// Error Responses
// -----------------------------------------------------------------------------

// Error sends an error response.
func Error(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorWithDetails sends an error response with field-level details.
func ErrorWithDetails(w http.ResponseWriter, status int, code, message string, details map[string]string) {
	writeJSON(w, status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest sends a 400 Bad Request error.
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized sends a 401 Unauthorized error.
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden sends a 403 Forbidden error.
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound sends a 404 Not Found error.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict sends a 409 Conflict error.
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "CONFLICT", message)
}

// UnprocessableEntity sends a 422 Unprocessable Entity error.
func UnprocessableEntity(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message)
}

// ValidationError sends a 422 error with validation details.
func ValidationError(w http.ResponseWriter, message string, details map[string]string) {
	ErrorWithDetails(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, details)
}

// InternalError sends a 500 Internal Server Error.
// Note: Always log the actual error before calling this.
func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
}

// InternalErrorWithLog sends a 500 error and logs the actual error.
func InternalErrorWithLog(w http.ResponseWriter, err error) {
	// Log is intentionally not imported here to keep pkg framework-agnostic.
	// The error should be logged by the caller before this is called.
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
}

// ServiceUnavailable sends a 503 Service Unavailable error.
func ServiceUnavailable(w http.ResponseWriter, message string) {
	Error(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}

// -----------------------------------------------------------------------------
// Internal Helpers
// -----------------------------------------------------------------------------

// writeJSON writes a JSON response to the http.ResponseWriter.
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// If encoding fails, attempt to write a basic error
		http.Error(w, `{"success":false,"error":{"code":"INTERNAL_ERROR","message":"Failed to encode response"}}`, http.StatusInternalServerError)
	}
}

// -----------------------------------------------------------------------------
// Deprecated Functions (for backward compatibility)
// -----------------------------------------------------------------------------

// JSON is deprecated. Use Success or OK instead.
// Kept for backward compatibility.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	Success(w, status, data)
}

// JSONWithMeta is deprecated. Use SuccessWithMeta instead.
// Kept for backward compatibility.
func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	SuccessWithMeta(w, status, data, meta)
}
