package main

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
