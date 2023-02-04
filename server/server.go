package server

import (
	"fmt"
	"kickcore/api"
	"kickcore/cache"
	"kickcore/logging"
	"net"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type Handler func(ctx *fasthttp.RequestCtx, cli *api.Session, c *cache.Cache) error

type ServeMux struct {
	APIClient *api.Session
	Cache     *cache.Cache

	// Logger
	Logger *logging.FileLogger

	// ErrorHandler
	ErrorHandler func(ctx *fasthttp.RequestCtx, err error)

	// Logs functions ping speed
	LogSpeed bool

	handlers map[string]Handler
}

func (m *ServeMux) callHandler(ctx *fasthttp.RequestCtx, h Handler) {
	var pingspeed time.Time
	var fname string

	if m.LogSpeed {
		pingspeed = time.Now()
	}

	ctx.SetConnectionClose()
	err := h(ctx, m.APIClient, m.Cache)

	if m.LogSpeed && m.Logger != nil {
		stop := time.Since(pingspeed)

		fobj := runtime.FuncForPC(reflect.ValueOf(h).Pointer())
		if fobj != nil {
			fname = strings.Split(fobj.Name(), ".")[1]
		}

		m.Logger.Log(logging.LEVEL_INFO, "(%s speed): %v", fname, stop)
	}

	if err != nil {

		if fname == "" {
			fobj := runtime.FuncForPC(reflect.ValueOf(h).Pointer())

			if fobj != nil {
				fname = strings.Split(fobj.Name(), ".")[1]
			}
		}

		if m.Logger != nil {
			m.Logger.Log(logging.LEVEL_ERROR, "(%s): %s", fname, err.Error())
		} else {
			fmt.Printf("ERROR: (%s): %s", fname, err.Error())
		}

		if m.ErrorHandler != nil {
			m.ErrorHandler(ctx, err)
		} else {
			ctx.Error(`{"code":500,"message":"internal server error"}`, fasthttp.StatusInternalServerError)
		}
	}
}

func (m *ServeMux) Init() {
	if m.handlers == nil {
		m.handlers = make(map[string]Handler)
	}

	for _, v := range URLs {
		m.handlers[v[0].(string)] = v[1].(func(*fasthttp.RequestCtx, *api.Session, *cache.Cache) error)
	}
}

func (m *ServeMux) AddHandler(path string, f Handler) {
	if m.handlers == nil {
		m.handlers = make(map[string]Handler)
	}

	m.handlers[path] = f
}

func (m *ServeMux) HandleHTTP(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	callback, ok := m.handlers[path]

	if !ok {
		if m.ErrorHandler == nil {
			ctx.Error(`{"code":404,"message":"Page Not Found"}`, 404)
		} else {
			m.ErrorHandler(ctx, nil)
		}
		return
	}

	if m.Logger != nil {
		m.Logger.Log(
			logging.LEVEL_INFO, "%s - \"%s %s\"", ctx.RemoteIP().String(), string(ctx.Method()), path,
		)
	}

	m.callHandler(ctx, callback)
}

func Serve(addr string, s *fasthttp.Server, mux *ServeMux) error {
	if addr == "" {
		addr = "127.0.0.1:8070"

		if mux.Logger != nil {
			mux.Logger.Log(
				logging.LEVEL_INFO, "Listening address is empty! It changed to '%s' automatcally.", addr,
			)
		}
	}

	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	s.Handler = mux.HandleHTTP
	s.ErrorHandler = nil
	mux.Init()

	fmt.Printf(
		"+ SERVER RUNNING http://%s/ ...\n", addr,
	)

	return s.Serve(l)
}

func ServeTLS(addr, certFile, keyFile string, s *fasthttp.Server, mux *ServeMux) error {
	if addr == "" {
		addr = "127.0.0.1:8070"

		if mux.Logger != nil {
			mux.Logger.Log(
				logging.LEVEL_INFO, "Listening address is empty! It changed to '%s' automatcally.", addr,
			)
		}
	}

	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	s.Handler = mux.HandleHTTP
	s.ErrorHandler = nil
	mux.Init()

	fmt.Printf(
		"+ SERVER RUNNING https://%s/ ...\n", addr,
	)

	return s.ServeTLS(l, certFile, keyFile)
}
