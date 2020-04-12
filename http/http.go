package http

import (
	"fmt"
	"net/http"

	"github.com/eriktate/divulge"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func helloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}
}

func BuildRouter(ps divulge.PostService, logger *logrus.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/post", PostPost(ps, logger))
	r.Get("/post", GetPosts(ps, logger))
	r.Get("/post/{postID}", GetPost(ps, logger))

	return r
}

// A Server knows how to listen and respond to HTTP requests.
type Server struct {
	host string
	port uint

	handler http.Handler
}

// New returns a new http.Server.
func New(host string, port uint, handler http.Handler) *Server {
	return &Server{
		host:    host,
		port:    port,
		handler: handler,
	}
}

// Start an HTTP server.
func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), s.handler)
}

func ok(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

func badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func serverError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}
