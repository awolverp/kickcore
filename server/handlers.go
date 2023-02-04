package server

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/awolverp/kickcore/api"
	"github.com/awolverp/kickcore/cache"

	"github.com/valyala/fasthttp"
)

// Memory usage: cat /proc/pid/smaps | grep -i pss |  awk '{Total+=$2} END {print Total/1024" MB"}'

var URLs = [][2]interface{}{
	// Redirects to /doc/
	{"/", indexPage},

	// Memory usage
	{"/stats/mem", memoryUsage}, // unit

	{"/api/search", searchAPI},                  // q
	{"/api/search/advanced", advancedSearchAPI}, // q, filter, offset, limit

	{"/api/competitions-list", getCompetitionsList},                  // type
	{"/api/competition/weeks", getCompetitionWeeks},                  // id
	{"/api/competition/standing-table", getCompetitionStandingTable}, // id
	{"/api/competition/matches/week", getMatchesByWeekNumber},        // id, n

	{"/api/match/info", getMatchInfo},  // id
	{"/api/matches", getMatchesByDate}, // days, slugs

	{"/api/transfers/regions", getTransfersRegions}, // -
	{"/api/transfers", getTransfers},                // sid
}

func indexPage(ctx *fasthttp.RequestCtx, _ *api.Session, _ *cache.Cache) error {
	ctx.Redirect("https://github.com/awolverp/kickcore/", 301)
	return nil
}

func memoryUsage(ctx *fasthttp.RequestCtx, _ *api.Session, _ *cache.Cache) error {
	var data_unit string = "B" // Bytes

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "unit", Optional: true, Object: &data_unit},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	datamap := map[string]interface{}{
		"alloc":       toUnit(m.Alloc, &data_unit),
		"total_alloc": toUnit(m.TotalAlloc, &data_unit),
		"sys":         toUnit(m.Sys, &data_unit),
		"num_gc":      m.NumGC,
	}

	datamap["all"] = datamap["alloc"].(uint64) + datamap["sys"].(uint64)
	datamap["code"] = 200
	datamap["unit"] = data_unit

	data, _ := api.ToBytes(datamap)

	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(200)
	ctx.Write(data)
	return nil
}

func advancedSearchAPI(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		query  string
		filter int
		offset uint16
		limit  uint16
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "q", Optional: false, Object: &query},
			{Name: "filter", Optional: false, Object: &filter},
			{Name: "offset", Optional: true, Object: &offset},
			{Name: "limit", Optional: true, Object: &limit},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	if len(query) < 4 {
		e := api.StatusCodeError{Code: 400, Msg: "query is too short: len(q) < 4"}
		ret, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetBody(ret)
		return nil
	}

	if filter > 3 {
		e := api.StatusCodeError{Code: 400, Msg: "filter must be between 0-3"}
		ret, _ := e.MarshalJSON()
		ctx.SetStatusCode(e.Code)
		ctx.SetBody(ret)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.ADVANCED_SEARCH,
		cache.GenerateKey(
			query,
			strconv.Itoa(filter),
			strconv.Itoa(int(offset)),
			strconv.Itoa(int(limit)),
		),
		func() (interface{}, error) {
			return cli.AdvancedSearch(query, uint8(filter), offset, limit)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getCompetitionStandingTable(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		current_id string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "id", Optional: false, Object: &current_id},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.COMPETITION_STANDING_TABLE,
		cache.GenerateKey(
			current_id,
		),
		func() (interface{}, error) {
			return cli.GetCompetitionStandingTable(current_id)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getCompetitionWeeks(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		current_id string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "id", Optional: false, Object: &current_id},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.COMPETITION_WEEKS,
		cache.GenerateKey(
			current_id,
		),
		func() (interface{}, error) {
			return cli.GetCompetitionWeeks(current_id)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getCompetitionsList(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		c_type string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "type", Optional: true, Object: &c_type},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.COMPETITIONS_LIST,
		cache.GenerateKey(
			c_type,
		),
		func() (interface{}, error) {
			return cli.GetCompetitionsList(c_type)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getMatchInfo(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		match_id string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "id", Optional: false, Object: &match_id},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.MATCH_INFO,
		cache.GenerateKey(
			match_id,
		),
		func() (interface{}, error) {
			return cli.GetMatchInfo(match_id)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getMatchesByDate(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		days    int
		slugs_q string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "days", Optional: true, Object: &days},
			{Name: "slugs", Optional: true, Object: &slugs_q},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	slugs := strings.Split(slugs_q, ",")

	data, _, err := cacheObject.CacheFuncJSON(
		cache.MATCHES_BY_DATE,
		cache.GenerateKey(
			strconv.Itoa(days), slugs_q,
		),
		func() (interface{}, error) {
			return cli.GetMatchesByDate(days, slugs...)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getMatchesByWeekNumber(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		weeknumber int
		id         string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "id", Optional: false, Object: &id},
			{Name: "n", Optional: false, Object: &weeknumber},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.MATCHES_BY_WEEKNUMBER,
		cache.GenerateKey(
			id, strconv.Itoa(weeknumber),
		),
		func() (interface{}, error) {
			return cli.GetMatchesByWeekNumber(id, uint32(weeknumber))
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getTransfers(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		sid string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "sid", Optional: false, Object: &sid},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.TRANSFERS,
		cache.GenerateKey(
			sid,
		),
		func() (interface{}, error) {
			return cli.GetTransfers(sid)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func getTransfersRegions(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	data, _, err := cacheObject.CacheFuncJSON(
		cache.TRANSFERS_REGIONS,
		"",
		func() (interface{}, error) {
			return cli.GetTransfersRegions()
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}

func searchAPI(ctx *fasthttp.RequestCtx, cli *api.Session, cacheObject *cache.Cache) error {
	ctx.SetContentType("application/json; charset=utf-8")

	var (
		q string
	)

	err := queryArgsParser(
		ctx.QueryArgs(),
		[]queryConfig{
			{Name: "q", Optional: false, Object: &q},
		},
	)
	if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
		return nil
	}

	data, _, err := cacheObject.CacheFuncJSON(
		cache.SEARCH,
		cache.GenerateKey(q),
		func() (interface{}, error) {
			return cli.Search(q, 0)
		},
	)

	if data != nil {
		ctx.SetStatusCode(200)
		ctx.SetBody(data)
	} else if err != nil {
		i, b, _ := api.ErrToBytes(err)
		ctx.SetStatusCode(i)
		ctx.SetBody(b)
	}

	return err
}
