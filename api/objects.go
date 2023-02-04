package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// compile-time type checks
var (
	_ error = (*StatusCodeError)(nil)
)

type StatusCodeError struct {
	Msg  string `json:"message"`
	Code int    `json:"code"`
}

func (e *StatusCodeError) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("status code [%d]", e.Code)
	}

	msg := e.Msg
	if strings.Contains(msg, "\n") {
		msg = strings.ReplaceAll(msg, "\n", "  ")

		if len(msg) > 30 {
			msg = msg[0:27] + "..."
		}
	}
	return fmt.Sprintf("status code [%d]: %s", e.Code, msg)
}

func escapeString(s string) string { return strings.Replace(s, `"`, `\"`, -1) }

func (e *StatusCodeError) MarshalJSON() ([]byte, error) {
	var msg string

	if e.Msg != "" {
		msg = e.Msg
		if strings.Contains(msg, "\n") {
			msg = strings.ReplaceAll(msg, "\n", "  ")

			if len(msg) > 30 {
				msg = msg[0:27] + "..."
			}
		}
	}

	return []byte(fmt.Sprintf(`{"code":%d,"message":"%s"}`, e.Code, escapeString(msg))), nil
}

func IsStatusCodeError(target error) bool { _, ok := target.(*StatusCodeError); return ok }

func ErrToBytes(obj interface{}) (int, []byte, error) {
	if err, ok := obj.(error); ok {
		if statusCode, ok := obj.(*StatusCodeError); ok {
			data, err := statusCode.MarshalJSON()
			return statusCode.Code, data, err
		}

		if statusCode, ok := obj.(json.Marshaler); ok {
			data, err := statusCode.MarshalJSON()
			return 500, data, err
		}

		data, err := (&StatusCodeError{Msg: err.Error(), Code: 500}).MarshalJSON()
		return 500, data, err
	}

	data, err := json.Marshal(obj)
	return 500, data, err
}

