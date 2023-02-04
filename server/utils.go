package server

import (
	"strconv"
	"strings"

	"github.com/awolverp/kickcore/api"

	"github.com/valyala/fasthttp"
)

type queryConfig struct {
	Name     string
	Optional bool
	Object   interface{}
}

func queryArgsParser(args *fasthttp.Args, q []queryConfig) error {
	for _, value := range q {
		obj := args.Peek(value.Name)

		if obj == nil {
			if !value.Optional {
				return &api.StatusCodeError{Code: 400, Msg: "'" + value.Name + "' parameter required."}
			}
			continue
		}

		switch valueObject := value.Object.(type) {
		case (*[]byte):
			(*valueObject) = obj

		case *int:
			i, err := strconv.Atoi(string(obj))
			if err != nil {
				(*valueObject) = 0
			} else {
				(*valueObject) = i
			}

		case *uint8:
			i, err := strconv.Atoi(string(obj))
			if err != nil {
				(*valueObject) = 0
			} else {
				(*valueObject) = uint8(i)
			}

		case *uint16:
			i, err := strconv.Atoi(string(obj))
			if err != nil {
				(*valueObject) = 0
			} else {
				(*valueObject) = uint16(i)
			}

		case *uint32:
			i, err := strconv.ParseUint(string(obj), 10, 0)
			if err != nil {
				(*valueObject) = 0
			} else {
				(*valueObject) = uint32(i)
			}

		case *string:
			(*valueObject) = string(obj)
		}
	}

	return nil
}

func toUnit(i uint64, unit *string) uint64 {
	*unit = strings.ToLower(*unit)

	switch *unit {
	case "kb", "kilobytes":
		i = i / 1024

	case "mb", "megabytes":
		i = i / 1024 / 1024

	default:
		*unit = "b"
	}

	return i
}
