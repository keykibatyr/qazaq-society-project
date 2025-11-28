package controllers

import "net/http"

type Template interface {
	ExecuteTemplate(w http.ResponseWriter, r *http.Request, data interface{})
}