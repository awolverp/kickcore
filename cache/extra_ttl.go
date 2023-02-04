package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func GenerateKey(s ...string) string {
	return strings.Join(s, "-")
}

type APICacheKey struct {
	Key      string
	ExtraTTL int64
}

// Implemented for kickcore/api package
var (
	EMPTY_APIKEY APICacheKey = APICacheKey{Key: "", ExtraTTL: 0}

	ADVANCED_SEARCH APICacheKey = APICacheKey{
		Key: "0", ExtraTTL: 0,
	}

	COMPETITION_STANDING_TABLE APICacheKey = APICacheKey{
		Key: "1", ExtraTTL: 0,
	}

	COMPETITION_WEEKS APICacheKey = APICacheKey{
		Key: "2", ExtraTTL: 0,
	}

	COMPETITIONS_LIST APICacheKey = APICacheKey{
		Key: "3", ExtraTTL: 0,
	}

	MATCH_INFO APICacheKey = APICacheKey{
		Key: "4", ExtraTTL: 0,
	}

	MATCHES_BY_DATE APICacheKey = APICacheKey{
		Key: "5", ExtraTTL: 0,
	}

	MATCHES_BY_WEEKNUMBER APICacheKey = APICacheKey{
		Key: "6", ExtraTTL: 0,
	}

	TRANSFERS APICacheKey = APICacheKey{
		Key: "7", ExtraTTL: 0,
	}

	TRANSFERS_REGIONS APICacheKey = APICacheKey{
		Key: "8", ExtraTTL: 0,
	}

	SEARCH APICacheKey = APICacheKey{
		Key: "9", ExtraTTL: 0,
	}
)

var mapVars = map[string](*APICacheKey){
	"ADVANCED_SEARCH":            &ADVANCED_SEARCH,
	"COMPETITION_STANDING_TABLE": &COMPETITION_STANDING_TABLE,
	"COMPETITION_WEEKS":          &COMPETITION_WEEKS,
	"COMPETITIONS_LIST":          &COMPETITIONS_LIST,
	"MATCH_INFO":                 &MATCH_INFO,
	"MATCHES_BY_DATE":            &MATCHES_BY_DATE,
	"MATCHES_BY_WEEKNUMBER":      &MATCHES_BY_WEEKNUMBER,
	"TRANSFERS":                  &TRANSFERS,
	"TRANSFERS_REGIONS":          &TRANSFERS_REGIONS,
	"SEARCH":                     &SEARCH,
}

func ReadExtraTTL(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var content map[string]interface{}

	err = json.Unmarshal(data, &content)
	if err != nil {
		return err
	}

	for key, value := range mapVars {
		obj, ok := content[key]
		if !ok {
			continue
		}

		var d int64

		switch objvalue := obj.(type) {
		case int:
			d = int64(objvalue)

		case string:
			var dur time.Duration
			dur, err = time.ParseDuration(objvalue)
			if err == nil {
				d = int64(dur.Seconds())
			}
		}

		if err != nil {
			return fmt.Errorf("expire ttl: cannot parse '%v' ('%s'): %s", obj, key, err.Error())
		}

		(*value).ExtraTTL = d
	}

	return nil
}