func ToBytes(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

type Team struct {
	ID    string `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Logo  string `json:"logo"`
	// Thumbnail  string `json:"thumbnail"`
	// IsActive   bool   `json:"is_active"`
	// FullTitle  string `json:"full_title"`
	// IsNational bool   `json:"is_national"`
	Country *struct {
		Name string `json:"name"`
	} `json:"country"`
	ToBeDecided bool `json:"to_be_decided"`
}

type Stadium struct {
	Name    string `json:"name"`
	Country struct {
		Name string `json:"name"`
	} `json:"country"`
	City string `json:"city"`
	// Latitude  float64     `json:"latitude"`
	// Longitude float64     `json:"longitude"`
	Capacity int `json:"capacity"`
	// OpenedAt  interface{} `json:"opened_at"`
}

type AdvancedSuggestsInterface interface {
	Offset() uint16
	Limit() uint16

	/*
		returns ((s.offset + s.limit) < uint16(s.Count))

		To see how much more exists: (uint16(s.Count) - (s.offset + s.limit))
	*/
	HasMore() bool
}

type TeamSuggests struct {
	Count   int    `json:"count"`
	Results []Team `json:"results"`

	offset, limit uint16
}

func (s TeamSuggests) Offset() uint16 { return s.offset }
func (s TeamSuggests) Limit() uint16  { return s.limit }
func (s TeamSuggests) HasMore() bool  { return ((s.offset + s.limit) < uint16(s.Count)) }

type PlayerSuggests struct {
	Count   int `json:"count"`
	Results []struct {
		ID string `json:"id"`
		// Slug     string `json:"slug"`
		Fullname string `json:"fullname"`
		Image    string `json:"image"`
		// Position struct {
		// 	Key   string `json:"key"`
		// 	Value string `json:"value"`
		// } `json:"position"`
		KitNumber int `json:"kit_number"`
		// Person    struct {
		// 	ID       string `json:"id"`
		// 	Fullname string `json:"fullname"`
		// 	Image    string `json:"image"`
		// } `json:"person"`
	} `json:"results"`

	offset, limit uint16
}

func (s PlayerSuggests) Offset() uint16 { return s.offset }
func (s PlayerSuggests) Limit() uint16  { return s.limit }
func (s PlayerSuggests) HasMore() bool  { return ((s.offset + s.limit) < uint16(s.Count)) }

type CoacheSuggests struct {
	Count   int `json:"count"`
	Results []struct {
		ID       string `json:"id"`
		Fullname string `json:"fullname"`
		Person   struct {
			// Fullname string `json:"fullname"`
			// ID       string `json:"id"`
			Image string `json:"image"`
		} `json:"person"`
		// Slug string `json:"slug"`
	} `json:"results"`

	offset, limit uint16
}

func (s CoacheSuggests) Offset() uint16 { return s.offset }
func (s CoacheSuggests) Limit() uint16  { return s.limit }
func (s CoacheSuggests) HasMore() bool  { return ((s.offset + s.limit) < uint16(s.Count)) }

type CompetitionSuggests struct {
	Count   int `json:"count"`
	Results []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Slug      string `json:"slug"`
		Logo      string `json:"logo"`
		Thumbnail string `json:"thumbnail"`
		SeoSlug   string `json:"seo_slug"`
	} `json:"results"`

	offset, limit uint16
}

func (s CompetitionSuggests) Offset() uint16 { return s.offset }
func (s CompetitionSuggests) Limit() uint16  { return s.limit }
func (s CompetitionSuggests) HasMore() bool  { return ((s.offset + s.limit) < uint16(s.Count)) }

type StandingTable []struct {
	Team struct {
		// ID    string `json:"id"`
		// Slug  string `json:"slug"`
		Title string `json:"title"`
		// Logo       string `json:"logo"`
		// Thumbnail  string `json:"thumbnail"`
		// IsActive   bool   `json:"is_active"`
		// FullTitle  string `json:"full_title"`
		// IsNational bool   `json:"is_national"`
		// Country struct {
		// 	Name string `json:"name"`
		// } `json:"country"`
		ToBeDecided bool `json:"to_be_decided"`
	} `json:"team"`
	// Form []struct {
	// 	State   string `json:"state"`
	// 	MatchID string `json:"match_id"`
	// 	Title   string `json:"title"`
	// } `json:"form"`
	Rank           int `json:"rank"`
	Score          int `json:"score"`
	PlayedMatches  int `json:"played_matches"`
	WonMatches     int `json:"won_matches"`
	LostMatches    int `json:"lost_matches"`
	ScoredGoals    int `json:"scored_goals"`
	ConcededGoals  int `json:"conceded_goals"`
	RedCards       int `json:"red_cards"`
	YellowCards    int `json:"yellow_cards"`
	GoalDifference int `json:"goal_difference"`
	TotalCards     int `json:"total_cards"`
	RankChange     int `json:"rank_change"`
}

type CompetitionWeeks struct {
	ID          string `json:"id"`
	Competition string `json:"competition"`
	// StartTime   interface{} `json:"start_time"`
	// SeoSlug string `json:"seo_slug"`
	Slug string `json:"slug"`
	// EndTime     int         `json:"end_time"`
	CurrentWeek struct {
		WeekNumber uint32 `json:"week_number"`
	} `json:"current_week"`

	Weeks []struct {
		WeekNumber uint32 `json:"week_number"`
	} `json:"weeks"`
}

type CompetitionsList []struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Current struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Slug  string `json:"slug"`
		Logo  string `json:"logo"`
		// Thumbnail string `json:"thumbnail"`
		SeoSlug string `json:"seo_slug"`
	} `json:"current"`
}

type MatchInfo struct {
	ID string `json:"id"`

	HomeTeam *Team `json:"home_team"`
	AwayTeam *Team `json:"away_team"`

	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`

	HoldsAt int `json:"holds_at"`
	// IsPostponed bool   `json:"is_postponed"`
	IsFinished bool `json:"is_finished"`
	Status     struct {
		StatusID   int    `json:"status_id"`
		Title      string `json:"title"`
		StatusType string `json:"status_type"`
	} `json:"status"`

	Minute int `json:"minute"`
	// Slug             string `json:"slug"`

	HomePenaltyScore int `json:"home_penalty_score"`
	AwayPenaltyScore int `json:"away_penalty_score"`

	// RoundType        struct {
	// 	Name        string `json:"name"`
	// 	Value       int    `json:"value"`
	// 	IsKnockout  bool   `json:"is_knockout"`
	// 	DisplayName string `json:"display_name"`
	// } `json:"round_type"`
	// Spectators            int    `json:"spectators"`

	ToBeDecided bool `json:"to_be_decided"`

	BroadcastChannel string `json:"broadcast_channel"`

	CompetitionTrendStage struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		// StageType string      `json:"stage_type"`
		// StartTime int         `json:"start_time"`
		// EndTime   int         `json:"end_time"`
		// Order     interface{} `json:"order"`
		// IsDefault bool        `json:"is_default"`
	} `json:"competition_trend_stage"`

	CompetitionTrend struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Slug  string `json:"slug"`
		// Logo      string `json:"logo"`
		// Thumbnail string `json:"thumbnail"`
		// SeoSlug   string `json:"seo_slug"`
	} `json:"competition_trend"`

	Stadium *Stadium `json:"stadium"`

	// HasStats       bool          `json:"has_stats"`
	// HasLineups     bool          `json:"has_lineups"`
	// RelatedMatches []interface{} `json:"related_matches"`

	HeaderEvents []struct {
		ID     string `json:"id"`
		Player struct {
			ID string `json:"id"`
			// Slug     string `json:"slug"`
			Fullname string `json:"fullname"`
			// Image    string `json:"image"`
			// Position struct {
			// 	Key   string `json:"key"`
			// 	Value string `json:"value"`
			// } `json:"position"`
			KitNumber int `json:"kit_number"`
			// Person    struct {
			// 	ID       string `json:"id"`
			// 	Fullname string `json:"fullname"`
			// 	Image    string `json:"image"`
			// } `json:"person"`
		} `json:"player"`

		Team      *Team `json:"team"`
		EventType struct {
			ShortForm string `json:"short_form"`
			Title     string `json:"title"`
			// EmptyValue interface{} `json:"empty_value"`
		} `json:"event_type"`

		Minute     int `json:"minute"`
		MinutePlus int `json:"minute_plus"`
		// Period     string `json:"period"`
		// Reason             interface{} `json:"reason"`
		// SortOrder          interface{} `json:"sort_order"`
		// MatchEventRelation struct {
		// 	Relation string `json:"relation"`
		// 	Type     string `json:"type"`
		// } `json:"match_event_relation"`
	} `json:"header_events"`

	// StateTimelines []struct {
	// 	State string `json:"state"`
	// 	Time  int    `json:"time"`
	// } `json:"state_timelines"`
	// HasRelatedPost bool `json:"has_related_post"`
}

