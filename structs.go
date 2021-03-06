package simplerequest

import (
	"net/http"
	"net/url"
	"time"
)

// HTTP request methods that are useful to us

// GET whole or partial response from URL
const GET = http.MethodGet

// HEAD is used to fetch headers only from URL
const HEAD = http.MethodHead

// Default byte-range for initial requests
// use first 0.5mb to find <title> tag...
const partByterange = "bytes=0-500"

const br = "bytes=0-"

// SimpleRequest structure to be turned into a HTTP request proper
type SimpleRequest struct {
	Method     string        // Method to use to query the URL
	URL        *url.URL      // URL we're requesting
	agent      string        // Agent string to include in the request
	byterange  string        // Byterange to download partial content
	accept     string        // Content-type to accept for apis
	timeout    time.Duration // Timeout for requests to fail
	noredirect bool          // Should the client follow redirects, or inspect first response
}

// Agent is used to supply a user agent name to the request
// useful to distinguish a spam request, or mock another client
func (sr *SimpleRequest) Agent(agent string) {
	sr.agent = agent
}

// Byterange is used to enable partial download of content
// for a GET request
func (sr *SimpleRequest) Byterange(limit string) {
	sr.byterange = br + limit
}

// Accept allows us to control the response sent back
func (sr *SimpleRequest) Accept(accept string) {
	sr.accept = accept
}

// SimpleResponse to wrap useful components of a HTTP response
type SimpleResponse struct {
	Data           string      // The payload from the response
	Header         http.Header // The header from the response
	PrettyResponse string      // Pretty printed response
	PrettyRequest  string      // Pretty printed request
	StatusText     string      // Status text from the response
	StatusCode     int         // Status code from the response
	Location       *url.URL    // Location pointed to bu resource
}
