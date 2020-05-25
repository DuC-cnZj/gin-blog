package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/youngduc/go-blog/hello/config"
	"github.com/youngduc/go-blog/hello/models"
	"github.com/youngduc/go-blog/hello/models/dao"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

var oauthCnf *oauth2.Config

func Init() {
	oauthCnf = config.Config.Oauth
}

func Redirect(c *gin.Context) {
	url := oauthCnf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	c.Redirect(http.StatusMovedPermanently, url)
}

func HandleProviderCallback(code string) {
	//ctx := context.Background()
	//tok, err := oauthCnf.Exchange(ctx, code)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//client := oauthCnf.Client(ctx, tok)
	//resp, err := client.Get("https://api.github.com/user?access_token=" + tok.AccessToken)
	//if err != nil {
	//	log.Println("err", err)
	//	return
	//}
	var user models.GithubUser
	var err error
	//err = json.NewDecoder(resp.Body).Decode(&user)
	err = json.NewDecoder(strings.NewReader(`
{
  "login": "DuC-cnZj",
  "id": 23514869,
  "node_id": "MDQ6VXNlcjIzNTE0ODY5",
  "avatar_url": "https://avatars0.githubusercontent.com/u/23514869?v=4",
  "gravatar_id": "",
  "url": "https://api.github.com/users/DuC-cnZj",
  "html_url": "https://github.com/DuC-cnZj",
  "followers_url": "https://api.github.com/users/DuC-cnZj/followers",
  "following_url": "https://api.github.com/users/DuC-cnZj/following{/other_user}",
  "gists_url": "https://api.github.com/users/DuC-cnZj/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/DuC-cnZj/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/DuC-cnZj/subscriptions",
  "organizations_url": "https://api.github.com/users/DuC-cnZj/orgs",
  "repos_url": "https://api.github.com/users/DuC-cnZj/repos",
  "events_url": "https://api.github.com/users/DuC-cnZj/events{/privacy}",
  "received_events_url": "https://api.github.com/users/DuC-cnZj/received_events",
  "type": "User",
  "site_admin": false,
  "name": "duc",
  "company": null,
  "blog": "",
  "location": "HangZhou, china",
  "email": "1025434218@qq.com",
  "hireable": null,
  "bio": "https://whoops-cn.club",
  "public_repos": 25,
  "public_gists": 1,
  "followers": 8,
  "following": 13,
  "created_at": "2016-11-17T03:15:33Z",
  "updated_at": "2020-05-24T13:02:40Z",
  "private_gists": 1,
  "total_private_repos": 5,
  "owned_private_repos": 5,
  "disk_usage": 46986,
  "collaborators": 1,
  "two_factor_authentication": false,
  "plan": {
    "name": "free",
    "space": 976562499,
    "collaborators": 0,
    "private_repos": 10000
  }
}
`)).Decode(&user)
	log.Println(user)
	if err != nil {
		return
	}

	dao.Dao.SaveSocialiteUser(user)
}
