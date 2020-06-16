package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v6"
	"github.com/youngduc/go-blog/config"
	"github.com/youngduc/go-blog/models"
	"github.com/youngduc/go-blog/services"
	"github.com/youngduc/go-blog/utils"
	"github.com/youngduc/go-blog/utils/errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ArticleController struct {
}

func (*ArticleController) Index(ctx *gin.Context) {
	var (
		article = &models.Article{}
		page    = utils.GetQueryIntValueWithDefault(ctx, "page", 1)
		perPage = utils.GetQueryIntValueWithDefault(ctx, "page_size", 15)
	)

	Success(ctx, 200, article.Paginate(page, perPage))
}

func (a *ArticleController) Show(ctx *gin.Context) {
	var (
		trending services.Trending
		article  = &models.Article{}
	)

	id, _ := strconv.Atoi(ctx.Param("id"))
	key := a.cacheKey("article:" + strconv.Itoa(id))
	s, e := redisClient.Get(key).Result()

	if e == redis.Nil {
		e := article.Find(id)

		if e != nil && e == gorm.ErrRecordNotFound {
			Fail(ctx, &errors.ModelNotFound{Id: id, Model: "article", Code: http.StatusNotFound})
			return
		}

		var c models.ArticleContent
		e = json.Unmarshal([]byte(article.Content), &c)

		if e != nil {
			Fail(ctx, e.(errors.BaseError))
			return
		}

		article.Content = c.Html
		if article.TopAt != nil && !article.TopAt.IsZero() {
			article.IsTop = true
		}

		bytes, _ := json.Marshal(article)
		_, e = redisClient.Set(key, string(bytes), 86400*time.Second).Result()

		if e != nil {
			Fail(ctx, e.(errors.BaseError))
			return
		}
	} else {
		e := json.Unmarshal([]byte(s), &article)

		if e != nil {
			Fail(ctx, &errors.ModelNotFound{Code: 404})
			return
		}

		article.ContentMd = ""
	}

	trending.Push(article.Id)

	Success(ctx, 200, gin.H{
		"data": article,
	})
}

func (*ArticleController) Search(ctx *gin.Context) {
	type HH struct {
		Highlight models.Highlight
		Hits      *elastic.SearchHit
	}

	var (
		highIdMap     = map[string]HH{}
		hitArticleIds []int
		articles      []*models.Article
	)

	es := config.Cfg.ES
	multiMatch := es.MultiMatch
	highlight := es.Highlight
	query := elastic.NewMultiMatchQuery(fmt.Sprintf(multiMatch.Query.MultiMatch.Query, ctx.Query("q")), multiMatch.Query.MultiMatch.Fields...).Analyzer(multiMatch.Query.MultiMatch.Analyzer)

	result, e := esClient.Search().
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

	for _, v := range result.Hits.Hits {
		var h = HH{}
		for field, highlight := range v.Highlight {
			switch field {
			case "title":
				h.Highlight.Title = strings.Join(highlight, "......")
			case "desc":
				h.Highlight.Desc = strings.Join(highlight, "......")
			case "content":
				h.Highlight.Content = strings.Join(highlight, "......")
			case "article_category.name":
				h.Highlight.Category = strings.Join(highlight, "......")
			case "tags":
				h.Highlight.Tags = strings.Join(highlight, ", ")
			}
		}
		h.Hits = v
		highIdMap[v.Id] = h
	}

	for _, v := range result.Hits.Hits {
		i, _ := strconv.Atoi(v.Id)
		hitArticleIds = append(hitArticleIds, i)
	}

	dbClient.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Select([]string{"id", "author_id", "category_id", "`desc`", "title", "head_image", "created_at", "display"}).
		Where("id in (?)", hitArticleIds).
		Where("display = ?", true).
		Find(&articles)

	for _, v := range articles {
		if data, ok := highIdMap[strconv.Itoa(v.Id)]; ok {
			v.Highlight = data.Highlight
		}
	}

	sort.Slice(articles, func(i, j int) bool {
		a := highIdMap[strconv.Itoa(articles[i].Id)].Hits.Score
		b := highIdMap[strconv.Itoa(articles[j].Id)].Hits.Score

		return *a > *b
	})

	Success(ctx, 200, gin.H{
		"data": articles,
	})
}

func (*ArticleController) Home(ctx *gin.Context) {
	var articles []models.Article
	dbClient.
		Preload("Category").
		Select([]string{"category_id", "author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("id DESC").
		Limit(3).
		Find(&articles)

	Success(ctx, http.StatusOK, gin.H{
		"data": articles,
	})
}

func (*ArticleController) Newest(ctx *gin.Context) {
	var articles []models.Article
	dbClient.
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

	Success(ctx, http.StatusOK, gin.H{
		"data": articles,
	})
}

func (*ArticleController) Popular(ctx *gin.Context) {
	var articles []models.Article
	dbClient.
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("RAND()").
		Limit(8).
		Find(&articles)

	Success(ctx, http.StatusOK, gin.H{
		"data": articles,
	})
}

func (*ArticleController) Trending(ctx *gin.Context) {
	var (
		trending        services.Trending
		tendingArticles []models.Article
		m               = map[int]int{}
	)

	get := trending.Get()
	for k, v := range get {
		m[v] = k
	}
	dbClient.Select("id,head_image,title,`desc`,display,created_at").
		Where("id in (?)", get).
		Find(&tendingArticles)

	sort.Slice(tendingArticles, func(i, j int) bool {
		return m[tendingArticles[i].Id] < m[tendingArticles[j].Id]
	})

	Success(ctx, http.StatusOK, gin.H{
		"data": tendingArticles,
	})
}

func (*ArticleController) Top(ctx *gin.Context) {
	var articles []*models.Article
	dbClient.
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Where("top_at is not null").
		Order("top_at DESC").
		Find(&articles)
	for _, v := range articles {
		v.IsTop = true
	}

	Success(ctx, http.StatusOK, gin.H{
		"data": articles,
	})
}

func (*ArticleController) cacheKey(key string) string {
	return fmt.Sprintf("%s%s", config.Cfg.RedisPrefix, key)
}
