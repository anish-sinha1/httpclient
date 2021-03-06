package httpc

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
)

func (c *httpClient) getRequestBody(contentType string, body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	switch strings.ToLower(contentType) {
	case "application/json":
		return json.Marshal(body)
	case "application/xml":
		return xml.Marshal(body)
	default:
		return json.Marshal(body)
	}
}

func (c *httpClient) do(method string, url string, headers http.Header, body interface{}) (*http.Response, error) {
	client := http.Client{}

	finalHeaders := c.getRequestHeaders(headers)
	requestBody, err := c.getRequestBody(finalHeaders.Get("Content-Type"), body)
	if err != nil {
		return nil, errors.New("error parsing request body")
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.New("unable to create http request")
	}

	req.Header = finalHeaders
	return client.Do(req)
}

func (c *httpClient) getRequestHeaders(requestHeaders http.Header) http.Header {

	// add common headers
	result := make(http.Header)
	for header, value := range c.Headers {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}

	// add custom headers
	for header, value := range requestHeaders {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}
	return result
}
