package handler

import (
	"net/http"

	"github.com/muflihunaf/boilerplate-go/pkg/response"
)

// Re-export response types for handler convenience
type Response = response.Response
type ErrorInfo = response.ErrorInfo
type Meta = response.Meta

// Re-export response functions
var (
	JSON         = response.JSON
	JSONWithMeta = response.JSONWithMeta
	Error        = response.Error
	OK           = response.OK
	Created      = response.Created
	NoContent    = response.NoContent
	Paginated    = response.Paginated
)

// Common error responses
func BadRequest(w http.ResponseWriter, message string) {
	response.BadRequest(w, message)
}

func NotFound(w http.ResponseWriter, message string) {
	response.NotFound(w, message)
}

func InternalError(w http.ResponseWriter) {
	response.InternalError(w)
}

func Unauthorized(w http.ResponseWriter, message string) {
	response.Unauthorized(w, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	response.Forbidden(w, message)
}

func Conflict(w http.ResponseWriter, message string) {
	response.Conflict(w, message)
}

func ValidationError(w http.ResponseWriter, message string, details map[string]string) {
	response.ValidationError(w, message, details)
}
