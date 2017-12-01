package mockhttp

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMockClient(t *testing.T) {
	client := NewResponseMock(http.StatusOK, map[string]string{"Content-Type": "text/plain"}, []byte("hello")).MakeClient()
	res, err := client.Get("http://example.com/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "hello" {
		t.Errorf("body should be hello, but got %s", string(body))
	}
}

func TestMockTransport(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	transport := NewResponseMock(http.StatusOK, map[string]string{"Content-Type": "text/plain"}, []byte("hello")).MakeTransport()
	res, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Request != req {
		t.Errorf("Request should be original request, but got %+v", res.Request)
	}

	transport.MockError = errors.New("mocked")
	res, err = transport.RoundTrip(req)
	if err != transport.MockError {
		t.Fatal(err)
	}
	if res != nil {
		t.Errorf("Response should be nil, but got %+v", res)
	}
}

func TestResponseMock(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := NewResponseMock(http.StatusOK, map[string]string{"Content-Type": "text/plain"}, []byte("hello")).MakeResponse(req)
	if res.Status != "200 OK" {
		t.Errorf("Status should be 200 OK, but got %s", res.Status)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("StatusCode should be 200, but got %d", res.StatusCode)
	}
	if res.Proto != "HTTP/1.0" {
		t.Errorf("Proto should be HTTP/1.0, but got %d", res.Proto)
	}
	if res.ProtoMajor != 1 {
		t.Errorf("ProtoMajor should be 1, but got %d", res.ProtoMajor)
	}
	if res.ProtoMinor != 0 {
		t.Errorf("ProtoMinor should be 0, but got %d", res.ProtoMinor)
	}

	if res.Header == nil {
		t.Error("Header should not be nil")
	}
	if headers := len(res.Header); headers != 2 {
		t.Errorf("Headers count should be 2, but got %d", headers)
	}
	if ct := res.Header.Get("Content-Type"); ct != "text/plain" {
		t.Errorf("Content-Type header should be text/plain, but got %s", ct)
	}
	if cl := res.Header.Get("Content-Length"); cl != "5" {
		t.Errorf("Content-Length header should be 5, but got %s", cl)
	}

	if res.Body == nil {
		t.Error("Body should not be nil")
	}
	if body, err := ioutil.ReadAll(res.Body); body == nil || err != nil || string(body) != "hello" {
		t.Errorf(`Body should be "hello", but got: "%s"`, string(body))
	}

	if res.ContentLength != 5 {
		t.Errorf("ContentLength should be 5, but got %d", res.ContentLength)
	}
	if res.TransferEncoding == nil {
		t.Error("TransferEncoding should not be nil")
	}
	if len(res.TransferEncoding) != 0 {
		t.Errorf("TransferEncoding should be empty, but got: %+v", res.TransferEncoding)
	}
	if res.Close != false {
		t.Error("Close should be false, but got true")
	}
	if res.Uncompressed != false {
		t.Error("Uncompressed should be false, but got true")
	}
	if res.Trailer != nil {
		t.Error("Trailer should not be nil")
	}
	if res.Request != req {
		t.Errorf("Request should be original request, but got %+v", res.Request)
	}
	if res.TLS != nil {
		t.Error("TLS should not be nil")
	}

	t.Run("Empty", func(t *testing.T) {
		res := NewResponseMock(http.StatusOK, nil, nil).MakeResponse(req)
		if res.Header == nil {
			t.Error("Header should not be nil")
		}
		if headers := len(res.Header); headers != 1 {
			t.Errorf("Headers count should be 1, but got %d", headers)
		}
		if cl := res.Header.Get("Content-Length"); cl != "0" {
			t.Errorf("Content-Length header should be 0, but got %s", cl)
		}
		if res.ContentLength != 0 {
			t.Errorf("ContentLength should be 0, but got %d", res.ContentLength)
		}

		if res.Body == nil {
			t.Error("Body should not be nil")
		}
		if body, err := ioutil.ReadAll(res.Body); body == nil || err != nil || len(body) != 0 {
			t.Errorf(`Body should be empty, but got: "%s"`, string(body))
		}
	})

	t.Run("StatusNoContent", func(t *testing.T) {
		res := NewResponseMock(http.StatusNoContent, nil, nil).MakeResponse(req)
		if res.Header == nil {
			t.Error("Header should not be nil")
		}
		if headers := len(res.Header); headers != 0 {
			t.Errorf("Headers count should be 0, but got %d", headers)
		}
		if res.ContentLength != 0 {
			t.Errorf("ContentLength should be 0, but got %d", res.ContentLength)
		}

		if res.Body == nil {
			t.Error("Body should not be nil")
		}
		if body, err := ioutil.ReadAll(res.Body); body == nil || err != nil || len(body) != 0 {
			t.Errorf(`Body should be empty, but got: "%s"`, string(body))
		}
	})

	t.Run("StatusNotModified", func(t *testing.T) {
		res := NewResponseMock(http.StatusNotModified, nil, []byte("should ignore this body")).MakeResponse(req)
		if res.Header == nil {
			t.Error("Header should not be nil")
		}
		if headers := len(res.Header); headers != 0 {
			t.Errorf("Headers count should be 0, but got %d", headers)
		}
		if res.ContentLength != 0 {
			t.Errorf("ContentLength should be 0, but got %d", res.ContentLength)
		}

		if res.Body == nil {
			t.Error("Body should not be nil")
		}
		if body, err := ioutil.ReadAll(res.Body); body == nil || err != nil || len(body) != 0 {
			t.Errorf(`Body should be empty, but got: "%s"`, string(body))
		}
	})
}
