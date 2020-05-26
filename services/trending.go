package services

import (
	"github.com/youngduc/go-blog/models/dao"
	"strconv"
)

type Trending struct {
}

func (t *Trending) Get() []int {
	var res = make([]int, 12)
	all := t.AllKeys()
	invisible := t.GetInvisibleIds()

	var flag bool
	for _, v := range all {
		flag = true
		for _, v1 := range invisible {
			if v==v1 {
				flag = false
			}
		}

		if flag {
			i, e := strconv.Atoi(v)
			if e==nil {
				res = append(res, i)
			}
		}
	}

	return res[:12]
}

func (t *Trending) Push(id int) {
	dao.Dao.Redis.ZIncrBy(t.CacheKey(), 1, strconv.Itoa(id))
}
func (t *Trending) Remove(id int) {
	dao.Dao.Redis.ZRem(t.CacheKey(), strconv.Itoa(id))
}
func (*Trending) CacheKey() string {
	return "trending_articles"
}
func (t *Trending) Reset() {
	dao.Dao.Redis.Del(t.CacheKey(), t.InvisibleKey())
}
func (t *Trending) HasKey(id int) bool {
	i := dao.Dao.Redis.ZRank(t.CacheKey(), strconv.Itoa(id)).Val()

	return i != 0
}
func (t *Trending) GetInvisibleIds() []string {
	return dao.Dao.Redis.SMembers(t.InvisibleKey()).Val()
}
func (t *Trending) AddInvisible(id int) bool {
	return dao.Dao.Redis.SAdd(t.InvisibleKey(), strconv.Itoa(id)).Val() != 0
}
func (t *Trending) RemoveInvisible(id int) bool {
	return dao.Dao.Redis.SRem(t.InvisibleKey(), strconv.Itoa(id)).Val() != 0
}
func (*Trending) InvisibleKey() string {
	return "invisible_articles"
}
func (t *Trending) AllKeys() []string {
	return dao.Dao.Redis.ZRevRange(t.InvisibleKey(), 0, -1).Val()

}