type MatchBase struct {
	ID            string `json:"id"`
	HomeTeam      *Team  `json:"home_team"`
	AwayTeam      *Team  `json:"away_team"`
	HomeScore     int    `json:"home_score"`
	AwayScore     int    `json:"away_score"`
	StatusDetails struct {
		StatusID   int    `json:"status_id"`
		Title      string `json:"title"`
		StatusType string `json:"status_type"`
	} `json:"status_details"`
	HoldsAt   int `json:"holds_at"`
	StartedAt int `json:"started_at"`
	// IsActive         bool        `json:"is_active"`
	// IsPostponed      bool        `json:"is_postponed"`
	BroadcastChannel string `json:"broadcast_channel"`
	IsFinished       bool   `json:"is_finished"`
	Competition      struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		// Slug      string `json:"slug"`
		// Logo      string `json:"logo"`
		// Thumbnail string `json:"thumbnail"`
		// SeoSlug   string `json:"seo_slug"`
	} `json:"competition"`
	WeekNumber int `json:"week_number"`
	Minute     int `json:"minute"`
	// Slug       string `json:"slug"`
	Stadium          *Stadium `json:"stadium"`
	HomePenaltyScore int      `json:"home_penalty_score"`
	AwayPenaltyScore int      `json:"away_penalty_score"`
	// Spectators       int `json:"spectators"`
	HasStanding bool `json:"has_standing"`
	// HasStats                bool          `json:"has_stats"`
	// HasLineups              bool          `json:"has_lineups"`
	CompetitionTrendStageID string `json:"competition_trend_stage_id"`
	// RoundType               struct {
	// 	Name        string `json:"name"`
	// 	Value       int    `json:"value"`
	// 	IsKnockout  bool   `json:"is_knockout"`
	// 	DisplayName string `json:"display_name"`
	// } `json:"round_type"`
	// StateTimelines []interface{} `json:"state_timelines"`
}

type CompetitionMatches []struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Slug  string `json:"slug"`
	// Logo       string `json:"logo"`
	// Thumbnail  string `json:"thumbnail"`
	// SeoSlug    string `json:"seo_slug"`
	// IsInFilter bool   `json:"is_in_filter"`
	Matches []MatchBase `json:"matches"`
}

