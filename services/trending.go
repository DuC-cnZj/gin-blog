package services

import (
	"github.com/youngduc/go-blog/config"
	"strconv"
)

type Trending struct {
}

func (t *Trending) Get() []int {
	var res = make([]int, 0)
	all := t.AllKeys()
	invisible := t.GetInvisibleIds()

	var flag bool
	for _, v := range all {
		flag = false
		for _, v1 := range invisible {
			if v == v1 {
				flag = true
			}
		}

		if !flag {
			i, e := strconv.Atoi(v)
			if e == nil && i != 0 {
				res = append(res, i)
			}
		}
	}

	if len(res) > 12 {
		return res[:12]
	}

	return res
}

func (t *Trending) Push(id int) {
	config.Conn.RedisClient.ZIncrBy(t.CacheKey(), 1, strconv.Itoa(id))
}
func (t *Trending) Remove(id int) {
	config.Conn.RedisClient.ZRem(t.CacheKey(), strconv.Itoa(id))
}
func (*Trending) CacheKey() string {
	return "trending_articles"
}
func (t *Trending) Reset() {
	config.Conn.RedisClient.Del(t.CacheKey(), t.InvisibleKey())
}
func (t *Trending) HasKey(id int) bool {
	i := config.Conn.RedisClient.ZRank(t.CacheKey(), strconv.Itoa(id)).Val()

	return i != 0
}
func (t *Trending) GetInvisibleIds() []string {
	return config.Conn.RedisClient.SMembers(t.InvisibleKey()).Val()
}
func (t *Trending) AddInvisible(id int) bool {
	return config.Conn.RedisClient.SAdd(t.InvisibleKey(), strconv.Itoa(id)).Val() != 0
}
func (t *Trending) RemoveInvisible(id int) bool {
	return config.Conn.RedisClient.SRem(t.InvisibleKey(), strconv.Itoa(id)).Val() != 0
}
func (*Trending) InvisibleKey() string {
	return "invisible_articles"
}
func (t *Trending) AllKeys() []string {
	return config.Conn.RedisClient.ZRevRange(t.CacheKey(), 0, -1).Val()

}
