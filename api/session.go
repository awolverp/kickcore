package api

import (
	"encoding/json"
	"kickcore/logging"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 UOS"
)

type Session struct {
	logger *logging.FileLogger
	app    fasthttp.Client
}

func NewSession(logger *logging.FileLogger, readTimeout, writeTimeout time.Duration) *Session {
	s := new(Session)
	s.app = fasthttp.Client{
		Name:         defaultUserAgent,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	s.logger = logger

	return s
}

func (s *Session) log(level int, msg string, args ...interface{}) int {
	if s.logger != nil {
		return s.logger.Log(level, msg, args...)
	}
	return 0
}

type RequestConfig struct {
	Method, URI, Accept, Referer string

	CloseConnection bool
}

func (s *Session) Request(r RequestConfig, f func(*fasthttp.Response) error) error {
	// Request
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(r.Method)
	req.SetRequestURI(r.URI)

	if r.Referer != "" {
		req.Header.SetReferer(r.Referer)
	}

	if r.Accept != "" {
		req.Header.Set("Accept", r.Accept)
	} else {
		req.Header.Set("Accept", "*/*")
	}

	if r.CloseConnection {
		req.SetConnectionClose()
	}

	// Response
	resp := fasthttp.AcquireResponse()

	// Send
	s.log(
		logging.LEVEL_DEBUG, "HTTP Request: '%s %s' ...", r.Method, req.URI().Path(),
	)
	err := s.app.Do(req, resp)

	// Release
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		return err
	}

	return f(resp)
}

func (s *Session) RequestJSON(req RequestConfig, obj interface{}) error {
	req.Accept = "application/json"
	return s.Request(req, func(r *fasthttp.Response) error {
		body := r.Body()
		code := r.StatusCode()

		if code != 200 {
			err := &StatusCodeError{Code: code}
			// unneccerry to check nil around the range
			for _, v := range body {
				err.Msg = string(v)
			}

			return err
		}

		return json.Unmarshal(body, &obj)
	})
}