type Transfers []struct {
	ID string `json:"id"`
	// Slug       string `json:"slug"`
	Title string `json:"title"`
	Logo  string `json:"logo"`
	// Thumbnail  string `json:"thumbnail"`
	// IsActive   bool   `json:"is_active"`
	// FullTitle  string `json:"full_title"`
	// IsNational bool   `json:"is_national"`
	Country struct {
		Name string `json:"name"`
	} `json:"country"`
	ToBeDecided bool `json:"to_be_decided"`
	// Order       int  `json:"order"`
	InTransfers []struct {
		TransferSeason string `json:"transfer_season"`
		Player         struct {
			ID       string `json:"id"`
			Slug     string `json:"slug"`
			Fullname string `json:"fullname"`
			// Image    string `json:"image"`
			// Position struct {
			// 	Key   string `json:"key"`
			// 	Value string `json:"value"`
			// } `json:"position"`
			KitNumber int `json:"kit_number"`
			// Person    struct {
			// 	ID       string `json:"id"`
			// 	Fullname string `json:"fullname"`
			// 	Image    string `json:"image"`
			// } `json:"person"`
		} `json:"player"`
		FromTeam       Team `json:"from_team"`
		TransferTime   int  `json:"transfer_time"`
		TransferStatus struct {
			StatusID int    `json:"status_id"`
			Name     string `json:"name"`
		} `json:"transfer_status"`
		IsImportant bool `json:"is_important"`
	} `json:"in_transfers"`
	OutTransfers []struct {
		TransferSeason string `json:"transfer_season"`
		Player         struct {
			ID       string `json:"id"`
			Slug     string `json:"slug"`
			Fullname string `json:"fullname"`
			// Image    string `json:"image"`
			// Position struct {
			// 	Key   string `json:"key"`
			// 	Value string `json:"value"`
			// } `json:"position"`
			KitNumber int `json:"kit_number"`
			// Person    struct {
			// 	ID       string `json:"id"`
			// 	Fullname string `json:"fullname"`
			// 	Image    string `json:"image"`
			// } `json:"person"`
		} `json:"player"`
		ToTeam         *Team `json:"to_team"`
		TransferTime   int   `json:"transfer_time"`
		TransferStatus struct {
			StatusID int    `json:"status_id"`
			Name     string `json:"name"`
		} `json:"transfer_status"`
		IsImportant bool `json:"is_important"`
	} `json:"out_transfers"`
}

type TransfersRegions struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		// Order   int    `json:"order"`
		Seasons []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			// Order int    `json:"order"`
		} `json:"seasons"`
	} `json:"results"`
}

