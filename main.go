package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/awolverp/kickcore/internal/kickcore"
	"github.com/awolverp/kickcore/logging"
)

var (
	logConfig   logging.Config
	coreConfig  kickcore.ConfigCore
	ListenAddr  string
	showVersion bool
	showUrls    bool
)

var core kickcore.Core

func main() {
	sigchannel := make(chan os.Signal, 1)
	donechannel := make(chan struct{})

	go kickcore_server(donechannel)

	signal.Notify(sigchannel, os.Interrupt)

	select {
	case <-donechannel:
		signal.Stop(sigchannel)
		close(donechannel)
		os.Exit(0)

	case <-sigchannel:
		fmt.Printf("Please wait (for 5 seconds), don't try again (open connections %d) ...\n", core.OpenConnections())
		core.Shutdown()

		signal.Stop(sigchannel)
		close(sigchannel)

		close(donechannel)

		os.Exit(0)
	}
}

func kickcore_server(done chan<- struct{}) {
	flag.Usage = func() { fmt.Printf(helpUsage, kickcore.Version(), os.Args[0]) }

	// server
	flag.StringVar(&ListenAddr, "l", "127.0.0.1:9090", "")
	flag.DurationVar(&coreConfig.ServerReadTimeout, "server-timeout:read", time.Second*30, "")
	flag.DurationVar(&coreConfig.ServerWriteTimeout, "server-timeout:write", time.Second*30, "")
	flag.BoolVar(&coreConfig.ReduceServerMemoryUsage, "reduce-memory-usage", false, "")
	flag.BoolVar(&coreConfig.ServerGetOnly, "get-only", false, "")

	// cache
	flag.BoolVar(&coreConfig.DisableCaching, "disable-cache", false, "")
	flag.DurationVar(&coreConfig.CacheExpirationMachineInterval, "expire:interval", time.Minute, "")
	flag.StringVar(&coreConfig.CacheExtraTTLFilename, "expire:ttl", "extra_ttl.json", "")
	flag.StringVar(&coreConfig.CacheSQLiteDSN, "sqlite:dsn", "db.sqlite3", "")
	flag.DurationVar(&coreConfig.CacheSQLiteTimeout, "sqlite:timeout", time.Minute, "")

	// api client timeouts
	flag.DurationVar(&coreConfig.APIClientReadTimeout, "client-timeout:read", time.Second*20, "")
	flag.DurationVar(&coreConfig.APIClientWriteTimeout, "client-timeout:write", time.Second*20, "")

	// logging options
	flag.IntVar(&coreConfig.LoggingLevel, "v", kickcore.LOGGING_WARNING, "")
	flag.StringVar(&logConfig.Filename, "log:file", "", "")
	flag.BoolVar(&logConfig.Append, "log:append", false, "")
	flag.BoolVar(&coreConfig.ServerLogSpeed, "log:speed", false, "")

	// other
	flag.BoolVar(&showUrls, "urls", false, "")
	flag.BoolVar(&showVersion, "version", false, "")

	flag.Parse()

	if showVersion {
		fmt.Printf(
			"kickcore %s [%s] - build on %s/%s\n",
			kickcore.Version(), runtime.Version(),
			runtime.GOOS, runtime.GOARCH,
		)
		done <- struct{}{}
		return
	}

	if showUrls {
		for _, value := range core.Urls() {
			var fname string
			fobj := runtime.FuncForPC(reflect.ValueOf(value[1]).Pointer())

			if fobj != nil {
				fname = strings.Split(fobj.Name(), ".")[1]
			}

			fmt.Printf(
				"%s - func %s\n", value[0].(string), fname,
			)
		}

		done <- struct{}{}
		return
	}

	fmt.Printf(
		"KickIt Core Server %s (C) / by aWolverP - [%s] on %s/%s\n\n",
		kickcore.Version(), runtime.Version(), runtime.GOOS, runtime.GOARCH,
	)

	coreConfig.LoggingConfig = &logConfig

	if err := core.Init(&coreConfig); err != nil {
		fmt.Println("ERROR", err)
		done <- struct{}{}
		return
	}

	if err := core.Serve(ListenAddr); err != nil {
		fmt.Println("ERROR", err)
	}
	done <- struct{}{}
}

var helpUsage = `NAME
       kickcore %s - KickCore Server (C)

USAGE
       %s [OPTIONS]

DESCRIPTION
       kickcore (C) is a Football API server written in golang language.

OPTIONS
  *Server
      -l=address     (default "127.0.0.1:9090")
            Server listening address. If the port in the address
            parameter is empty or "0", as in "127.0.0.1:" or
            "[::1]:0", a port number is automatically chosen. The
            Addr method of Listener can be used to discover the
            chosen port.
        
      -server-timeout:read=duration     (default 30s)
            is the amount of time allowed to read the full request
            including body. The connection's read deadline is
            reset when the connection opens, or for keep-alive
            connections after the first byte has been read.

      -server-timeout:write=duration     (default 30s)
            is the maximum duration before timing out writes of the
            response. It is reset after the request handler has
            returned.
        
      -reduce-memory-usage
            reduces memory usage at the cost of higher CPU usage.
            Try enabling this option only if the server consumes too
            much memory serving mostly idle keep-alive connections.
            This may reduce memory usage by more than 50%%.

      -get-only
            Rejects all non-GET requests. This option is useful as
            anti-DoS protection for servers accepting only GET
            requests.

  *Cache
      -disable-cache
            Disable cache. It slows down this server and maybe banned
            from original football API.
        
      -expire:interval=duration     (default 1m)
            The Cache expiration machine checks the cache for expired
            objects after any interval time.
        
      -expire:ttl=filename     (default "extra_ttl.json")
            Configuration file of Time-To-Live of cached objects.
            file format must be JSON, like 'extra_ttl.json'.
            if set empty, all objects are deleted after every -expire-interval.

      -sqlite:dsn=dsn     (default "db.sqlite3")
            SQLite path address.
        
      -sqlite:timeout=duration     (default 1m)
            SQLite connecting timeout.

  *API Client
      -client-timeout:read=duration     (default 20s)
            Maximum duration for full response reading (including body)
        
      -client-timeout:write=duration     (default 20s)
            Maximum duration for full request writing (including body).

  *Logging
      -v=[0-4]     (default 1)
            Logging verbose level.

              0   Critical level.
              1   Error level.
              2   Warning level.
              3   Information level.
              4   Debugging level.
        
      -log:file=filename     (default "")
            Logging filename. if set, the logs are written in the file.
        
      -log:append
            If set, not truncate file and append new logs to file.
        
      -log:speed
            Show server handlers ping speed. (needs -v 3 or 4)

  *Other
      -version  Print version and exit.

      -urls  Print URLs and exit.
`
