package mockhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Transport struct {
	*ResponseMock
	MockError error
}

var _ http.RoundTripper = &Transport{}

func (c *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.MockError != nil {
		return nil, c.MockError
	}

	return c.ResponseMock.MakeResponse(req), nil
}

var emptyBytes = []byte{}

func NewResponseMock(statusCode int, headers map[string]string, body []byte) *ResponseMock {
	if headers == nil {
		headers = map[string]string{}
	}
	if body == nil {
		body = emptyBytes
	}
	return &ResponseMock{StatusCode: statusCode, HeadersMap: headers, Body: body}
}

type ResponseMock struct {
	StatusCode int
	HeadersMap map[string]string
	Body       []byte
}

func (r *ResponseMock) MakeTransport() *Transport {
	return &Transport{ResponseMock: r}
}

func (r *ResponseMock) MakeClient() *http.Client {
	return &http.Client{Transport: r.MakeTransport()}
}

func (r *ResponseMock) MakeResponse(req *http.Request) *http.Response {
	status := strconv.Itoa(r.StatusCode) + " " + http.StatusText(r.StatusCode)
	header := http.Header{}
	for name, value := range r.HeadersMap {
		header.Set(name, value)
	}

	contentLength := len(r.Body)
	header.Set("Content-Length", strconv.Itoa(contentLength))

	res := &http.Response{
		Status:           status,
		StatusCode:       r.StatusCode,
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           header,
		Body:             ioutil.NopCloser(bytes.NewReader(r.Body)),
		ContentLength:    int64(contentLength),
		TransferEncoding: []string{},
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          req,
		TLS:              nil,
	}

	// should no set Content-Length header when 204 or 304
	if r.StatusCode == http.StatusNoContent || r.StatusCode == http.StatusNotModified {
		if res.ContentLength != 0 {
			res.Body = ioutil.NopCloser(bytes.NewReader(emptyBytes))
			res.ContentLength = 0
		}
		header.Del("Content-Length")
	}

	return res
}
