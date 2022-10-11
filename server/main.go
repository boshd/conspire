package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cohere-ai/cohere-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	cohere *cohere.Client
	router *mux.Router
}

func main() {
	s := NewServer()

	s.ListenAndServe()
}

func (s *Server) ListenAndServe() error {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Connection", "Upgrade", "Sec-WebSocket-Extensions", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Accept-Encoding", "Accept-Language"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowCreds := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	return http.ListenAndServe(
		fmt.Sprintf(":%d", 8080),
		// "localhost.cer.pem",
		// "localhost.key.pem",
		handlers.CORS(headers, origins, methods, allowCreds)(s.router),
	)
}

func printEndpoints(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		fmt.Printf("%v %s\n", methods, path)
		return nil
	})
}
func NewServer() *Server {
	s := new(Server)
	log.Println("configuring cohere client...")
	s.configureCohereClient()

	log.Println("configuring http router w/ the following routes...")
	s.configureRouter()
	printEndpoints(s.router)
	return s
}

func (s *Server) configureCohereClient() {
	co, err := cohere.CreateClient("8Qfv0c6ffS5AIh5n8y2lAhEcalAwxlit6bWQajsg")
	if err != nil {
		fmt.Println(err)
		return
	}
	s.cohere = co
}

func (s *Server) configureRouter() error {
	s.router = mux.NewRouter().StrictSlash(true)
	r, err := NewRoutes(s.cohere)

	if err != nil {
		return err
	}

	s.bindRoutes(r)

	return nil
}

// bindRoutes adds all routes to the server's router.
func (s *Server) bindRoutes(r []Route) {
	for _, route := range r {
		httpHandler := s.makeHTTPHandler(route)

		s.router.
			Methods(route.Method()).
			Path(route.Pattern()).
			HandlerFunc(httpHandler)
	}
}

// makeHTTPHandler creates a http.HandlerFunc from a custom http function and logs the error if
// exists: func(http.ResponseWriter, *http.Request) error.
func (s *Server) makeHTTPHandler(route Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc := s.handleWithMiddleware(route)
		// handlerFunc := route.HandlerFunc()

		err := handlerFunc(w, r)

		if err != nil {
			log.Printf("Handler [%s][%s] returned error: %s", r.Method, r.URL.Path, err)
		}
	}
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (s *Server) handleWithMiddleware(route Route) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			handler = route.HandlerFunc()
			ctx     = r.Context()
		)

		// ctx = context.WithValue(ctx, "cookieStore", s.cookieStore)
		// ctx = context.WithValue(ctx, "config", s.cfg)
		r = r.WithContext(ctx)

		handler = HandleHTTPError(handler)

		// if route.RequiresAuth() {
		// 	handler = middleware.ValidateAuth(handler)
		// }

		return handler(w, r)
	}
}

// HandleHTTPError sets the appropriate headers to the response if a http
// handler returned an error. This is used for different types of errors
// are returned.
func HandleHTTPError(h HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := h(w, r)

		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}

		return err
	}
}

func (s *Server) Run() {
	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
