package models

const (
	UserableTypeSocialiteUser = "App\\SocialiteUser"
	UserableTypeUser          = "App\\User"
)

type CommentAuthor struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Avatar    string   `json:"avatar"`
}

type Comment struct {
	Model

	Visitor      string `json:"visitor"`
	Body         string `json:"body" gorm:"-"`
	Content      string `json:"content"`
	ArticleId    int64  `json:"article_id"`
	CommentId    int64  `json:"comment_id"`
	UserableId   int64  `json:"userable_id" gorm:"default:0"`
	UserableType string `json:"userable_type" gorm:"default:''"`

	Author CommentAuthor `json:"author"`

	Replies []*Comment `json:"replies"`
	//Article *Article
}
