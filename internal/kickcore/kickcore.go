package kickcore

import (
	"context"
	"errors"
	"kickcore/api"
	"kickcore/cache"
	"kickcore/cache/noncache"
	"kickcore/cache/sqlite"
	"kickcore/logging"
	"kickcore/server"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

// Logging level's
//
// DEBUG < INFO < WARN[ING] < ERROR < CRITICAL
const (
	LOGGING_CRITICAL = logging.LEVEL_CRITICAL
	LOGGING_ERROR    = logging.LEVEL_ERROR
	LOGGING_WARNING  = logging.LEVEL_WARNING
	LOGGING_INFO     = logging.LEVEL_INFO
	LOGGING_DEBUG    = logging.LEVEL_DEBUG

	LOGGING_WARN = logging.LEVEL_WARN
)

type Core struct {
	// Logging system
	logger *logging.FileLogger

	// API Client
	api_client *api.Session

	// Cache system
	cache_struct *cache.Cache
	expirator    *cache.ExpirationMachine

	server_app fasthttp.Server
	server_mux server.ServeMux
}

type ConfigCore struct {
	LoggingLevel  int
	LoggingConfig *logging.Config

	APIClientReadTimeout  time.Duration
	APIClientWriteTimeout time.Duration

	// CacheSystem    string
	DisableCaching                 bool
	CacheSQLiteTimeout             time.Duration
	CacheExpirationMachineInterval time.Duration
	CacheExtraTTLFilename          string
	CacheSQLiteDSN                 string

	ServerReadTimeout       time.Duration
	ServerWriteTimeout      time.Duration
	ReduceServerMemoryUsage bool
	ServerGetOnly           bool

	ServerLogSpeed bool
}

func (core *Core) Init(c *ConfigCore) error {
	if c == nil {
		return errors.New("argument (*ConfigCore) is nil")
	}

	var err error

	core.logger, err = logging.NewLogger(c.LoggingLevel, c.LoggingConfig)
	if err != nil {
		return err
	}

	core.api_client = api.NewSession(core.logger, c.APIClientReadTimeout, c.APIClientWriteTimeout)

	if c.DisableCaching {
		core.cache_struct, _ = cache.NewCache(noncache.Connect())
	} else {
		if c.CacheSQLiteDSN == "" {
			c.CacheSQLiteDSN = "db.sqlite3"
		}

		core.cache_struct, err = cache.NewCache(sqlite.Connect(c.CacheSQLiteDSN, c.CacheSQLiteTimeout))
	}
	if err != nil {
		return err
	}

	core.cache_struct.Logger = core.logger

	if c.CacheExtraTTLFilename != "" {
		err = cache.ReadExtraTTL(c.CacheExtraTTLFilename)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				core.logger.Log(LOGGING_WARNING, "Extra TTL file not found: '%s'!", c.CacheExtraTTLFilename)
				err = os.WriteFile(c.CacheExtraTTLFilename, defaultExtraTTL, 0655)
				if err != nil {
					return err
				}

				core.logger.Log(LOGGING_WARNING, "Extra TTL file created with default content: '%s'", c.CacheExtraTTLFilename)
				err = cache.ReadExtraTTL(c.CacheExtraTTLFilename)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	var disable_cache_expiration bool = (c.CacheExpirationMachineInterval <= 0)

	if !disable_cache_expiration && !c.DisableCaching {
		core.expirator = &cache.ExpirationMachine{
			ActiveCache: core.cache_struct,
		}

		err = core.expirator.Start(c.CacheExpirationMachineInterval)
		if err != nil {
			return err
		}
	}

	core.server_app = fasthttp.Server{
		Name:              "awolverp",
		ReadTimeout:       c.ServerReadTimeout,
		WriteTimeout:      c.ServerWriteTimeout,
		ReduceMemoryUsage: c.ReduceServerMemoryUsage,
		GetOnly:           c.ServerGetOnly,
		CloseOnShutdown:   true,
		DisableKeepalive:  true,
	}

	core.server_mux = server.ServeMux{
		APIClient: core.api_client,
		Cache:     core.cache_struct,
		Logger:    core.logger,
		LogSpeed:  c.ServerLogSpeed,
	}
	core.server_mux.Init()

	return nil
}

func (core *Core) Urls() [][2]interface{} { return server.URLs }

func (core *Core) OpenConnections() int32 { return core.server_app.GetOpenConnectionsCount() }

func (core *Core) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if core.expirator != nil {
		core.expirator.Stop()
	}
	return core.server_app.ShutdownWithContext(ctx)
}

func (core *Core) Serve(addr string) error {
	return server.Serve(addr, &core.server_app, &core.server_mux)
}

var (
	version string = "2.4.9"
)

func Version() string { return "v" + version }

var defaultExtraTTL = []byte(`{
	"ADVANCED_SEARCH":            "24h",
	"COMPETITION_STANDING_TABLE": "5h",
	"COMPETITION_WEEKS":          "1h",
	"COMPETITIONS_LIST":          "24h",
	"MATCH_INFO":                 "1m",
	"MATCHES_BY_DATE":            "2m",
	"MATCHES_BY_WEEKNUMBER":      "2m",
	"TRANSFERS":                  "12h",
	"TRANSFERS_REGIONS":          "24h",
	"SEARCH":                     "24h"
}`)
