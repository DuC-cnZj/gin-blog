package dao

import (
	"github.com/youngduc/go-blog/hello/models"
	"log"
	"strconv"
	"time"
)

func (dao *dao) SaveSocialiteUser(githubUser models.GithubUser)  {
	log.Println(githubUser)
	var su models.SocialiteUser
	identifier :=strconv.Itoa(githubUser.ID)
	dao.db.
		Where("identity_type = ?", "github").
		Where("identifier = ?", identifier).
		Find(&su)
	
	if su.Id != 0 {
		su.LastLoginAt = models.JSONTime{
			Time: time.Now(),
		}
	} else  {
		su.Identifier = identifier
		su.IdentityType = "github"
		su.Avatar =  &githubUser.AvatarURL
		su.Url =  &githubUser.URL
		su.Name =  &githubUser.Name
		su.LastLoginAt = models.JSONTime{
			Time: time.Now(),
		}
	}

	dao.db.Save(&su)
}