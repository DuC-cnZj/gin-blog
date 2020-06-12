package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/olivere/elastic/v6"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"log"
	"sync"
	"time"
)

var (
	Config *config
	once   sync.Once
)

type es struct {
	SearchQuery string
	Highlight   *elastic.Highlight
	MultiMatch  multiMatch
	Host        string
}

type config struct {
	App         *app
	DB          *db
	ES          es
	Oauth       *oauth2.Config
	Redis       *redis.Options
	RedisPrefix string
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
	once.Do(func() {
		Config = &config{
			DB:          InitDB(),
			App:         InitApp(),
			ES:          initES(),
			Oauth:       initOauth(),
			Redis:       InitRedis(),
			RedisPrefix: viper.GetString("REDIS_PREFIX"),
		}
	})

	return Config
}

func initOauth() *oauth2.Config {
	redirectURL := viper.GetString("OAUTH_REDIRECT_URL")
	clientID := viper.GetString("OAUTH_CLIENT_ID")
	clientSecret := viper.GetString("OAUTH_CLIENT_SECRET")

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: redirectURL,
		Scopes:      []string{"user", "repo"},
	}
}

func InitRedis() *redis.Options {
	return &redis.Options{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST"), viper.GetString("REDIS_PORT")),
		DB:   viper.GetInt("REDIS_DB"),
	}
}

func initES() es {
	var es es
	es.Host = viper.GetString("ES_HOST")
	log.Println("es hosts:", es.Host)
	es.SearchQuery = esSearchConfig()
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
		Debug:        viper.GetBool("DEBUG"),
		Domain:       viper.GetString("Domain"),
		RunMode:      viper.GetString("RUN_MODE"),
		PageSize:     viper.GetInt("PAGE_SIZE"),
		JwtSecret:    viper.GetString("JWT_SECRET"),
		HttpPort:     viper.GetInt("HTTP_PORT"),
		FrontDomain:  viper.GetString("FRONT_DOMAIN"),
		ReadTimeout:  time.Duration(viper.GetInt64("READ_TIMEOUT")),
		WriteTimeout: time.Duration(viper.GetInt64("WRITE_TIMEOUT")),
	}
}

func esSearchConfig() string {
	return `
{
  "query" : {
    "multi_match": {
      "query": "%s",
      "fields": ["content", "title^2", "desc^2", "tags^3", "article_category.name^4", "author.name^5"],
      "analyzer": "ik_smart"
    }
  },
  "fields":{
    "title":{
      "type":"plain",
      "pre_tags":"<span style='background-color:#bfa;padding:1px;'>",
      "post_tags":"</span>"
    },
    "tags":{
      "type":"plain",
      "pre_tags":"<span style='background-color:#bfa;padding:1px;'>",
      "post_tags":"</span>"
    },
    "article_category.name":{
      "type":"plain",
      "pre_tags":"<span style='background-color:#bfa;padding:1px;'>",
      "post_tags":"</span>"
    },
    "content":{
      "type":"plain",
      "pre_tags":"<span style='background-color:#bfa;padding:1px;'>",
      "post_tags":"</span>",
      "fragment_size":10,
      "number_of_fragments":2
    },
    "desc":{
      "type":"plain",
      "fragment_size":10,
      "number_of_fragments":2
    }
  },
  "pre_tags":"<span style='color:red'>",
  "post_tags":"</span>"
}

`
}

type app struct {
	Debug       bool
	Domain      string
	RunMode     string
	PageSize    int
	JwtSecret   string
	FrontDomain string

	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
