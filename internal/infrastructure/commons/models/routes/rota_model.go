package model

import "net/http"

type RouteModel struct {
	URI              string
	Method           string
	Func             func(w http.ResponseWriter, r *http.Request)
	HasAuthenticated bool
}