type Suggests struct {
	Teams   TeamSuggests   `json:"teams"`
	Players PlayerSuggests `json:"players"`
	Coaches CoacheSuggests `json:"coaches"`
	// News struct {
	// 	Count   int `json:"count"`
	// 	Results []struct {
	// 		Code         int         `json:"code"`
	// 		Title        string      `json:"title"`
	// 		SubTitle     string      `json:"sub_title"`
	// 		SuperTitle   interface{} `json:"super_title"`
	// 		PrimaryMedia string      `json:"primary_media"`
	// 		CreatedAt    int         `json:"created_at"`
	// 		PublishedAt  int         `json:"published_at"`
	// 		Author       struct {
	// 			ID        string      `json:"id"`
	// 			FullName  string      `json:"full_name"`
	// 			AvatarID  int         `json:"avatar_id"`
	// 			Image     interface{} `json:"image"`
	// 			Thumbnail interface{} `json:"thumbnail"`
	// 		} `json:"author"`
	// 		Medias []struct {
	// 			Media struct {
	// 				ID         string      `json:"id"`
	// 				File       string      `json:"file"`
	// 				Thumbnail  string      `json:"thumbnail"`
	// 				MediaType  string      `json:"media_type"`
	// 				Title      string      `json:"title"`
	// 				AparatLink interface{} `json:"aparat_link"`
	// 				CoverImage interface{} `json:"cover_image"`
	// 				ArvanLink  interface{} `json:"arvan_link"`
	// 				ArvanHls   interface{} `json:"arvan_hls"`
	// 				ArvanMpd   interface{} `json:"arvan_mpd"`
	// 				ArvanAd    interface{} `json:"arvan_ad"`
	// 			} `json:"media"`
	// 			IsPrimary bool `json:"is_primary"`
	// 		} `json:"medias"`
	// 		IsPublished  bool   `json:"is_published"`
	// 		Slug         string `json:"slug"`
	// 		Link         string `json:"link"`
	// 		ID           string `json:"id"`
	// 		ViewCount    int    `json:"view_count"`
	// 		HitCount     int    `json:"hit_count"`
	// 		CommentCount int    `json:"comment_count"`
	// 		Tags         []struct {
	// 			ID        string `json:"id"`
	// 			Title     string `json:"title"`
	// 			IsVisible bool   `json:"is_visible"`
	// 			Instance  struct {
	// 				InstanceData struct {
	// 					ID         string `json:"id"`
	// 					Slug       string `json:"slug"`
	// 					Title      string `json:"title"`
	// 					Logo       string `json:"logo"`
	// 					Thumbnail  string `json:"thumbnail"`
	// 					IsActive   bool   `json:"is_active"`
	// 					FullTitle  string `json:"full_title"`
	// 					IsNational bool   `json:"is_national"`
	// 					Country    struct {
	// 						Name    string `json:"name"`
	// 						Flag1X1 string `json:"flag_1x1"`
	// 						Flag4X3 string `json:"flag_4x3"`
	// 					} `json:"country"`
	// 					ToBeDecided bool `json:"to_be_decided"`
	// 				} `json:"instance_data"`
	// 				InstanceType string `json:"instance_type"`
	// 			} `json:"instance"`
	// 		} `json:"tags"`
	// 		PostType string `json:"post_type"`
	// 	} `json:"results"`
	// } `json:"news"`
	// Videos struct {
	// 	Count   int `json:"count"`
	// 	Results []struct {
	// 		Code         int         `json:"code"`
	// 		Title        string      `json:"title"`
	// 		SubTitle     string      `json:"sub_title"`
	// 		SuperTitle   interface{} `json:"super_title"`
	// 		PrimaryMedia string      `json:"primary_media"`
	// 		CreatedAt    int         `json:"created_at"`
	// 		PublishedAt  int         `json:"published_at"`
	// 		Author       struct {
	// 			ID        string      `json:"id"`
	// 			FullName  string      `json:"full_name"`
	// 			AvatarID  int         `json:"avatar_id"`
	// 			Image     interface{} `json:"image"`
	// 			Thumbnail interface{} `json:"thumbnail"`
	// 		} `json:"author"`
	// 		Medias []struct {
	// 			Media struct {
	// 				ID         string      `json:"id"`
	// 				File       string      `json:"file"`
	// 				Thumbnail  string      `json:"thumbnail"`
	// 				MediaType  string      `json:"media_type"`
	// 				Title      string      `json:"title"`
	// 				AparatLink interface{} `json:"aparat_link"`
	// 				CoverImage interface{} `json:"cover_image"`
	// 				ArvanLink  interface{} `json:"arvan_link"`
	// 				ArvanHls   interface{} `json:"arvan_hls"`
	// 				ArvanMpd   interface{} `json:"arvan_mpd"`
	// 				ArvanAd    interface{} `json:"arvan_ad"`
	// 			} `json:"media"`
	// 			IsPrimary bool `json:"is_primary"`
	// 		} `json:"medias"`
	// 		IsPublished  bool   `json:"is_published"`
	// 		Slug         string `json:"slug"`
	// 		Link         string `json:"link"`
	// 		ID           string `json:"id"`
	// 		ViewCount    int    `json:"view_count"`
	// 		HitCount     int    `json:"hit_count"`
	// 		CommentCount int    `json:"comment_count"`
	// 		Tags         []struct {
	// 			ID        string `json:"id"`
	// 			Title     string `json:"title"`
	// 			IsVisible bool   `json:"is_visible"`
	// 			Instance  struct {
	// 				InstanceData struct {
	// 					ID         string `json:"id"`
	// 					Slug       string `json:"slug"`
	// 					Title      string `json:"title"`
	// 					Logo       string `json:"logo"`
	// 					Thumbnail  string `json:"thumbnail"`
	// 					IsActive   bool   `json:"is_active"`
	// 					FullTitle  string `json:"full_title"`
	// 					IsNational bool   `json:"is_national"`
	// 					Country    struct {
	// 						Name    string `json:"name"`
	// 						Flag1X1 string `json:"flag_1x1"`
	// 						Flag4X3 string `json:"flag_4x3"`
	// 					} `json:"country"`
	// 					ToBeDecided bool `json:"to_be_decided"`
	// 				} `json:"instance_data"`
	// 				InstanceType string `json:"instance_type"`
	// 			} `json:"instance"`
	// 		} `json:"tags"`
	// 		PostType string `json:"post_type"`
	// 	} `json:"results"`
	// } `json:"videos"`
	// Competitions CompetitionSuggests `json:"competitions"`
}
