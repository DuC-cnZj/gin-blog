package models

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
	Highlight Highlight
}

type Highlight struct {
	Title string `json:"title"`
	Tags string `json:"tags"`
	Category string `json:"category"`
	Content string `json:"content"`
	Desc string `json:"desc"`
}

//$table->increments('id');
//$table->integer('author_id');
//$table->integer('category_id');
//$table->longtext('content')->nullable();
//$table->string('desc');
//$table->string('title');
//$table->string('head_image');
//$table->boolean('display')->default(true);
//$table->timestamp('top_at')->nullable();
//$table->timestamps();
