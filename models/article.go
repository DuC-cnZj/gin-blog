package models

type ArticleContent struct {
	Md   string `json:"md"`
	Html string `json:"html"`
}

type Article struct {
	Model

	IsTop      bool      `json:"is_top"`
	AuthorId   int64     `json:"author_id"`
	HeadImage  string    `json:"head_image"`
	Content    string    `json:"content"`
	ContentMd  string    `json:"content_md"`
	Title      string    `json:"title"`
	Desc       string    `json:"desc"`
	Display    bool      `json:"display" gorm:"default:true"`
	CategoryId int64     `json:"category_id"`
	TopAt      *JSONTime `json:"top_at"`

	Author User `json:"author" gorm:"foreignkey:AuthorId"`

	Category Category `json:"category"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tag;"`

	Comments []*Comment `json:"comments"`
	//recommendArticles
	Highlight Highlight `json:"highlight"`
}

func (*Article) Paginate(page, perPage int) map[string]interface{} {
	var articles []Article
	var count int
	offset := (page - 1) * perPage

	db.
		Preload("Author").
		Select([]string{"author_id", "id", "top_at", "head_image", "title", "`desc`", "created_at"}).
		Where("display = ?", true).
		Order("id desc").
		Offset(offset).
		Limit(perPage).
		Find(&articles)

	db.Table("articles").Where("display = ?", true).Count(&count)

	return map[string]interface{}{
		"data": articles,
		"meta": Paginator{
			Total:       count,
			CurrentPage: page,
			PerPage:     perPage,
		},
		"links": map[string]string{},
	}
}

func (article *Article) Find(id int) error {
	return db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("id = ?", id).
		Where("display = ?", true).
		Find(article).
		Error
}

type Highlight struct {
	Title    string `json:"title"`
	Tags     string `json:"tags"`
	Category string `json:"category"`
	Content  string `json:"content"`
	Desc     string `json:"desc"`
}
