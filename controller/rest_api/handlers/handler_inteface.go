package handlers

import "net/http"

type Handler interface {
	URL() string                                      // URL returns the request URL to this handler
	Methods() []string                                // Methods returns allowed HTTP methods
	Authorize() error                                 // Authorize returns error if authentication fails
	ParseArgs(r *http.Request) (*http.Request, error) // ParseArgs parses and validates request arguments
	Process(r *http.Request) *http.Response
}
