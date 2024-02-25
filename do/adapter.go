package do

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"unicode/utf8"
)

type FuncEventHTTP struct {
	Body         string            `json:"body"`
	Headers      map[string]string `json:"headers"`
	BodyIsBase64 bool              `json:"isBase64Encoded"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	Query        string            `json:"queryString"`
}

type FuncEvent struct {
	HTTP FuncEventHTTP `json:"http"`
}

func (e *FuncEvent) Request() (*http.Request, error) {
	body := []byte(e.HTTP.Body)

	if e.HTTP.BodyIsBase64 {
		_, err := base64.StdEncoding.Decode(body, []byte(e.HTTP.Body))

		if err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(e.HTTP.Path)
	u.RawQuery = e.HTTP.Query

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(e.HTTP.Method, u.String(), bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	for header, value := range e.HTTP.Headers {
		req.Header.Set(header, value)
	}

	return req, nil
}

type FuncEventMap map[string]interface{}

func (em *FuncEventMap) Event() (*FuncEvent, error) {
	res := &FuncEvent{}

	str, err := json.Marshal(em)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(str, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

type FuncResponseMap map[string]interface{}

type FuncResponse struct {
	Body       string            `json:"body,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

func (s *FuncResponse) Map() FuncResponseMap {
	res := FuncResponseMap{}

	str, _ := json.Marshal(s)
	_ = json.Unmarshal(str, &res)

	return res
}

type FuncResponseWriter struct {
	Body       []byte      `json:"body,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
}

func (w *FuncResponseWriter) Header() http.Header {
	if w.Headers == nil {
		w.Headers = make(http.Header)
	}

	return w.Headers
}

func (w *FuncResponseWriter) Write(b []byte) (int, error) {
	w.Body = make([]byte, len(b))

	copy(w.Body, b)

	return len(b), nil
}

func (w *FuncResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

// Deprecated: use FuncResponse() instead. Will be removed in v2.0.
func (w *FuncResponseWriter) GetFuncResponse() *FuncResponse {
	return w.FuncResponse()
}

func (w *FuncResponseWriter) FuncResponse() *FuncResponse {
	res := &FuncResponse{}
	res.Headers = make(map[string]string)
	res.StatusCode = w.StatusCode

	for header := range w.Headers {
		res.Headers[header] = w.Headers.Get(header)
	}

	if utf8.Valid(w.Body) {
		res.Body = string(w.Body)
	} else {
		res.Body = base64.StdEncoding.EncodeToString(w.Body)
	}

	return res
}
