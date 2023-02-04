# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.4.10] - 2023-2-4
### Fixed
- Fixed go.mod error

## [2.4.9] - 2023-2-4
### Fixed
- Project main application

## [2.4.7] - 2023-2-3
### Fixed
- Fixed nil pointer bug on shutdown

## [2.4.6] - 2023-2-2
### Fixed
- 404 error bug fixed.

## [2.4.5] - 2023-1-30
### Fixed
- Fixed help message text.

### Changed
- More command-line options are renamed:
  - `-addr` to `-l`
  - `-server-read-timeout` to `-server-timeout:read`
  - `-server-write-timeout` to `-server-timeout:write`
  - `-expire-interval` to `-expire:interval`
  - `-expire-ttl` to `-expire:ttl`
  - `-sqlite-dsn` to `-sqlite:dsn`
  - `-sqlite-timeout` to `-sqlite:timeout`
  - `-client-read-timeout` to `-client-timeout:read`
  - `-client-write-timeout` to `-client-timeout:write`
  - `-verbose` to `-v`
  - `-log-file` to `-log:file`
  - `-log-append` to `-log:append`
  - `-log-speed` to `-log:speed`
  - `-V` to `-version`

## [2.3.4] - 2023-1-30
### Changed
- `-expire-ttl` is removed.
- `-extra-ttl` renamed to `-expire-ttl`.

### Fixed
- Fixed some error messages in `-extra-ttl`.

## [2.2.3] - 2023-1-29
### Added
- Now you can get memory status

### Fixed
- Fixed signal's bug.

## [2.1.2] - 2023-1-27
### Fixed
- Fixed help message problem.

## [2.1.1] - 2023-1-27
### Added
- Now can set extra ttl for each object by `-extra-ttl` command-line option.

### Changed
- Increase debugging logs in cache expiration machine.
- `-cache-expire-ttl` changed to `-expire-ttl`
- `-cache-expire-interval` changed to `-expire-interval`
- `-log-verbose` changed to `-verbose`
- `-log-filename` changed to `-log-file`
- `-version` changed to `-V`

## [2.0.1] - 2023-1-27
### Fixed
- Fixed SQLite cache expiration bugs.

## [2.0.0] - 2023-1-26
### Added
- Cache now can be disable.

### Changed
- Code structure changed.
- Server library changed from net/http to fasthttp to make server performance 10 times better than before.
- Cache strcuture changed to make develop easier.
