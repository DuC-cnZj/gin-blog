package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type content struct {
	Md   string `json:"md"`
	Html string `json:"html"`
}

func (dao *dao) IndexArticles(page, perPage int) map[string]interface{} {
	var articles []models.Article
	var count int
	offset := (page - 1) * perPage

	dao.DB.
		Preload("Author").
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("id desc").
		Offset(offset).
		Limit(perPage).
		Find(&articles)

	dao.DB.Table("articles").Where("display = ?", true).Count(&count)

	return map[string]interface{}{
		"data": articles,
		"meta": models.Paginator{
			Total:       count,
			CurrentPage: page,
			PerPage:     perPage,
		},
		"links": map[string]string{},
	}
}

func cacheKey(key string) string {
	return fmt.Sprintf("%s%s", config.Config.RedisPrefix, key)
}

func (dao *dao) GetArticleByIds(ids []int) []models.Article {
	var articles []models.Article
	dao.DB.Select("id,head_image,title,`desc`,display,created_at").
		Where("id in (?)", ids).
		Find(&articles)

	return articles
}

func (dao *dao) ShowArticle(id int) (*models.Article, BaseError) {
	key := cacheKey("article:" + strconv.Itoa(id))
	s, e := dao.Redis.Get(key).Result()
	if e == redis.Nil {
		article := &models.Article{}
		e := dao.DB.
			Preload("Author").
			Preload("Category").
			Preload("Tags").
			Where("id = ?", id).
			Where("display = ?", true).
			Find(article).
			Error

		if e != nil && e == gorm.ErrRecordNotFound {
			return nil, &ModelNotFound{Id: id, Model: "article", Code: http.StatusNotFound}
		}

		var c content
		e = json.Unmarshal([]byte(article.Content), &c)

		if e != nil {
			return nil, e.(BaseError)
		}

		article.Content = c.Html
		article.ContentMd = c.Md
		if article.TopAt != nil && !article.TopAt.IsZero() {
			article.IsTop = true
		}

		bytes, _ := json.Marshal(article)
		_, e = dao.Redis.Set(key, string(bytes), 86400*time.Second).Result()

		if e != nil {
			return nil, e.(BaseError)
		}

		return article, nil
	} else {
		//log.Println(e)
		var article models.Article

		e := json.Unmarshal([]byte(s), &article)

		if e != nil {
			log.Println(e)
			return nil, &ModelNotFound{
				Code: 404,
			}
		}

		return &article, nil
	}
}

func (dao *dao) HomeArticles() []models.Article {
	var articles []models.Article
	dao.DB.
		Preload("Category").
		Select([]string{"category_id", "author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("id DESC").
		Limit(3).
		Find(&articles)

	return articles
}

func (dao *dao) TopArticles() []*models.Article {
	var articles []*models.Article
	dao.DB.
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Where("top_at is not null").
		Order("top_at DESC").
		Find(&articles)
	for _, v := range articles {
		v.IsTop = true
	}

	return articles
}

func (dao *dao) NewestArticles() []models.Article {
	var articles []models.Article
	dao.DB.
		Select("author_id,id,top_at,head_image,title,`desc`,created_at").
		Where("display = ?", true).
		Order("id DESC").
		Limit(13).
		Find(&articles)
	for _, v := range articles {
		if v.TopAt != nil && !v.TopAt.IsZero() {
			v.IsTop = true
		}
	}

	return articles
}

func (dao *dao) PopularArticles() []models.Article {
	var articles []models.Article
	dao.DB.
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("RAND()").
		Limit(8).
		Find(&articles)

	return articles
}

func (dao *dao) Search(q string) []*models.Article {
	es := config.Config.ES
	multiMatch := es.MultiMatch
	highlight := es.Highlight
	query := elastic.NewMultiMatchQuery(fmt.Sprintf(multiMatch.Query.MultiMatch.Query, q), multiMatch.Query.MultiMatch.Fields...).Analyzer(multiMatch.Query.MultiMatch.Analyzer)

	result, e := dao.ES.Search().
		Index("article_index").
		Highlight(highlight).
		Query(query).
		FetchSource(false).
		Size(10000).
		Pretty(true).
		Do(context.Background())
	if e != nil {
		log.Println("err", e)
	}

	var highIdMap = map[string]models.Highlight{}
	for _, v := range result.Hits.Hits {
		var h models.Highlight
		for field, highlight := range v.Highlight {
			switch field {
			case "title":
				h.Title = strings.Join(highlight, "......")
			case "desc":
				h.Desc = strings.Join(highlight, "......")
			case "content":
				h.Content = strings.Join(highlight, "......")
			case "article_category.name":
				h.Category = strings.Join(highlight, "......")
			case "tags":
				h.Tags = strings.Join(highlight, ", ")
			}
		}
		highIdMap[v.Id] = h
	}

	var hitArticleIds []int

	for _, v := range result.Hits.Hits {
		i, _ := strconv.Atoi(v.Id)
		hitArticleIds = append(hitArticleIds, i)
	}

	var articles []*models.Article
	dao.DB.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Select([]string{"id", "author_id", "category_id", "`desc`", "title", "head_image", "created_at", "display"}).
		Where("id in (?)", hitArticleIds).
		Order("id DESC").
		Where("display = ?", true).
		Find(&articles)
	for _, v := range articles {
		if data, ok := highIdMap[strconv.Itoa(v.Id)]; ok {
			v.Highlight = data
		}
	}

	return articles
}
