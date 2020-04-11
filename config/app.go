package config

import (
	"encoding/json"
	"github.com/olivere/elastic/v6"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"time"
)

var Config *config

type es struct {
	SearchQuery string
	Highlight   *elastic.Highlight
	MultiMatch  multiMatch
}

type config struct {
	App *app
	DB  *db
	ES  es
}

type multiMatch struct {
	Query struct {
		MultiMatch struct {
			Query    string   `json:"query"`
			Fields   []string `json:"fields"`
			Analyzer string   `json:"analyzer"`
		} `json:"multi_match"`
	} `json:"query"`
	Fields map[string]struct {
		HType             string `json:"type"`
		PreTags           string `json:"pre_tags"`
		PostTags          string `json:"post_tags"`
		FragmentSize      int    `json:"fragment_size"`
		NumberOfFragments int    `json:"number_of_fragments"`
	}
	PreTags  string `json:"pre_tags"`
	PostTags string `json:"post_tags"`
}

func Init() *config {
	Config = &config{
		DB:  InitDB(),
		App: InitApp(),
		ES:  initES(),
	}

	return Config
}

func initES() es {
	var es es
	bytes, e := ioutil.ReadFile("./config/es_search.json")
	if e != nil {
		log.Fatal(e)
	}
	es.SearchQuery = string(bytes)
	var multiMatch multiMatch
	json.Unmarshal([]byte(es.SearchQuery), &multiMatch)
	var Fields []*elastic.HighlighterField
	for name, v := range multiMatch.Fields {
		var f = elastic.HighlighterField{
			Name: name,
		}
		f.PostTags(v.PostTags).
			PreTags(v.PreTags).
			FragmentSize(v.FragmentSize).
			NumOfFragments(v.NumberOfFragments)
		Fields = append(Fields, &f)
	}
	es.MultiMatch = multiMatch
	es.Highlight = elastic.
		NewHighlight().
		Fields(Fields...).
		PreTags(multiMatch.PreTags).
		PostTags(multiMatch.PostTags)
	return es
}

func InitApp() *app {
	return &app{
		RunMode:      viper.GetString("RUN_MODE"),
		PageSize:     viper.GetInt("PAGE_SIZE"),
		JwtSecret:    viper.GetString("JWT_SECRET"),
		HttpPort:     viper.GetInt("HTTP_PORT"),
		ReadTimeout:  time.Duration(viper.GetInt64("READ_TIMEOUT")),
		WriteTimeout: time.Duration(viper.GetInt64("WRITE_TIMEOUT")),
	}
}

type app struct {
	RunMode   string
	PageSize  int
	JwtSecret string

	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
