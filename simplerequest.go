package simplerequest

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Create is a mechanism to make a suitable http request header
// to find some information out about a web resouse.
func Create(method string, reqURL string) (SimpleRequest, error) {
	var sr SimpleRequest
	sr.Method = method
	req, err := url.Parse(reqURL)
	if err != nil {
		return sr, errors.Wrap(err, "url parse failed in CreateSimpleRequest")
	}
	sr.URL = req
	return sr, nil
}

// Default will create a default object to make work easier for all.
func Default(reqURL *url.URL) SimpleRequest {
	// we're not concerned about error here, as internally, we've
	// already parsed the URL which is the only source of potential
	// error in CreateSimpleRequest
	sr, _ := Create(GET, reqURL.String())
	return sr
}

// Do will perform the request for us from the
// SimpleRequest object itself.
func (sr *SimpleRequest) Do() (SimpleResponse, error) {
	resp, err := sr.handlehttp(sr.Method, sr.URL)
	return resp, err
}

// Timeout sets the client timeout in seconds...
func (sr *SimpleRequest) Timeout(duration time.Duration) {
	sr.timeout = time.Duration(duration * time.Second)
}

// Redirect tells the client whether or not to follow redirects...
func (sr *SimpleRequest) NoRedirect(redirect bool) {
	sr.noredirect = redirect
}

// GetHeader can be used to retrieve headers from the server's response
func (sr *SimpleResponse) GetHeader(key string) string {
	return sr.Header.Get(key)
}

// prettyRequest is used internally to pretty print the request
func prettyRequest(sr SimpleResponse, req *http.Request) SimpleResponse {
	pr, _ := httputil.DumpRequest(req, false)
	sr.PrettyRequest = string(pr)
	return sr
}

// prettyResponse is used internally to pretty print the response
func prettyResponse(sr SimpleResponse, resp *http.Response) SimpleResponse {
	pr, _ := httputil.DumpResponse(resp, false)
	sr.PrettyResponse = string(pr)
	return sr
}

// status is used internally to populate status elements of SimpleResponse
func status(sr SimpleResponse, resp *http.Response) SimpleResponse {
	sr.StatusText = http.StatusText(resp.StatusCode)
	sr.StatusCode = resp.StatusCode
	return sr
}

// Handle the HTTP request and response functions of the calling application.
func (sr *SimpleRequest) handlehttp(method string, reqURL *url.URL) (SimpleResponse, error) {

	var simpleresponse SimpleResponse

	var client http.Client
	if sr.timeout != 0 {
		client.Timeout = sr.timeout
		if sr.noredirect {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
	}

	req, err := http.NewRequest(method, reqURL.String(), nil)
	if err != nil {
		return simpleresponse, errors.Wrap(err, "request generation failed")
	}

	if sr.agent != "" {
		req.Header.Add("User-Agent", sr.agent)
	}

	if sr.byterange != "" {
		req.Header.Add("Range", sr.byterange)
	}

	if sr.accept != "" {
		req.Header.Add("Accept", sr.accept)
	}

	resp, err := client.Do(req)
	if err != nil {
		return simpleresponse, errors.Wrap(err, "client.Do failed")
	}

	// once we've closed the body we can't do anything else
	// with the packet content...
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return simpleresponse, errors.Wrap(err, "reading http response body")
	}

	// Populate the remainder of the SimpleResponse elements
	simpleresponse.Header = resp.Header
	simpleresponse.Data = string(data)
	simpleresponse = prettyRequest(simpleresponse, req)
	simpleresponse = prettyResponse(simpleresponse, resp)
	simpleresponse = status(simpleresponse, resp)

	return simpleresponse, err
}
