package handler

import (
	"github.com/muflihunaf/boilerplate-go/pkg/response"
	"net/http"
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

func ValidationError(w http.ResponseWriter, message string) {
	response.ValidationError(w, message)
}
