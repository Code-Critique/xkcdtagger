package redis

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/Code-Critique/xkcdtagger"
	"github.com/go-redis/redis"
)

// Client represents the redis client
type Client struct {
	*redis.Client
}

// NewClient returns a new client
func NewClient() *Client {
	return &Client{
		Client: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0, // use default DB
		}),
	}
}

// ListComics lists all comics in redis
func (c *Client) ListComics() ([]xkcdtagger.Comic, error) {
	sc := c.HVals("comics")
	strings, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var out []xkcdtagger.Comic

	for _, s := range strings {

		log.Println(s)

		var item xkcdtagger.Comic

		err := json.Unmarshal([]byte(s), &item)

		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	return out, nil
}

// GetComic gets comic by id
func (c *Client) GetComic(id xkcdtagger.ComicID) (*xkcdtagger.Comic, error) {
	key := strconv.Itoa(int(id))

	sc := c.HGet("comics", key)

	s, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var item xkcdtagger.Comic

	err = json.Unmarshal([]byte(s), &item)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetTagsForComic gets tags for comic
func (c *Client) GetTagsForComic(id xkcdtagger.ComicID) ([]xkcdtagger.Tag, error) {
	key := strconv.Itoa(int(id))

	sc := c.HGet("comics", key)

	s, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var item xkcdtagger.Comic

	err = json.Unmarshal([]byte(s), &item)

	if err != nil {
		return nil, err
	}

	return item.Tags, nil
}

// GetTags gets all tags
func (c *Client) GetTags() ([]xkcdtagger.Tag, error) {
	sc := c.SMembers("tags")

	strings, err := sc.Result()

	if err != nil {
		return nil, err
	}

	var out []xkcdtagger.Tag

	for _, s := range strings {
		var item xkcdtagger.Tag

		err := json.Unmarshal([]byte(s), &item)

		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	return out, nil
}
