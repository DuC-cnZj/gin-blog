package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/olivere/elastic/v6"
	"github.com/youngduc/go-blog/hello/config"
	"github.com/youngduc/go-blog/hello/models"
	"log"
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

	dao.db.
		Preload("Author").
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Offset(offset).
		Limit(perPage).
		Find(&articles)

	dao.db.Table("articles").Where("display = ?", true).Count(&count)

	return map[string]interface{}{
		"data": articles,
		"meta": models.Paginator{
			Total:       count,
			CurrentPage: page,
			PerPage:     perPage,
		},
	}
}

func cacheKey(key string) string {
	return fmt.Sprintf("%s%s", config.Config.RedisPrefix, key)
}

func (dao *dao) ShowArticle(id int) interface{} {
	key := cacheKey("article:" + strconv.Itoa(id))
	s, e := dao.redis.Get(key).Result()
	if e == redis.Nil {
		article := &models.Article{}
		dao.db.
			Preload("Author").
			Preload("Category").
			Preload("Tags").
			Where("id = ?", id).
			Where("display = ?", true).
			Find(article)

		var c content
		e := json.Unmarshal([]byte(article.Content), &c)

		if e != nil {
			log.Fatal(e)
		}

		article.Content = c.Html
		article.ContentMd = c.Md
		if article.TopAt != nil && !article.TopAt.IsZero() {
			article.IsTop = true
		}

		bytes, _ := json.Marshal(article)
		result, e := dao.redis.Set(key, string(bytes), 86400*time.Second).Result()
		if e != nil {
			log.Println(e)
		}
		log.Println("redis result: ", result)

		return article
	} else {
		var article models.Article

		e := json.Unmarshal([]byte(s), &article)

		log.Println(s, e)

		return article
	}
}

func (dao *dao) HomeArticles() []models.Article {
	var articles []models.Article
	dao.db.
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("id DESC").
		Limit(3).
		Find(&articles)

	return articles
}

func (dao *dao) TopArticles() []*models.Article {
	var articles []*models.Article
	dao.db.
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

func (dao *dao) NewestArticles() []*models.Article {
	var articles []*models.Article
	dao.db.
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("id DESC").
		Limit(13).
		Find(&articles)
	for _, v := range articles {
		if !v.TopAt.IsZero() {
			v.IsTop = true
		}
	}

	return articles
}

func (dao *dao) PopularArticles() []models.Article {
	var articles []models.Article
	dao.db.
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
	log.Println(query)

	result, e := dao.es.Search().
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
	dao.db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("id in (?)", hitArticleIds).
		Where("display = ?", true).
		Find(&articles)
	for _, v := range articles {
		var c content
		e := json.Unmarshal([]byte(v.Content), &c)

		if e != nil {
			log.Fatal(e)
		}
		v.Content = c.Html
		v.ContentMd = c.Md
		if data, ok := highIdMap[strconv.Itoa(v.Id)]; ok {
			v.Highlight = data
		}
	}

	return articles
}
