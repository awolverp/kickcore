package main

import (
	"flag"
	"fmt"
	"kickcore"
	"kickcore/logging"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"time"
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
