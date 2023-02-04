package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var hostAddr = string([]byte{0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x66, 0x6f, 0x6f, 0x74, 0x62, 0x61, 0x6c, 0x6c, 0x33, 0x36, 0x30, 0x2e, 0x69, 0x72})

// Searches with result filtering
//
// Parameters:
//   - q: Query ( assert len(q) > 3 )
//   - filter: Search for what? teams:0, players:1, coaches:2, or competitions:3
//   - offset: result offset
//   - limit: result limit ( default 10 )
func (cli *Session) AdvancedSearch(q string, filter uint8, offset, limit uint16) (AdvancedSuggestsInterface, error) {
	if len(q) < 4 {
		return nil, errors.New("query is too short: len(q) < 4")
	}
	q = url.PathEscape(q)

	if limit == 0 {
		limit = 10
	}

	var (
		ret      AdvancedSuggestsInterface
		filter_q string
	)

	switch filter {
	case 0:
		v := TeamSuggests{offset: offset, limit: limit}
		filter_q = "teams"
		ret = &v
	case 1:
		v := PlayerSuggests{offset: offset, limit: limit}
		filter_q = "players"
		ret = &v
	case 2:
		v := CoacheSuggests{offset: offset, limit: limit}
		filter_q = "coaches"
		ret = &v
	case 3:
		v := CompetitionSuggests{offset: offset, limit: limit}
		filter_q = "competitions"
		ret = &v
	default:
		v := TeamSuggests{offset: offset, limit: limit}
		filter_q = "teams"
		ret = &v
		fmt.Printf("(*Client).AdvancedSearch: Warning: unknown filter: %d, It is setting 0 automatically.\n", filter)
	}

	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + fmt.Sprintf("/api/search/%s/?q=%s&offset=%d&limit=%d", filter_q, q, offset, limit),
			Referer: hostAddr + "/search/", CloseConnection: true,
		},
		&ret,
	)
	if err != nil {
		return nil, err
	}

	return ret, err
}

// Returns standing table of a competition.
//
// Parameters:
//   - current_id: competition current id.
func (cli *Session) GetCompetitionStandingTable(current_id string) (StandingTable, error) {
	if current_id == "" {
		return nil, errors.New("current_id is empty")
	}

	var obj StandingTable
	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/standing-table/" + current_id + "/",
			Referer: hostAddr, CloseConnection: true,
		},
		&obj,
	)
	return obj, err
}

// Returns competition weeks.
//
// Parameters:
//   - current_id: competition current id.
func (cli *Session) GetCompetitionWeeks(current_id string) (*CompetitionWeeks, error) {
	if current_id == "" {
		return nil, errors.New("current_id is empty")
	}

	var obj CompetitionWeeks
	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/competition-trends/" + current_id + "/",
			Referer: hostAddr, CloseConnection: true,
		},
		&obj,
	)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

// Returns the competitions that are supported.
//
// Parameters:
//   - c_type: competition type, (C)lub competitions or (N)ational competitions ( can be zero ).
func (cli *Session) GetCompetitionsList(c_type string) (CompetitionsList, error) {
	var c_type_q string
	if c_type != "" {
		c_type_q = "&type=" + c_type
	}

	var rawobj map[string]json.RawMessage

	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/base/competitions/defaults/" + c_type_q,
			Referer: hostAddr, CloseConnection: true,
		},
		&rawobj,
	)
	if err != nil {
		return nil, err
	}

	var obj CompetitionsList = CompetitionsList{}
	if v, ok := rawobj["competitions"]; ok {
		err = json.Unmarshal(v, &obj)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Returns the match information.
//
// Parameters:
//   - match_id: the match id.
func (cli *Session) GetMatchInfo(match_id string) (*MatchInfo, error) {
	var obj MatchInfo
	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/base/v2/matches/" + match_id + "/info/",
			Referer: hostAddr, CloseConnection: true,
		},
		&obj,
	)
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// Returns the matches of that date.
//
// Parameters:
//   - days: days after today ( can be zero or negative ).
//   - slugs: the slugs of competitions that matches you want.
func (cli *Session) GetMatchesByDate(days int, slugs ...string) (CompetitionMatches, error) {
	var slug_q string
	if len(slugs) > 0 {
		slug_q = "&slugs=" + url.QueryEscape(strings.Join(slugs, ","))
	}

	date := time.Now()
	if days != 0 {
		date = date.AddDate(0, 0, days)
	}

	var obj CompetitionMatches = nil
	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/competition-trends/matches-by-date/?date=" + date.Format("2006-01-02") + slug_q,
			Referer: hostAddr, CloseConnection: true,
		},
		&obj,
	)

	return obj, err
}

// Returns the matches of that week.
//
// Parameters:
//   - week_number: week number.
//   - current_id: competition current id.
func (cli *Session) GetMatchesByWeekNumber(current_id string, week_number uint32) ([]MatchBase, error) {
	if current_id == "" {
		return nil, errors.New("current_id is empty")
	}

	var obj []MatchBase = nil
	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/competition-trends/" + current_id + "/weeks/" + strconv.Itoa(int(week_number)) + "/",
			Referer: hostAddr, CloseConnection: true,
		},
		&obj,
	)

	return obj, err
}

// Returns the transfers of a season.
//
// Parameters:
//   - season_id: season id.
func (cli *Session) GetTransfers(season_id string) (Transfers, error) {
	if season_id == "" {
		return nil, errors.New("season_id is empty")
	}

	var rawobj map[string]json.RawMessage

	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/transfers/transfer-seasons/" + season_id + "/transfers/",
			Referer: hostAddr, CloseConnection: true,
		},
		&rawobj,
	)
	if err != nil {
		return nil, err
	}

	var obj Transfers = Transfers{}
	if v, ok := rawobj["data"]; ok {
		err = json.Unmarshal(v, &obj)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Returns the regions that have transfer.
func (cli *Session) GetTransfersRegions() (*TransfersRegions, error) {
	var obj TransfersRegions

	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/transfers/regions/",
			Referer: hostAddr, CloseConnection: true,
		},
		&obj,
	)
	if err != nil {
		return nil, err
	}

	return &obj, err
}

// Searches the query.
//
// Parameters:
//   - q: query ( must be len(q) > 3).
//   - s_type: search type - Full-Search:1 or Simple-Search:0
//
// - Simple search: coaches, players, teams.
// - Full search: coaches, players, teams, competitions, news.
func (cli *Session) Search(q string, s_type uint8) (*Suggests, error) {
	if len(q) < 4 {
		return nil, errors.New("query is too short: len(q) < 4")
	}

	q = url.PathEscape(q)

	var s_type_q string
	switch s_type {
	case 1:
		s_type_q = "page"
	case 0:
		s_type_q = "bar"
	default:
		s_type_q = "bar"
		println("(*Client).Search: Warning: unknown 's_type'. It is choose Simple-Search automatically.")
	}

	var obj Suggests
	err := cli.RequestJSON(
		RequestConfig{
			Method: "GET", URI: hostAddr + "/api/search/suggest/?q=" + q + "&location=" + s_type_q,
			Referer: hostAddr, CloseConnection: true,
		},
		&obj,
	)
	if err != nil {
		return nil, err
	}

	return &obj, err
}
