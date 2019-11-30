package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var methodAll = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

type middleware func(http.Handler) http.Handler

type Router struct {
	mx          *httprouter.Router
	middlewares []middleware
}

func (r *Router) Use(middlewares ...middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *Router) GET(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodGet, path, r.chain(handle))
}

func (r *Router) HEAD(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodHead, path, r.chain(handle))
}

func (r *Router) POST(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodPost, path, r.chain(handle))
}

func (r *Router) PUT(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodPut, path, r.chain(handle))
}

func (r *Router) PATCH(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodPatch, path, r.chain(handle))
}

func (r *Router) DELETE(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodDelete, path, r.chain(handle))
}

func (r *Router) CONNECT(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodConnect, path, r.chain(handle))
}

func (r *Router) OPTIONS(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodOptions, path, r.chain(handle))
}

func (r *Router) TRACE(path string, handle http.HandlerFunc) {
	r.mx.Handler(http.MethodTrace, path, r.chain(handle))
}

func (r *Router) ALL(path string, handle http.HandlerFunc) {
	for _, m := range methodAll {
		r.mx.Handler(m, path, r.chain(handle))
	}
}

func (r *Router) MethodNotAllowed(handle http.HandlerFunc) {
	r.mx.MethodNotAllowed = r.chain(handle)
}

func (r *Router) PanicHandler(handle func(http.ResponseWriter, *http.Request, interface{})) {
	r.mx.PanicHandler = handle
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mx.ServeHTTP(w, req)
}

func (r *Router) chain(handle http.Handler) http.Handler {
	l := len(r.middlewares) - 1
	for i := range r.middlewares {
		handle = r.middlewares[l-i](handle)
	}
	return handle
}

func New() *Router {
	return &Router{
		mx:          httprouter.New(),
		middlewares: make([]middleware, 0),
	}
}
