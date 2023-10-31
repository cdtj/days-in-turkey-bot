package httpserver

import "net/http"

type HttpServerRouter interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handlerFn http.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
