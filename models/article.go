package models

import "github.com/youngduc/go-blog/hello/resources"

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
func Get(id int64) *resources.Article {
	var article resources.Article
	db.Table("articles").Select("id,author_id,category_id,content,`desc`,title,head_image,display,top_at,created_at").
		Where("id = ?", id).
		First(&article)

	return &article
}
