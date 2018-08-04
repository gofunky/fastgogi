package fastgogi

import (
	"bytes"
	"errors"
	"github.com/francoispqt/onelog"
	"github.com/valyala/fasthttp"
	"os"
	"sort"
	"strings"
	"time"
)

var log *onelog.Logger

func init() {
	log = onelog.New(
		os.Stdout,
		onelog.WARN|onelog.ERROR|onelog.FATAL,
	)
}

// UseLogger sets a custom Logger for the package.
func UseLogger(logger *onelog.Logger) {
	log = logger
}

const (
	// Version of FastGogi
	Version = "1.0.0"
	// StartComment indicates the beginning of the gitignore template
	StartComment = "# Created by "
	// EndComment indicates the end of the gitignore template
	EndComment       = "# End of "
	defaultUserAgent = "fastgogi/" + Version
	defaultHost      = "https://www.gitignore.io"
	slash            = "/"
	comma            = ","
	typePath         = "api"
	listPath         = "list"
	errorConst       = "error"
)

var (
	startBytes = []byte(StartComment)
	endBytes   = []byte(EndComment)
)

type (
	// fastGogiClient contains the http client and its parameters.
	fastGogiClient struct {
		client *fasthttp.Client
		*FastGogiOptions
	}
	// FastGogiOptions contain the http client parameters.
	FastGogiOptions struct {
		UserAgent string
		Host      string
	}
)

// NewClientWithOptions is the factory for a new fastGogiClient.
func NewClientWithOptions(options FastGogiOptions) (client *fastGogiClient) {
	client = &fastGogiClient{&fasthttp.Client{}, &options}
	if client.Host == "" {
		client.Host = defaultHost
	}
	if client.UserAgent == "" {
		client.UserAgent = defaultUserAgent
	}
	client.client.ReadTimeout = time.Second
	client.client.WriteTimeout = time.Second
	client.client.MaxIdleConnDuration = time.Second
	return client
}

// NewClient returns a default fastGogiClient.
func NewClient() (client *fastGogiClient) {
	return NewClientWithOptions(FastGogiOptions{})
}

// List all available gitignore types.
func (c *fastGogiClient) List() (types []string, err error) {
	uri := strings.Join([]string{c.Host, typePath, listPath}, slash)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetUserAgent(c.UserAgent)
	resp := fasthttp.AcquireResponse()
	err = c.client.Do(req, resp)
	if err != nil {
		log.ErrorWith("gitignore.io api request failed").
			String("uri", uri).Err(errorConst, err).Write()
		return nil, err
	}
	body := string(resp.Body())
	normalizedBody := strings.ToLower(strings.Replace(body, "\n", comma, -1))
	types = strings.Split(normalizedBody, comma)
	if !contains(types, "go") {
		log.ErrorWith("gitignore.io api response could not be parsed").
			String("uri", uri).Err(errorConst, err).Write()
		return nil, errors.New("api response parsing failed")
	}
	sort.Strings(types)
	return types, nil
}

// GetPath generates the URI to gitignore.io given the input types.
func (c *fastGogiClient) GetPath(includedTypes ...string) (URI string) {
	sort.Strings(includedTypes)
	renderedTypes := strings.Join(includedTypes, comma)
	rawURI := strings.Join([]string{c.Host, typePath, renderedTypes}, slash)
	URI = strings.ToLower(rawURI)
	return URI
}

// Get the gitignore content for the given input type.
func (c *fastGogiClient) Get(includedTypes ...string) (content []byte, err error) {
	URI := c.GetPath(includedTypes...)
	uriBytes := []byte(URI)
	req := fasthttp.AcquireRequest()
	req.SetRequestURIBytes(uriBytes)
	req.Header.SetUserAgent(c.UserAgent)
	resp := fasthttp.AcquireResponse()
	err = c.client.Do(req, resp)
	if err != nil {
		log.ErrorWith("gitignore.io api request failed").
			String("URI", URI).Err(errorConst, err).Write()
		return nil, err
	}
	body := resp.Body()
	gioStart := append(startBytes, uriBytes...)
	gioEnd := append(endBytes, uriBytes...)
	if !bytes.Contains(body, gioStart) || !bytes.Contains(body, gioEnd) {
		log.ErrorWith("gitignore.io api response is not valid").
			String("URI", URI).Err(errorConst, err).Write()
		return nil, errors.New("api response is invalid")
	}
	return body, nil
}

// contains detects whether the array s contains any e
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
