<h1 align=center>
    KickCore
</h1>

<p align=center>
    Is a football API Server (language: Persian)
</p>

<p align=center>
    <a href="./LICENSE.md">Copyright</a> -
    <a href="./CHANGELOG.md">Changelog</a>
</p>

--------

- **Content**
  - [**How It Works?**](#how-it-works)
  - [**Install**](#install)
  - [**API**](#api)
    - [**Search**](#search)
    - [**Advanced Search**](#adavnced-search)
    - [**Competitions List**](#list-of-competitions)
    - [**Competition Weeks**](#weeks-of-competition)
    - [**Standing Table**](#competition-standing-table)
    - [**Competition Matches**](#competition-matches-by-week)
    - [**Match Info**](#match-info)
    - [**Matches**](#matches)
    - [**Transfers Regions**](#transfers-regions)
    - [**Transfers**](#transfers)
    - [**Memory Stats**](#memory-stats-developer-api)
  - [**What is** `extra_ttl.json` **file?**](#how-to-write-expire-ttl-file)

## How It Works?
```
             |      |---------|                     |---------|
Request ---> | ---> |         | --- Not Exists ---> |   API   |
             |      |         |                     |---------|
             |      |  Cache  |                          |
             |      |         |                        Result
             |      |         | <------------------------|
             |      |---------|                          |
             |           |                               |
             |         Exists                            |
             |           |                               |
             |  <--------|-------------------------------|
```

## Install
### Build From Source
**Requirements**
- **Go** (version 1.7 or above)

```
go install github.com/awolverp/kickcore@latest
```

## API
KickCore API Documentation.

> **Note**: {url} is the host that kickcore uses e.g. 'http://127.0.0.1:9090'

### Search
Search.

```bash
curl "{url}/api/search"
```

**Query Params**
|  Key  | Value  | Description |
| ----- | ------ | ----------- |
|   q   | string | Search Query (q length must be > 4). |

### Adavnced Search
Search (you can filter result).

```bash
curl "{url}/api/search/advanced"
```

**Query Params**
|  Key   | Value   | Description |
| -----  | ------  | ----------- |
|   q    | string  | Search Query (q length must be > 4). |
| filter | integer | Filter result. zero means teams, 1 means players, 2 means coaches, 3 means competitions |
| limit  | integer | Optional. Result limit  |
| offset | integer | Optional. Result offset |

### List of competitions
Get list of competitions which are supported.

```bash
curl "{url}/api/competitions-list"
```

### Weeks of competition
Get weeks of a competition.

```bash
curl "{url}/api/competition/weeks"
```

### Competition Standing Table
Get standing table of a competition.

```bash
curl "{url}/api/competition/standing-table"
```

**Query Params**
|  Key  | Value  | Description |
| ----- | ------ | ----------- |
|  id   | string | Current ID of the competition |

### Competition Matches by week
Get standing table of a competition.

```bash
curl "{url}/api/competition/matches/week"
```

**Query Params**
|  Key  | Value   | Description |
| ----- | ------- | ----------- |
|  id   | string  | Current ID of the competition |
|  n    | integer | Week number |

### Match info
Get match information.

```bash
curl "{url}/api/match/info"
```

**Query Params**
|  Key  | Value   | Description |
| ----- | ------- | ----------- |
|  id   | string  | Match ID |

### Matches
Get matches by date.

```bash
curl "{url}/api/matches"
```

**Query Params**
|  Key  | Value   | Description |
| ----- | ------- | ----------- |
| days  | integer | Optional. Zero is today. 1 is tomorrow, 2 two days later, etc. (and you can pass nagative numbers). |

### Transfers Regions
Get regions (and seasons) which have transfers.

```bash
curl "{url}/api/transfers/regions"
```

### Transfers
Get transfers of a season.

```bash
curl "{url}/api/transfers"
```

**Query Params**
|  Key  | Value  | Description |
| ----- | ------ | ----------- |
| sid   | string | seasion ID |

### Memory Stats (Developer API)
Get memory usage/stats of script.

```bash
curl "{url}/stats/mem"
```

**Query Params**
|  Key  | Value  | Description |
| ----- | ------ | ----------- |
| unit  | string | Optional. Is the unit byte (b or byte, kb or kilobyte, mb or megabyte). default is byte. |

-----

## Questions

### How to write expire ttl file?
> **What is extra_ttl.json file?**

You can specify the expiration time (TTL) of each api
with expire ttl file.

**How to write?** Each api has key that you can use it to
specify expiration time.

#### Keys
- Advanced Search: `ADVANCED_SEARCH`
- Competition Standing Table: `COMPETITION_STANDING_TABLE`
- Competition Weeks: `COMEPTITION_WEEKS`
- List of competitions: `COMPETITIONS_LIST`
- Match info: `MATCH_INFO`
- Matches: `MATCHES_BY_DATE`
- Competitions match by week: `MATCHES_BY_WEEKNUMBER`
- Transfers: `TRANSFERS`
- Transfers Regions: `TRANSFERS_REGIONS`

> Other keys will ignored

**What value can be set?** Integer (means seconds) or 
string (duration, see `extra_ttl.json` file for examples)
