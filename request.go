package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	urlpkg "net/url"
	"strings"
)

// New creates a new RequestBuilder with the provided URL.
// It initializes the request with a GET method.
func New(url string) *RequestBuilder {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	return &RequestBuilder{req: req, err: err}
}

// RequestBuilder is a builder for http.Request.
// It provides methods to set up the request.
type RequestBuilder struct {
	retryTimes uint
	err        error
	req        *http.Request
	client     *http.Client
}

// Err returns the error that occurred while building the request.
func (r *RequestBuilder) Err() error {
	return r.err
}

// Method sets the HTTP method for the request.
func (r *RequestBuilder) Method(method string) *RequestBuilder {
	r.req.Method = method
	return r
}

// Get sets the HTTP method to GET.
func (r *RequestBuilder) Get() *RequestBuilder {
	return r.Method(http.MethodGet)
}

// Post sets the HTTP method to POST.
func (r *RequestBuilder) Post() *RequestBuilder {
	return r.Method(http.MethodPost)
}

// Put sets the HTTP method to PUT.
func (r *RequestBuilder) Put() *RequestBuilder {
	return r.Method(http.MethodPut)
}

// Patch sets the HTTP method to PATCH.
func (r *RequestBuilder) Patch() *RequestBuilder { return r.Method(http.MethodPatch) }

// Delete sets the HTTP method to DELETE.
func (r *RequestBuilder) Delete() *RequestBuilder { return r.Method(http.MethodDelete) }

// Head sets the HTTP method to HEAD.
func (r *RequestBuilder) Head() *RequestBuilder {
	return r.Method(http.MethodHead)
}

// Connect sets the HTTP method to CONNECT.
func (r *RequestBuilder) Connect() *RequestBuilder { return r.Method(http.MethodConnect) }

// Options sets the HTTP method to OPTIONS.
func (r *RequestBuilder) Options() *RequestBuilder { return r.Method(http.MethodOptions) }

// Trace sets the HTTP method to TRACE.
func (r *RequestBuilder) Trace() *RequestBuilder { return r.Method(http.MethodTrace) }

// Body sets the body for the request.
func (r *RequestBuilder) Body(body io.ReadCloser) *RequestBuilder {
	if r.err != nil {
		return r
	}
	r.req.Body = body
	switch v := body.(io.Reader).(type) {
	case *bytes.Buffer:
		r.req.ContentLength = int64(v.Len())
		buf := v.Bytes()
		r.req.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return io.NopCloser(r), nil
		}
	case *bytes.Reader:
		r.req.ContentLength = int64(v.Len())
		snapshot := *v
		r.req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	case *strings.Reader:
		r.req.ContentLength = int64(v.Len())
		snapshot := *v
		r.req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	default:
		// This is where we'd set it to -1 (at least
		// if body != NoBody) to mean unknown, but
		// that broke people during the Go 1.8 testing
		// period. People depend on it being 0 I
		// guess. Maybe retry later. See Issue 18117.
	}
	// For client requests, Request.ContentLength of 0
	// means either actually 0, or unknown. The only way
	// to explicitly say that the ContentLength is zero is
	// to set the Body to nil. But turns out too much code
	// depends on NewRequest returning a non-nil Body,
	// so we use a well-known ReadCloser variable instead
	// and have the http package also treat that sentinel
	// variable to mean explicitly zero.
	if r.req.GetBody != nil && r.req.ContentLength == 0 {
		r.req.Body = http.NoBody
		r.req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	}
	return r
}

// SetHeader sets a header for the request.
func (r *RequestBuilder) SetHeader(key, value string) *RequestBuilder {
	if r.err != nil {
		return r
	}
	r.req.Header.Set(key, value)
	return r
}

// Form sets form values for the request.
func (r *RequestBuilder) Form(values urlpkg.Values) *RequestBuilder {
	if r.err != nil {
		return r
	}
	r.req.Form = values
	return r
}

// Query sets query parameters for the request.
func (r *RequestBuilder) Query(queries map[string]string) *RequestBuilder {
	if r.err != nil {
		return r
	}
	query := urlpkg.Values{}
	for key, value := range queries {
		query.Add(key, value)
	}
	r.req.URL.RawQuery += query.Encode()
	return r
}

// AddQuery adds a single query parameter to the request.
func (r *RequestBuilder) AddQuery(key, value string) *RequestBuilder {
	if r.err != nil {
		return r
	}
	return r.Query(map[string]string{key: value})
}

// Json sets the body of the request to the JSON representation of v.
func (r *RequestBuilder) Json(v interface{}) *RequestBuilder {
	if r.err != nil {
		return r
	}
	data, err := json.Marshal(v)
	if err != nil {
		r.err = err
		return r
	}
	r.SetHeader("Content-Type", "application/json")
	return r.Body(io.NopCloser(bytes.NewBuffer(data)))
}

// PostForm sets the body of the request to the URL-encoded form data.
func (r *RequestBuilder) PostForm(values urlpkg.Values) *RequestBuilder {
	if r.err != nil {
		return r
	}
	r.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	body := strings.NewReader(values.Encode())
	return r.Body(io.NopCloser(body))
}

func (r *RequestBuilder) Retry(retryTimes uint) *RequestBuilder {
	r.retryTimes = retryTimes
	return r
}

// BuildWithContext builds the request with the provided context.
func (r *RequestBuilder) BuildWithContext(ctx context.Context) (*http.Request, error) {
	if r.err != nil {
		return nil, r.err
	}
	if ctx == context.Background() {
		return r.req, nil
	}
	return r.req.WithContext(ctx), nil
}

// Build builds the request with a background context.
func (r *RequestBuilder) Build() (*http.Request, error) {
	return r.BuildWithContext(context.Background())
}

// Do send the request and returns the response.
func (r *RequestBuilder) Do() (resp *http.Response, err error) {
	req, err := r.Build()
	if err != nil {
		return nil, err
	}

	var client = r.client
	if client == nil {
		client = http.DefaultClient
	}

	retryTimes := r.retryTimes
	if retryTimes == 0 {
		retryTimes = 1
	}

	for i := 0; i < int(retryTimes); i++ {
		resp, err = client.Do(req)
		if err == nil {
			return resp, nil
		}
	}
	return nil, err
}
