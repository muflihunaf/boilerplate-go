package handler

import (
	"net/http"

	"github.com/muflihunaf/boilerplate-go/pkg/response"
)

// Response helpers - thin wrappers for cleaner handler code.
// For full functionality, import pkg/response directly.

func JSON(w http.ResponseWriter, status int, data interface{}) {
	response.Success(w, status, data)
}

func OK(w http.ResponseWriter, data interface{}) {
	response.OK(w, data)
}

func Created(w http.ResponseWriter, data interface{}) {
	response.Created(w, data)
}

func BadRequest(w http.ResponseWriter, msg string) {
	response.BadRequest(w, msg)
}

func Unauthorized(w http.ResponseWriter, msg string) {
	response.Unauthorized(w, msg)
}

func Forbidden(w http.ResponseWriter, msg string) {
	response.Forbidden(w, msg)
}

func NotFound(w http.ResponseWriter, msg string) {
	response.NotFound(w, msg)
}

func Conflict(w http.ResponseWriter, msg string) {
	response.Conflict(w, msg)
}

func InternalError(w http.ResponseWriter) {
	response.InternalError(w)
}

func Error(w http.ResponseWriter, status int, code, msg string) {
	response.Error(w, status, code, msg)
}
