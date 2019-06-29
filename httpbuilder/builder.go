package httpbuilder

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultHTTPClient use a global client to get caching benefits
var DefaultHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

// HTTPBuilder used to build a fluent interface for http requests
type HTTPBuilder interface {
	// SetClient will set the underlying client
	SetClient(client *http.Client) HTTPBuilder

	// SetURL will change the current URL for the query
	SetURL(string) HTTPBuilder

	// SetMethod the method used to invoke
	SetMethod(string) HTTPBuilder

	// SetBody is the input of the command
	SetBody(interface{}) HTTPBuilder

	// AddHeader to the request
	AddHeader(string, string) HTTPBuilder

	// AddQuery to the request
	AddQuery(string, string) HTTPBuilder

	// SetOut is the output of the command
	SetOut(interface{}) HTTPBuilder

	// SetBearerToken will set the Authorization header with a bearer token
	SetBearerToken(string) HTTPBuilder

	// SetAuthToken will set the Authorization header
	SetAuthToken(string, string) HTTPBuilder

	// SetAuthBasic will set the Authorization header for basic auth
	SetAuthBasic(string, string) HTTPBuilder

	// SetLogger so we can print stuff
	SetLogger(func(string, ...interface{})) HTTPBuilder

	// Do the HTTP Request
	Do() (int, error)
}

type httpBuilder struct {
	client    *http.Client
	url       string
	method    string
	authType  string
	authToken string
	headers   map[string]string
	queries   map[string]string
	body      interface{}
	logger    func(string, ...interface{})
	out       interface{}
}

func (b *httpBuilder) SetClient(client *http.Client) HTTPBuilder {
	b.client = client
	return b
}

func (b *httpBuilder) SetURL(url string) HTTPBuilder {
	b.url = url
	return b
}

func (b *httpBuilder) SetMethod(method string) HTTPBuilder {
	b.method = method
	return b
}

func (b *httpBuilder) AddHeader(key, value string) HTTPBuilder {
	b.headers[key] = value
	return b
}

func (b *httpBuilder) AddQuery(key, value string) HTTPBuilder {
	b.queries[key] = value
	return b
}

func (b *httpBuilder) SetOut(d interface{}) HTTPBuilder {
	b.out = d
	return b
}

func (b *httpBuilder) SetBearerToken(token string) HTTPBuilder {
	b.authType = "Bearer"
	b.authToken = token
	return b
}

func (b *httpBuilder) SetAuthToken(authType, authToken string) HTTPBuilder {
	b.authType = authType
	b.authToken = authToken
	return b
}

func (b *httpBuilder) SetAuthBasic(username, password string) HTTPBuilder {
	raw := username + ":" + password
	authToken := base64.StdEncoding.EncodeToString([]byte(raw))

	b.authType = "Basic"
	b.authToken = authToken
	return b
}

func (b *httpBuilder) SetLogger(logger func(string, ...interface{})) HTTPBuilder {
	b.logger = logger
	return b
}

func (b *httpBuilder) SetBody(body interface{}) HTTPBuilder {
	b.body = body
	return b
}

// Do the query
func (b *httpBuilder) Do() (int, error) {
	var headers = make(map[string]string)

	var body io.Reader
	if b.body != nil {
		switch raw := b.body.(type) {
		case *string:
			body = strings.NewReader(*raw)
			headers["Content-Type"] = "text/plain"
			headers["Accept"] = "text/plain"
		case string:
			body = strings.NewReader(raw)
			headers["Content-Type"] = "text/plain"
			headers["Accept"] = "text/plain"
		case url.Values:
			body = strings.NewReader(raw.Encode())
			headers["Content-Type"] = "application/x-www-form-urlencoded"
			headers["Accept"] = "text/plain"
		default:
			input, _ := json.Marshal(raw)
			body = bytes.NewReader(input)
			headers["Content-Type"] = "application/json"
			headers["Accept"] = "application/json"
		}
	}

	b.logger("Method %s, URL %s", b.method, b.url)
	req, err := http.NewRequest(b.method, b.url, body)
	if err != nil {
		return 0, err
	}

	// Add headers
	for key, val := range b.headers {
		req.Header.Add(key, val)
	}
	// Add headers
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	query := req.URL.Query()
	// Add queries
	for key, val := range b.queries {
		query.Set(key, val)
	}
	req.URL.RawQuery = query.Encode()

	if b.authType != "" {
		auth := fmt.Sprintf("%s %s", b.authType, b.authToken)
		req.Header.Set("Authorization", strings.Trim(auth, " "))
	}

	res, err := b.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// check the status code
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return res.StatusCode, nil
	}

	// If we have an output decode it
	if b.out != nil {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return res.StatusCode, err
		}

		body := string(bodyBytes)
		b.logger(body)

		switch out := b.out.(type) {
		case *string:
			*out = body
			return res.StatusCode, nil
		default:
			return res.StatusCode, json.Unmarshal(bodyBytes, b.out)
		}
	}

	return res.StatusCode, nil
}

// New will create a new instance of a HTTPBuilder
func New() HTTPBuilder {
	return &httpBuilder{
		client:  DefaultHTTPClient,
		method:  http.MethodGet,
		headers: make(map[string]string),
		queries: make(map[string]string),
		logger:  func(string, ...interface{}) {},
	}
}
