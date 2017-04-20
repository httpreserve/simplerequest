package simplerequest

import (
	"net/http"
	"net/url"
	"time"
)

// HTTP request methods that are useful to us
const GET = http.MethodGet
const HEAD = http.MethodHead

// Default byte-range for initial requests
// use first 0.5mb to find <title> tag...
const partByterange = "bytes=0-500"

const br = "bytes=0-"

// SimpleRequest structure to be turned into a
// HTTP request proper in code.
type SimpleRequest struct {
	Method    string
	URL       *url.URL
	Proxy     bool
	agent     string
	byterange string
	timeout   time.Duration
}

func (sr *SimpleRequest) Agent(agent string) {
	sr.agent = agent
}

func (sr *SimpleRequest) Byterange(limit string) {
	sr.byterange = br + limit
}

type SimpleResponse struct {
	Data           string
	Header         http.Header
	PrettyResponse string
	PrettyRequest  string
	StatusText     string
	StatusCode     int
}
