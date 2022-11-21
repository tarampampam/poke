package addons

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	js "github.com/dop251/goja"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Fetch struct {
	ctx    context.Context
	client httpClient
}

func NewFetch(ctx context.Context, client httpClient) *Fetch {
	const defaultTimeout = time.Second * 60

	if client == nil {
		client = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return &Fetch{ctx: ctx, client: client}
}

func (f Fetch) Register(runtime *js.Runtime) error {
	return runtime.Set("fetchSync", f.fetch(runtime))
}

// https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
func (f Fetch) fetch(runtime *js.Runtime) func(call js.FunctionCall) js.Value { //nolint:funlen,gocognit
	return func(call js.FunctionCall) js.Value {
		if len(call.Arguments) == 0 {
			panic(runtime.ToValue("Wrong arguments count for the fetch function call"))
		}

		var (
			url     = call.Argument(0).String()
			options *js.Object
		)

		if len(call.Arguments) > 1 {
			options = call.Argument(1).ToObject(runtime)
		} else {
			options = runtime.NewObject()
		}

		var ( // defaults
			method            = http.MethodGet
			headers           = make(http.Header)
			body    io.Reader = http.NoBody
		)

		headers.Set("User-Agent", "Mozilla/5.0 (X11) Gecko/20100101 Firefox/106.0") // default user-agent

		if methodValue := options.Get("method"); methodValue != nil {
			method = strings.ToUpper(methodValue.String())
		}

		if headersValue := options.Get("headers"); headersValue != nil {
			if headersMap, isMap := headersValue.Export().(map[string]any); isMap {
				for name, headerValue := range headersMap {
					if value, isString := headerValue.(string); isString {
						headers.Set(name, value)
					}
				}
			}
		}

		if bodyValue := options.Get("body"); bodyValue != nil {
			body = bytes.NewBufferString(bodyValue.String())
		}

		var result = fetchResponse{
			runtime: runtime,
			Headers: make(map[string]string),
			URL:     url,
		}

		req, err := http.NewRequestWithContext(f.ctx, method, url, body)
		if err != nil {
			result.setStatusCode(http.StatusInternalServerError)
			result.Body = err.Error()

			return runtime.ToValue(result)
		}

		req.Header = headers

		req.UserAgent()

		resp, err := f.client.Do(req)
		if err != nil {
			result.setStatusCode(http.StatusInternalServerError)
			result.Body = err.Error()

			return runtime.ToValue(result)
		}

		defer func() { _ = resp.Body.Close() }()

		responseBody, _ := io.ReadAll(resp.Body)

		result.setStatusCode(resp.StatusCode)
		result.Body = string(responseBody)
		result.OK = resp.StatusCode >= 200 && resp.StatusCode < 300 //nolint:gomnd

		for name, v := range resp.Header {
			result.Headers[name] = strings.Join(v, ", ")
		}

		return runtime.ToValue(result)
	}
}

type fetchResponse struct { // https://developer.mozilla.org/en-US/docs/Web/API/Response
	runtime *js.Runtime

	// Body contents
	Body string `json:"body"`

	// The Headers object associated with the response
	Headers map[string]string `json:"headers"`

	// A boolean indicating whether the response was successful (status in the range 200 – 299) or not
	OK bool `json:"ok"`

	// Indicates whether the response is the result of a redirect (that is, its URL list has more than one entry)
	Redirected bool `json:"redirected"` // TODO: ❌ not implemented

	// The status code of the response (this will be 200 for a success)
	Status int `json:"status"`

	// The status message corresponding to the status code (e.g., OK for 200)
	StatusText string `json:"statusText"`

	// The type of the response (e.g., basic, cors)
	Type string `json:"type"` // TODO: ❌ not implemented

	// The URL of the response
	URL string `json:"url"`
}

func (r *fetchResponse) setStatusCode(code int) {
	r.Status = code
	r.StatusText = http.StatusText(code)
}

// ArrayBuffer returns body data as an ArrayBuffer.
// https://developer.mozilla.org/en-US/docs/Web/API/Response/arrayBuffer
func (r *fetchResponse) ArrayBuffer(_ js.FunctionCall) js.Value {
	return r.runtime.ToValue(r.runtime.NewArrayBuffer([]byte(r.Body)))
}

// Blob returns a Blob representation of the response body.
// https://developer.mozilla.org/en-US/docs/Web/API/Response/blob
func (r *fetchResponse) Blob(_ js.FunctionCall) js.Value { return js.Undefined() } // TODO: ❌ not implemented

// Clone returns a clone of a fetchResponse object.
// https://developer.mozilla.org/en-US/docs/Web/API/Response/clone
func (r *fetchResponse) Clone(_ js.FunctionCall) js.Value { return js.Undefined() } // TODO: ❌ not implemented

// FormData returns a FormData representation of the response body.
// https://developer.mozilla.org/en-US/docs/Web/API/Response/formData
func (r *fetchResponse) FormData(_ js.FunctionCall) js.Value { return js.Undefined() } // TODO: ❌ not implemented

// Json returns a result of parsing the response body text as JSON.
// https://developer.mozilla.org/en-US/docs/Web/API/Response/json
func (r *fetchResponse) Json(_ js.FunctionCall) js.Value {
	var value any

	if err := json.Unmarshal([]byte(r.Body), &value); err == nil {
		return r.runtime.ToValue(value)
	} else {
		panic(r.runtime.ToValue("Wrong JSON: " + err.Error()))
	}
}

// Text returns a text representation of the response body.
// https://developer.mozilla.org/en-US/docs/Web/API/Response/text
func (r *fetchResponse) Text(_ js.FunctionCall) js.Value {
	return r.runtime.ToValue(r.Body)
}
