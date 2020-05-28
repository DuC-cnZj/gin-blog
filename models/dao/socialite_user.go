package dao

import (
	"github.com/youngduc/go-blog/models"
	"log"
	"strconv"
	"time"
)

func (dao *dao) SaveSocialiteUser(githubUser *models.GithubUser) *models.SocialiteUser  {
	log.Println(githubUser)
	var su models.SocialiteUser
	identifier :=strconv.Itoa(githubUser.ID)
	dao.DB.
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

	dao.DB.Save(&su)

	return &su
}

func (dao *dao) FindSocialiteUser(id int) *models.SocialiteUser {
	var su models.SocialiteUser
	dao.DB.
		Where("id = ?", id).
		Find(&su)

	return &su
}