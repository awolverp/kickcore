package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Logging level's
//
// DEBUG < INFO < WARN[ING] < ERROR < CRITICAL
const (
	LEVEL_CRITICAL = iota
	LEVEL_ERROR
	LEVEL_WARNING
	LEVEL_INFO
	LEVEL_DEBUG

	LEVEL_WARN = LEVEL_WARNING
)

func GetLevelName(level int) string {
	switch level {
	case LEVEL_DEBUG:
		return "DEBUG"
	case LEVEL_INFO:
		return "INFO"
	case LEVEL_WARNING:
		return "WARNING"
	case LEVEL_ERROR:
		return "ERROR"
	case LEVEL_CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

type FileLogger struct {
	level    int
	filename string
	f        io.Writer
	locker   sync.Mutex
}

type Config struct {
	// Specifies that should be using file.
	Filename string

	// File object to write. if spectified, filename and append are ignored.
	FileObject io.Writer

	// Specifies that not truncrate the file.
	Append bool
}

func NewLogger(level int, c *Config) (*FileLogger, error) {
	if c == nil {
		l := new(FileLogger)
		l.level = level
		l.f = os.Stdout
		return l, nil
	}

	var thefile io.Writer

	if c.FileObject != nil {
		thefile = c.FileObject
	} else if c.Filename != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if c.Append {
			flags = flags | os.O_APPEND
		} else {
			flags = flags | os.O_TRUNC
		}

		f, err := os.OpenFile(c.Filename, flags, 0666)
		if err != nil {
			return nil, err
		}
		runtime.SetFinalizer(f, (*os.File).Close)

		thefile = f
	} else {
		thefile = os.Stdout
	}

	l := new(FileLogger)
	l.filename = c.Filename
	l.level = level
	l.f = thefile

	return l, nil
}

func MustNewLogger(level int, c *Config) *FileLogger {
	if c == nil {
		l := new(FileLogger)
		l.level = level
		l.f = os.Stdout
		return l
	}

	var thefile io.Writer

	if c.FileObject != nil {
		thefile = c.FileObject
	} else if c.Filename != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if c.Append {
			flags = flags | os.O_APPEND
		} else {
			flags = flags | os.O_TRUNC
		}

		f, err := os.OpenFile(c.Filename, flags, 0666)
		if err != nil {
			panic(err)
		}
		runtime.SetFinalizer(f, (*os.File).Close)

		thefile = f
	} else {
		thefile = os.Stdout
	}

	l := new(FileLogger)
	l.filename = c.Filename
	l.level = level
	l.f = thefile

	return l
}

func (n *FileLogger) Log(level int, msg string, args ...interface{}) int {
	var count int = 0

	issues_time := time.Now().Format("15:04:05")

	n.locker.Lock()
	if level <= n.level {
		var prefix string = issues_time

		if n.level == LEVEL_DEBUG {
			_, file, line, ok := runtime.Caller(1)
			if ok {
				prefix += " " + filepath.Base(file) + ":" + strconv.Itoa(line)
			}
		}

		if msg[len(msg)-1] != '\n' {
			msg += "\n"
		}

		count, _ = fmt.Fprintf(n.f, prefix+" ["+GetLevelName(level)+"] "+msg, args...)
	}
	n.locker.Unlock()

	return count
}

func (n *FileLogger) Fatal(msg string, args ...interface{}) {
	n.Log(LEVEL_CRITICAL, msg, args...)
	os.Exit(1)
}

func (n *FileLogger) Level() int { n.locker.Lock(); defer n.locker.Unlock(); return n.level }

func (n *FileLogger) SetLevel(level int) { n.locker.Lock(); n.level = level; n.locker.Unlock() }

func (n *FileLogger) String() string {
	n.locker.Lock()
	defer n.locker.Unlock()
	return "<FileLogger file=" + n.filename + " level=" + strconv.Itoa(n.level) + ">"
}
