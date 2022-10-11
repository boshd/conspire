package main

import (
	"net/http"

	"github.com/boshd/conspire/coherey"
	cohereservices "github.com/boshd/conspire/coherey"
	"github.com/cohere-ai/cohere-go"
)

type Route interface {
	Pattern() string
	Method() string
	HandlerFunc() func(http.ResponseWriter, *http.Request) error
}

type route struct {
	pattern     string
	method      string
	handlerFunc func(http.ResponseWriter, *http.Request) error
}

func (r *route) Pattern() string {
	return r.pattern
}

func (r *route) Method() string {
	return r.method
}

func (r *route) HandlerFunc() func(http.ResponseWriter, *http.Request) error {
	return r.handlerFunc
}

func NewRoutes(cohere *cohere.Client) ([]Route, error) {
	var (
		cohereService = cohereservices.NewService(cohere)
		cohereHandler = coherey.NewHandler(cohereService)
	)

	return []Route{
		&route{
			pattern:     "/theories",
			method:      "POST",
			handlerFunc: cohereHandler.HandleGenerateConspiracyTheories,
		},
	}, nil
}
